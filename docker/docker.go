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
	"log"
	"os"
	"path/filepath"

	types "github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	filters "github.com/docker/docker/api/types/filters"
	mount "github.com/docker/docker/api/types/mount"
	client "github.com/docker/docker/client"
	nat "github.com/docker/go-connections/nat"
	liferay "github.com/mdelapenya/lpn/liferay"
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
		return bytes.Buffer{}, fmt.Errorf("Could not build TAR header: %v", err)
	}

	b, err := ioutil.ReadFile(file.Name())
	tarWriter.Write(b)
	defer tarWriter.Close()

	return buffer, nil
}

// CheckDocker checks if Docker is installed
func CheckDocker() bool {
	_, err := GetDockerVersion()
	if err != nil {
		return false
	}

	return true
}

// CheckDockerContainerExists checks if the container is running
func CheckDockerContainerExists(image liferay.Image) bool {
	dockerClient := getDockerClient()

	containers, err := dockerClient.ContainerList(
		context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		return false
	}

	for _, container := range containers {
		containerName := "/" + image.GetContainerName()

		if containerName == container.Names[0] {
			return true
		}
	}

	return false
}

// CheckDockerImageExists checks if the image is already present
func CheckDockerImageExists(dockerImage string) bool {
	dockerClient := getDockerClient()

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

	log.Println("Deploying [" + path + "] to " + image.GetDeployFolder())

	_, err := dockerClient.ContainerStatPath(
		context.Background(), image.GetContainerName(), image.GetDeployFolder())
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer, err := buildTarForDeployment(file)
	if err != nil {
		return err
	}

	err = dockerClient.CopyToContainer(
		context.Background(), image.GetContainerName(), image.GetDeployFolder(),
		&buffer, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})

	if err == nil {
		targetFilePath := filepath.Join(image.GetDeployFolder(), filepath.Base(file.Name()))
		owner := image.GetUser()

		cmd := []string{"chown", owner + ":" + owner, targetFilePath}

		execCommandIntoContainer(image.GetContainerName(), cmd)
	}

	return err
}

func execCommandIntoContainer(containerName string, cmd []string) error {
	dockerClient := getDockerClient()

	response, err := dockerClient.ContainerExecCreate(
		context.Background(), containerName, types.ExecConfig{
			User:         "root",
			Tty:          false,
			AttachStdin:  false,
			AttachStderr: false,
			AttachStdout: false,
			Detach:       true,
			Cmd:          cmd,
		})

	if err != nil {
		return err
	}

	err = dockerClient.ContainerExecStart(
		context.Background(), response.ID, types.ExecStartCheck{
			Detach: true,
			Tty:    false,
		})

	return err
}

func getDockerClient() *client.Client {
	if instance != nil {
		return instance
	}

	instance, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	return instance
}

// GetDockerImageFromRunningContainer gets the image name of the container
func GetDockerImageFromRunningContainer(image liferay.Image) (string, error) {
	dockerClient := getDockerClient()

	containers, err := dockerClient.ContainerList(
		context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		return "", err
	}

	for _, container := range containers {
		containerName := "/" + image.GetContainerName()

		if containerName == container.Names[0] {
			return container.Image, nil
		}
	}

	return "", errors.New("We could not find the container among the running containers")
}

// GetDockerVersion returns the output of Docker version
func GetDockerVersion() (string, error) {
	dockerClient := getDockerClient()

	serverVersion, err := dockerClient.ServerVersion(context.Background())

	version := "Client version: " + dockerClient.ClientVersion() + "\n"
	version += "Server version: " + serverVersion.Version + "\n"
	version += "Go version: " + serverVersion.GoVersion

	return version, err
}

// inspect inspects a container
func inspect(containerName string) types.ContainerJSON {
	dockerClient := getDockerClient()

	containerJSON, err := dockerClient.ContainerInspect(context.Background(), containerName)
	if err != nil {
		log.Fatalln("The container [" + containerName + "] could not be inspected")
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
		types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}

// PsFilter Retrieves all containers with a label
func PsFilter(label string) ([]types.Container, error) {
	dockerClient := getDockerClient()

	filters := filters.NewArgs()
	filters.Add("label", label)

	return dockerClient.ContainerList(
		context.Background(), types.ContainerListOptions{
			Size:    true,
			All:     true,
			Since:   "container",
			Filters: filters,
		})
}

// PullDockerImage downloads the image
func PullDockerImage(dockerImage string) {
	log.Println("Pulling [" + dockerImage + "].")

	dockerClient := getDockerClient()

	out, err := dockerClient.ImagePull(
		context.Background(), dockerImage, types.ImagePullOptions{})

	if err == nil {
		parseImagePull(out)
	} else {
		log.Fatalf("The image [" + dockerImage + "] could not be pulled")
	}
}

func parseImagePull(pullResp io.ReadCloser) {
	d := json.NewDecoder(pullResp)
	for {
		var pullResult imagePullResponse
		if err := d.Decode(&pullResult); err != nil {
			break
		}

		fmt.Printf("%s %s %s\n", pullResult.ID, pullResult.Status, pullResult.Progress)
	}
}

// RemoveDockerContainer removes a running container
func RemoveDockerContainer(containerName string) error {
	dockerClient := getDockerClient()

	return dockerClient.ContainerRemove(
		context.Background(), containerName, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
}

// RemoveDockerImage removes a docker image
func RemoveDockerImage(dockerImageName string) error {
	dockerClient := getDockerClient()

	_, err := dockerClient.ImageRemove(
		context.Background(), dockerImageName,
		types.ImageRemoveOptions{
			Force: true,
		})

	if err == nil {
		log.Println("[" + dockerImageName + "] was deleted.")
	}

	return err
}

// RunDockerImage runs the image, setting the HTTP and GoGoShell ports for bundle, debug mode, and
// jvmMemory if needed
func RunDockerImage(
	image liferay.Image, httpPort int, gogoShellPort int, enableDebug bool, debugPort int,
	memory string, properties string) error {

	if CheckDockerContainerExists(image) {
		log.Println("The container [" + image.GetContainerName() + "] is running.")
		_ = RemoveDockerContainer(image.GetContainerName())
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
		jvmEnvVarName := ""

		switch imageType := image.(type) {
		case liferay.CE, liferay.Commerce, liferay.DXP, liferay.Nightly:
			jvmEnvVarName = "LIFERAY_JVM_OPTS"
		case liferay.Release:
			jvmEnvVarName = "JVM_TUNING_MEMORY"
		default:
			log.Fatalln("Non supported type", imageType)
		}

		environmentVariables = append(environmentVariables, jvmEnvVarName+"="+memory)
	}

	var mounts []mount.Mount

	if properties != "" {
		log.Println("Mounting " + properties + " as configuration file")

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: properties,
			Target: image.GetLiferayHome() + "/portal-ext.properties",
		})
	}

	PullDockerImage(image.GetFullyQualifiedName())

	dockerClient := getDockerClient()

	containerCreationResponse, err := dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        image.GetFullyQualifiedName(),
			Env:          environmentVariables,
			ExposedPorts: exposedPorts,
			Labels: map[string]string{
				"lpn-type": image.GetType(),
			},
		},
		&container.HostConfig{
			PortBindings: portBindings,
			Mounts:       mounts,
		},
		nil, image.GetContainerName())
	if err != nil {
		panic(err)
	}

	return dockerClient.ContainerStart(
		context.Background(), containerCreationResponse.ID, types.ContainerStartOptions{})
}

// StartDockerContainer stops the running container
func StartDockerContainer(containerName string) error {
	dockerClient := getDockerClient()

	return dockerClient.ContainerStart(
		context.Background(), containerName, types.ContainerStartOptions{})
}

// StopDockerContainer stops the running container
func StopDockerContainer(containerName string) error {
	dockerClient := getDockerClient()

	return dockerClient.ContainerStop(context.Background(), containerName, nil)
}

// ContainerInstance simple model for a container
type ContainerInstance struct {
	ID     string `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Status string `json:"status" binding:"required"`
}
