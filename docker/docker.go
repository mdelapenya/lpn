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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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

// CheckDocker checks if Docker is installed
func CheckDocker() bool {
	_, _, err := GetDockerVersion()
	if err != nil {
		return false
	}

	return true
}

// getContainersByLabel returns containers matching the lpn-container-name label
// Returns empty slice if no containers found or on error
func getContainersByLabel(containerName string) []types.Container {
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
		return []types.Container{}
	}

	return containers
}

// CheckDockerContainerExists checks if the container exists by label
// It looks for containers with the label "lpn-container-name" matching the provided name
func CheckDockerContainerExists(containerName string) bool {
	containers := getContainersByLabel(containerName)
	return len(containers) > 0
}

// GetContainerIDByLabel returns the container ID for a container with the given label
// Returns empty string if no container found
func GetContainerIDByLabel(containerName string) string {
	containers := getContainersByLabel(containerName)

	if len(containers) == 0 {
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

// getContainerWrapper creates a testcontainers.Container wrapper for an existing container
// This allows us to use testcontainers APIs on containers managed by lpn
func getContainerWrapper(containerName string) (testcontainers.Container, error) {
	ctx := context.Background()
	
	// Get provider
	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return nil, fmt.Errorf("could not create docker provider: %w", err)
	}

	// Get container ID
	containerID := GetContainerIDByLabel(containerName)
	if containerID == "" {
		return nil, fmt.Errorf("container not found: %s", containerName)
	}

	// Get container details
	dockerClient := getDockerClient()
	containerJSON, err := dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("could not inspect container: %w", err)
	}

	// Create container summary for ContainerFromType
	summary := containertypes.Summary{
		ID:     containerID,
		Image:  containerJSON.Config.Image,
		State:  containerJSON.State.Status,
		Labels: containerJSON.Config.Labels,
	}

	// Convert to testcontainers.Container
	container, err := provider.ContainerFromType(ctx, summary)
	if err != nil {
		return nil, fmt.Errorf("could not create container wrapper: %w", err)
	}

	return container, nil
}

// CopyFileToContainer copies a file to the running container using testcontainers API
func CopyFileToContainer(image liferay.Image, path string) error {
	ctx := context.Background()
	containerName := image.GetContainerName()
	deployFolder := image.GetDeployFolder()

	slog.Debug("Deploying file to target", "file", path, "target", deployFolder)

	// Get testcontainers.Container wrapper
	container, err := getContainerWrapper(containerName)
	if err != nil {
		slog.Error("Could not get container wrapper", "container", containerName, "error", err)
		return err
	}

	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		slog.Error("Could not open file to deploy", "file", path, "error", err)
		return err
	}

	// Use testcontainers API to copy file
	targetFilePath := filepath.Join(deployFolder, filepath.Base(path))
	err = container.CopyFileToContainer(ctx, path, targetFilePath, 0o777)
	if err != nil {
		slog.Error("Could not copy file to container", "container", containerName, "deployDir", deployFolder, "error", err)
		return err
	}

	// Change ownership of the deployed file
	owner := image.GetUser()
	cmd := []string{"chown", owner + ":" + owner, targetFilePath}
	
	_, _, err = container.Exec(ctx, cmd)
	if err != nil {
		slog.Error("Could not change file ownership", "container", containerName, "file", targetFilePath, "error", err)
		// Don't return error here as file was copied successfully
	}

	return nil
}

func getDockerClient() *client.Client {
	if instance != nil {
		return instance
	}

	ctx := context.Background()
	dockerClient, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		slog.Error("Could not get Docker client", "error", err)
		os.Exit(1)
	}

	// DockerClient embeds *client.Client, so we can use it directly
	instance = dockerClient.Client

	return instance
}

// GetDockerImageFromRunningContainer gets the image name of the container
func GetDockerImageFromRunningContainer(image liferay.Image) (string, error) {
	dockerClient := getDockerClient()

	containers, err := dockerClient.ContainerList(
		context.Background(), containertypes.ListOptions{All: true})

	if err != nil {
		slog.Error("Could not list all containers", "error", err)
		return "", err
	}

	for _, container := range containers {
		containerName := "/" + image.GetContainerName()

		if containerName == container.Names[0] {
			slog.Debug("Container found!", "container", image.GetContainerName())
			return container.Image, nil
		}
	}

	err = errors.New("We could not find the container among the running containers")
	slog.Debug("We could not find the container among the running containers", "container", image.GetContainerName(), "error", err)

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
		slog.Error("Container not found", "container", containerName)
		os.Exit(1)
	}

	containerJSON, err := dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		slog.Error("The container could not be inspected", "container", containerName, "error", err)
		os.Exit(1)
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
		slog.Error("Could not get container logs", "container", image.GetContainerName(), "error", err)
		os.Exit(1)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		slog.Error("Error following container logs", "container", image.GetContainerName(), "error", err)
		os.Exit(1)
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

	slog.Debug("Pulling Docker image.", "dockerImage", dockerImage)

	out, err := dockerClient.ImagePull(
		context.Background(), dockerImage, image.PullOptions{})

	if err == nil {
		parseImagePull(out)
	} else {
		slog.Error("The image could not be pulled", "dockerImage", dockerImage, "error", err)
		os.Exit(1)
	}
}

func parseImagePull(pullResp io.ReadCloser) {
	d := json.NewDecoder(pullResp)
	for {
		var pullResult imagePullResponse
		if err := d.Decode(&pullResult); err != nil {
			break
		}

		slog.Info("Image pull progress", "id", pullResult.ID, "status", pullResult.Status, "progress", pullResult.Progress)
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

		slog.Error("Could not filter container by label", "container", image.GetContainerName(), "label", label, "error", err)

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
			slog.Info("Container has been removed", "container", name)
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
		slog.Warn("Impossible to remove the image", "image", dockerImageName, "error", err)

		return err
	}

	slog.Info("Image has been removed", "image", dockerImageName)

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
		slog.Debug("Not starting a new container because it's already running", "container", containerName)

		return nil
	}

	// Create mount path for data persistence
	volumePath := filepath.Join(internal.LpnWorkspace, containerName)
	os.MkdirAll(volumePath, os.ModePerm)
	
	slog.Debug("Mounting database data folder", "container", containerName, "volume", volumePath)

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
			// Add labels for identification
			testcontainers.WithLabels(map[string]string{
				"lpn-container-name": containerName,
				"db-type":            image.GetType(),
				"lpn-type":           image.GetLpnType(),
			}),
			// Note: mysql.Run() has a built-in wait strategy that checks for MySQL readiness
			// No need to override it with a custom wait strategy
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
			// Add labels for identification
			testcontainers.WithLabels(map[string]string{
				"lpn-container-name": containerName,
				"db-type":            image.GetType(),
				"lpn-type":           image.GetLpnType(),
			}),
			// Note: postgres.Run() has a built-in wait strategy that checks for PostgreSQL readiness
			// No need to override it with a custom wait strategy
		)

	default:
		return fmt.Errorf("unsupported database type: %s", image.GetType())
	}

	if err != nil {
		slog.Error("Could not create database container", "container", containerName, "image", image.GetFullyQualifiedName(), "error", err)
		os.Exit(1)
	}

	slog.Debug("Database container has been started", "container", containerName, "image", image.GetFullyQualifiedName())

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
		slog.Debug("The container is running.", "container", image.GetContainerName())

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
			slog.Error("Non supported type", "imageType", imageType)
			os.Exit(1)
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
		slog.Error("Could not create container", "container", image.GetContainerName(), "image", image.GetFullyQualifiedName(), "env", environmentVariables, "ports", exposedPorts, "portBindings", portBindings, "error", err)
		os.Exit(1)
	}

	err = dockerClient.ContainerStart(
		context.Background(), containerCreationResponse.ID, containertypes.StartOptions{})
	if err == nil {
		slog.Debug("Container has been started", "container", image.GetContainerName(), "image", image.GetFullyQualifiedName(), "env", environmentVariables, "ports", exposedPorts, "portBindings", portBindings)
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
			slog.Info("Database container has been started", "container", name)
		}
	}

	err = dockerClient.ContainerStart(
		context.Background(), image.GetContainerName(), containertypes.StartOptions{})
	if err == nil {
		slog.Info("Container has been started", "container", image.GetContainerName())
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
			slog.Info("Container has been stopped", "container", name)
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
