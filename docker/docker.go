// Copyright (c) 2000-present Liferay, Inc. All rights reserved.
//
// This library is free software; you can redistribute it and/or modify it under
// the terms of the GNU Lesser General Public License as published by the Free
// Software Foundation; either version 2.1 of the License, or (at your option)
// any later version.
//
// This library is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more
// details.

package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	types "github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	filters "github.com/docker/docker/api/types/filters"
	mount "github.com/docker/docker/api/types/mount"
	client "github.com/docker/docker/client"
	nat "github.com/docker/go-connections/nat"
	internal "github.com/mdelapenya/lpn/internal"
	liferay "github.com/mdelapenya/lpn/liferay"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var instance *client.Client

type imagePullResponse struct {
	ID             string `json:"id"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
	Status string `json:"status"`
}

func buildPortBinding(port string, ip string) []nat.PortBinding {
	return []nat.PortBinding{
		nat.PortBinding{
			HostPort: port,
			HostIP:   ip,
		},
	}
}

func buildTarForDeployment(file *os.File) (bytes.Buffer, error) {
	fileInfo, _ := file.Stat()

	var buffer bytes.Buffer
	tarWriter := tar.NewWriter(&buffer)
	err := tarWriter.WriteHeader(&tar.Header{
		Name: fileInfo.Name(),
		Mode: 0777,
		Size: int64(fileInfo.Size()),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"fileInfoName": fileInfo.Name(),
			"size":         fileInfo.Size(),
			"error":        err,
		}).Error("Could not build TAR header")
		return bytes.Buffer{}, fmt.Errorf("Could not build TAR header: %v", err)
	}

	b, err := ioutil.ReadFile(file.Name())
	tarWriter.Write(b)
	defer tarWriter.Close()

	return buffer, nil
}

// CheckDocker checks if Docker is installed
func CheckDocker() bool {
	_, _, err := GetDockerVersion()
	if err != nil {
		return false
	}

	return true
}

// CheckDockerContainerExists checks if the container exists by label
// It looks for containers with the label "lpn-container-name" matching the provided name
func CheckDockerContainerExists(containerName string) bool {
	dockerClient := getDockerClient()

	// Use label filter to find containers by lpn-container-name label
	containers, err := dockerClient.ContainerList(
		context.Background(), containertypes.ListOptions{
			All: true,
			Filters: filters.NewArgs(
				filters.Arg("label", "lpn-container-name="+containerName),
			),
		})

	if err != nil {
		return false
	}

	return len(containers) > 0
}

// GetContainerIDByLabel returns the container ID for a container with the given label
// Returns empty string if no container found
func GetContainerIDByLabel(containerName string) string {
	dockerClient := getDockerClient()

	// Use label filter to find containers by lpn-container-name label
	containers, err := dockerClient.ContainerList(
		context.Background(), containertypes.ListOptions{
			All: true,
			Filters: filters.NewArgs(
				filters.Arg("label", "lpn-container-name="+containerName),
			),
		})

	if err != nil || len(containers) == 0 {
		return ""
	}

	// Return the ID of the first matching container
	return containers[0].ID
}

// CheckDockerImageExists checks if the image is already present
func CheckDockerImageExists(dockerImage string) bool {
	dockerClient := getDockerClient()

	dockerImage = strings.ReplaceAll(dockerImage, "docker.io/", "")

	imageInspect, _, err := dockerClient.ImageInspectWithRaw(context.Background(), dockerImage)

	if err != nil {
		return false
	}

	for i := range imageInspect.RepoTags {
		tag := imageInspect.RepoTags[i]

		if dockerImage == tag {
			return true
		}
	}
	return false
}

// CopyFileToContainer copies a file to the running container
func CopyFileToContainer(image liferay.Image, path string) error {
	dockerClient := getDockerClient()

	log.WithFields(log.Fields{
		"file":   path,
		"target": image.GetDeployFolder(),
	}).Debug("Deploying [" + path + "] to " + image.GetDeployFolder())

	_, err := dockerClient.ContainerStatPath(
		context.Background(), image.GetContainerName(), image.GetDeployFolder())
	if err != nil {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
			"target":    image.GetDeployFolder(),
			"error":     err,
		}).Error("Could not get directory in the container")
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  path,
			"error": err,
		}).Error("Could not open file to deploy")
		return err
	}
	defer file.Close()

	buffer, err := buildTarForDeployment(file)
	if err != nil {
		return err
	}

	err = dockerClient.CopyToContainer(
		context.Background(), image.GetContainerName(), image.GetDeployFolder(),
		&buffer, containertypes.CopyToContainerOptions{AllowOverwriteDirWithFile: true})

	if err == nil {
		targetFilePath := filepath.Join(image.GetDeployFolder(), filepath.Base(file.Name()))
		owner := image.GetUser()

		cmd := []string{"chown", owner + ":" + owner, targetFilePath}

		execCommandIntoContainer(image.GetContainerName(), cmd)
	} else {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
			"deployDir": image.GetDeployFolder(),
			"error":     err,
		}).Error("Could not copy file to container")
	}

	return err
}

func execCommandIntoContainer(containerName string, cmd []string) error {
	dockerClient := getDockerClient()

	containerID := GetContainerIDByLabel(containerName)
	if containerID == "" {
		err := fmt.Errorf("container not found")
		log.WithFields(log.Fields{
			"container": containerName,
			"error":     err,
		}).Error("Could not find container")
		return err
	}

	response, err := dockerClient.ContainerExecCreate(
		context.Background(), containerID, containertypes.ExecOptions{
			User:         "root",
			Tty:          false,
			AttachStdin:  false,
			AttachStderr: false,
			AttachStdout: false,
			Detach:       true,
			Cmd:          cmd,
		})

	if err != nil {
		log.WithFields(log.Fields{
			"container": containerName,
			"cmd":       cmd,
			"error":     err,
		}).Error("Could not create command in the container")
		return err
	}

	err = dockerClient.ContainerExecStart(
		context.Background(), response.ID, containertypes.ExecStartOptions{
			Detach: true,
			Tty:    false,
		})
	if err != nil {
		log.WithFields(log.Fields{
			"container": containerName,
			"cmd":       cmd,
			"detach":    true,
			"tty":       false,
			"error":     err,
		}).Error("Could not start command in the container")
	}

	return err
}

func getDockerClient() *client.Client {
	if instance != nil {
		return instance
	}

	instance, err := client.NewEnvClient()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not get Docker client")
	}

	return instance
}

// GetDockerImageFromRunningContainer gets the image name of the container
func GetDockerImageFromRunningContainer(image liferay.Image) (string, error) {
	dockerClient := getDockerClient()

	containers, err := dockerClient.ContainerList(
		context.Background(), containertypes.ListOptions{All: true})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not list all containers")
		return "", err
	}

	for _, container := range containers {
		containerName := "/" + image.GetContainerName()

		if containerName == container.Names[0] {
			log.WithFields(log.Fields{
				"container": image.GetContainerName(),
			}).Debug("Container found!")
			return container.Image, nil
		}
	}

	err = errors.New("We could not find the container among the running containers")
	log.WithFields(log.Fields{
		"container": image.GetContainerName(),
		"error":     err,
	}).Debug("We could not find the container among the running containers")

	return "", err
}

// GetDockerVersion returns the output of Docker version
func GetDockerVersion() (string, types.Version, error) {
	dockerClient := getDockerClient()

	serverVersion, err := dockerClient.ServerVersion(context.Background())

	return dockerClient.ClientVersion(), serverVersion, err
}

// inspect inspects a container by label
func inspect(containerName string) types.ContainerJSON {
	dockerClient := getDockerClient()

	containerID := GetContainerIDByLabel(containerName)
	if containerID == "" {
		log.WithFields(log.Fields{
			"container": containerName,
		}).Fatal("Container not found")
	}

	containerJSON, err := dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.WithFields(log.Fields{
			"container": containerName,
			"error":     err,
		}).Fatal("The container could not be inspected")
	}

	return containerJSON
}

// GetTomcatPort gets Tomcat port from running instance
func GetTomcatPort(image liferay.Image) string {
	containerJSON := inspect(image.GetContainerName())

	hostConfig := containerJSON.HostConfig

	portBindings := hostConfig.PortBindings

	tomcatPortBinding := portBindings["8080/tcp"]

	return tomcatPortBinding[0].HostPort
}

// LogContainer show logs of a container in tail mode
func LogContainer(image liferay.Image) {
	dockerClient := getDockerClient()

	reader, err := dockerClient.ContainerLogs(
		context.Background(), image.GetContainerName(),
		containertypes.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
			"error":     err,
		}).Fatal("Could not get container logs")
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
			"error":     err,
		}).Fatal("Error following container logs")
	}
}

// PsFilterByLabel Retrieves all containers with a label
func PsFilterByLabel(label string) ([]types.Container, error) {
	dockerClient := getDockerClient()

	filters := filters.NewArgs()
	filters.Add("label", label)

	return dockerClient.ContainerList(
		context.Background(), containertypes.ListOptions{
			Size:    true,
			All:     true,
			Since:   "container",
			Filters: filters,
		})
}

// PullDockerImage downloads the image
func PullDockerImage(dockerImage string) {
	dockerClient := getDockerClient()

	log.WithFields(log.Fields{
		"dockerImage": dockerImage,
	}).Debug("Pulling Docker image.")

	out, err := dockerClient.ImagePull(
		context.Background(), dockerImage, image.PullOptions{})

	if err == nil {
		parseImagePull(out)
	} else {
		log.WithFields(log.Fields{
			"dockerImage": dockerImage,
			"error":       err,
		}).Fatal("The image could not be pulled")
	}
}

func parseImagePull(pullResp io.ReadCloser) {
	d := json.NewDecoder(pullResp)
	for {
		var pullResult imagePullResponse
		if err := d.Decode(&pullResult); err != nil {
			break
		}

		log.WithFields(log.Fields{
			"id":       pullResult.ID,
			"status":   pullResult.Status,
			"progress": pullResult.Progress,
		}).Infof("%s %s %s\n", pullResult.ID, pullResult.Status, pullResult.Progress)
	}
}

// RemoveDockerContainer removes a running container, and its stack
func RemoveDockerContainer(image liferay.Image) error {
	dockerClient := getDockerClient()

	var err error
	label := "lpn-type=" + image.GetType()

	containers, err := PsFilterByLabel(label)

	if len(containers) == 0 {
		err = errors.New("Error response from daemon: No such container: " + image.GetContainerName())

		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
			"label":     label,
			"error":     err,
		}).Error("Could not filter container by label")

		return err
	}

	for _, container := range containers {
		name := strings.TrimLeft(container.Names[0], "/")
		err = dockerClient.ContainerRemove(
			context.Background(), name, containertypes.RemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			})
		if err == nil {
			log.WithFields(log.Fields{
				"container": name,
			}).Info("Container has been removed")
		}
	}

	return err
}

// RemoveDockerImage removes a docker image
func RemoveDockerImage(dockerImageName string) error {
	dockerClient := getDockerClient()

	_, err := dockerClient.ImageRemove(
		context.Background(), dockerImageName,
		image.RemoveOptions{
			Force: true,
		})
	if err != nil {
		log.WithFields(log.Fields{
			"image": dockerImageName,
			"error": err,
		}).Warn("Impossible to remove the image")

		return err
	}

	log.WithFields(log.Fields{
		"image": dockerImageName,
	}).Info("Image has been removed")

	return nil
}

// RunDatabaseDockerImage runs the image, setting the HTTP port and a volume for the data folder
// Now uses testcontainers-go for container management with persistent, reusable containers.
//
// Key features:
// - WithReuseByName: Container persists and is reused across invocations
// - No automatic cleanup: Container lifecycle managed explicitly
// - Production-ready with built-in wait strategies
//
// For additional persistence guarantees, you can disable Ryuk globally:
// - Set TESTCONTAINERS_RYUK_DISABLED=true environment variable, OR
// - Create .testcontainers.properties with ryuk.disabled=true
func RunDatabaseDockerImage(image DatabaseImage) error {
	ctx := context.Background()
	containerName := image.GetContainerName()

	// Check if container already exists and is running
	if CheckDockerContainerExists(containerName) {
		log.WithFields(log.Fields{
			"container": containerName,
		}).Debug("Not starting a new container because it's already running")

		return nil
	}

	// Create mount path for data persistence
	volumePath := filepath.Join(internal.LpnWorkspace, containerName)
	os.MkdirAll(volumePath, os.ModePerm)
	
	log.WithFields(log.Fields{
		"container": containerName,
		"volume":    volumePath,
	}).Debug("Mounting database data folder")

	var container testcontainers.Container
	var err error

	// Use appropriate module based on database type
	switch image.GetType() {
	case "mysql":
		container, err = mysql.Run(
			ctx,
			image.GetFullyQualifiedName(),
			mysql.WithDatabase(DBName),
			mysql.WithUsername(DBUser),
			mysql.WithPassword(DBPassword),
			// Mount volume for data persistence
			testcontainers.WithMounts(
				testcontainers.BindMount(volumePath, testcontainers.ContainerMountTarget(image.GetDataFolder())),
			),
			// Add labels for identification - use label instead of container name
			testcontainers.WithLabels(map[string]string{
				"lpn-container-name": containerName,
				"db-type":            image.GetType(),
				"lpn-type":           image.GetLpnType(),
			}),
			// Wait for database to be ready
			testcontainers.WithWaitStrategy(
				wait.ForLog("port: 3306  MySQL Community Server"),
			),
		)

	case "postgresql":
		container, err = postgres.Run(
			ctx,
			image.GetFullyQualifiedName(),
			postgres.WithDatabase(DBName),
			postgres.WithUsername(DBUser),
			postgres.WithPassword(DBPassword),
			// Mount volume for data persistence
			testcontainers.WithMounts(
				testcontainers.BindMount(volumePath, testcontainers.ContainerMountTarget(image.GetDataFolder())),
			),
			// Add labels for identification - use label instead of container name
			testcontainers.WithLabels(map[string]string{
				"lpn-container-name": containerName,
				"db-type":            image.GetType(),
				"lpn-type":           image.GetLpnType(),
			}),
		)

	default:
		return fmt.Errorf("unsupported database type: %s", image.GetType())
	}

	if err != nil {
		log.WithFields(log.Fields{
			"container": containerName,
			"image":     image.GetFullyQualifiedName(),
			"error":     err,
		}).Fatal("Could not create database container")
		return err
	}

	log.WithFields(log.Fields{
		"container": containerName,
		"image":     image.GetFullyQualifiedName(),
	}).Debug("Database container has been started")

	// Store container reference for later use (optional)
	_ = container

	return nil
}

// RunLiferayDockerImage runs the image, setting the HTTP and GoGoShell ports for bundle, debug mode, and
// jvmMemory if needed
func RunLiferayDockerImage(
	image liferay.Image, database DatabaseImage, httpPort int, gogoShellPort int, enableDebug bool,
	debugPort int, memory string) error {

	if CheckDockerContainerExists(image.GetContainerName()) {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
		}).Debug("The container is running.")

		_ = RemoveDockerContainer(image)
	}

	port := fmt.Sprintf("%d", httpPort)
	gogoPort := fmt.Sprintf("%d", gogoShellPort)
	debuggerPort := fmt.Sprintf("%d", debugPort)

	environmentVariables := []string{}

	exposedPorts := map[nat.Port]struct{}{
		"8080/tcp":  {},
		"11311/tcp": {},
	}

	portBindings := make(map[nat.Port][]nat.PortBinding)

	portBindings["8080/tcp"] = buildPortBinding(port, "0.0.0.0")
	portBindings["11311/tcp"] = buildPortBinding(gogoPort, "0.0.0.0")

	if enableDebug {
		var port9000 struct{}
		exposedPorts["9000/tcp"] = port9000

		portBindings["9000/tcp"] = buildPortBinding(debuggerPort, "0.0.0.0")

		debugEnvVarName := ""

		switch imageType := image.(type) {
		case liferay.CE, liferay.Commerce, liferay.DXP, liferay.Nightly:
			debugEnvVarName = "LIFERAY_JPDA_ENABLED"
		case liferay.Release:
			debugEnvVarName = "DEBUG_MODE"
		default:
			log.Fatalln("Non supported type", imageType)
		}

		environmentVariables = append(environmentVariables, debugEnvVarName+"=true")
	}

	if memory != "" {
		environmentVariables = append(environmentVariables, "LIFERAY_JVM_OPTS="+memory)
	}

	PullDockerImage(image.GetFullyQualifiedName())

	dockerClient := getDockerClient()

	links := []string{}

	if database != nil {
		link := database.GetContainerName() + ":" + "db"
		links = append(links, link)

		RunDatabaseDockerImage(database)

		environmentVariables = append(environmentVariables, "LIFERAY_JDBC_PERIOD_DEFAULT_PERIOD_DRIVER_UPPERCASEC_LASS_UPPERCASEN_AME="+database.GetJDBCConnection().DriverClassName)
		environmentVariables = append(environmentVariables, "LIFERAY_JDBC_PERIOD_DEFAULT_PERIOD_PASSWORD="+database.GetJDBCConnection().Password)
		environmentVariables = append(environmentVariables, "LIFERAY_JDBC_PERIOD_DEFAULT_PERIOD_URL="+database.GetJDBCConnection().URL)
		environmentVariables = append(environmentVariables, "LIFERAY_JDBC_PERIOD_DEFAULT_PERIOD_USERNAME="+database.GetJDBCConnection().User)

		// retry JDBC in case the database is slower
		environmentVariables = append(environmentVariables, "LIFERAY_RETRY_PERIOD_JDBC_PERIOD_ON_PERIOD_STARTUP_PERIOD_DELAY=5")
		environmentVariables = append(environmentVariables, "LIFERAY_RETRY_PERIOD_JDBC_PERIOD_ON_PERIOD_STARTUP_PERIOD_MAX_PERIOD_RETRIES=5")
	}

	containerCreationResponse, err := dockerClient.ContainerCreate(
		context.Background(),
		&containertypes.Config{
			Image:        image.GetFullyQualifiedName(),
			Env:          environmentVariables,
			ExposedPorts: exposedPorts,
			Labels: map[string]string{
				"lpn-type": image.GetType(),
			},
		},
		&containertypes.HostConfig{
			Links:        links,
			PortBindings: portBindings,
			Mounts:       []mount.Mount{},
		},
		nil, // NetworkingConfig not needed, using legacy Links
		nil, // Platform not specified, use default
		image.GetContainerName())
	if err != nil {
		log.WithFields(log.Fields{
			"container":    image.GetContainerName(),
			"image":        image.GetFullyQualifiedName(),
			"env":          environmentVariables,
			"ports":        exposedPorts,
			"portBindings": portBindings,
			"error":        err,
		}).Fatal("Could not create container")
	}

	err = dockerClient.ContainerStart(
		context.Background(), containerCreationResponse.ID, containertypes.StartOptions{})
	if err == nil {
		log.WithFields(log.Fields{
			"container":    image.GetContainerName(),
			"image":        image.GetFullyQualifiedName(),
			"env":          environmentVariables,
			"ports":        exposedPorts,
			"portBindings": portBindings,
		}).Debug("Container has been started")
	}

	return err
}

// StartDockerContainer starts the stopped container
func StartDockerContainer(image liferay.Image) error {
	dockerClient := getDockerClient()

	var err error

	containers, err := PsFilterByLabel("lpn-type=" + image.GetType())

	if len(containers) == 0 {
		return errors.New("Error response from daemon: No such container: lpn-" + image.GetType())
	}

	for _, container := range containers {
		name := strings.TrimLeft(container.Names[0], "/")

		if name == image.GetContainerName() {
			// as we are using docker links for communications,
			// we need lpn instance to be started last
			continue
		}

		err = dockerClient.ContainerStart(
			context.Background(), name, containertypes.StartOptions{})
		if err == nil {
			log.WithFields(log.Fields{
				"container": name,
			}).Info("Database container has been started")
		}
	}

	err = dockerClient.ContainerStart(
		context.Background(), image.GetContainerName(), containertypes.StartOptions{})
	if err == nil {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
		}).Info("Container has been started")
	}

	return err
}

// StopDockerContainer stops the running container
func StopDockerContainer(image liferay.Image) error {
	dockerClient := getDockerClient()

	var err error

	containers, err := PsFilterByLabel("lpn-type=" + image.GetType())

	if len(containers) == 0 {
		return errors.New("Error response from daemon: No such container: lpn-" + image.GetType())
	}

	for _, container := range containers {
		name := strings.TrimLeft(container.Names[0], "/")
		err = dockerClient.ContainerStop(context.Background(), name, containertypes.StopOptions{})
		if err == nil {
			log.WithFields(log.Fields{
				"container": name,
			}).Info("Container has been stopped")
		}
	}

	return err
}

// ContainerInstance simple model for a container
type ContainerInstance struct {
	ID     string `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Status string `json:"status" binding:"required"`
}
