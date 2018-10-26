package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	types "github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	filters "github.com/docker/docker/api/types/filters"
	mount "github.com/docker/docker/api/types/mount"
	client "github.com/docker/docker/client"
	nat "github.com/docker/go-connections/nat"
	liferay "github.com/mdelapenya/lpn/liferay"
	shell "github.com/mdelapenya/lpn/shell"
)

// dockerBinary represents the name of the binary to execute Docker commands
const dockerBinary = "docker"

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

	reader, err := os.Open(path)
	if err != nil {
		return err
	}

	return dockerClient.CopyToContainer(
		context.Background(), image.GetContainerName(), image.GetDeployFolder(), reader,
		types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
}

func getDockerClient() *client.Client {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	return dockerClient
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

	cmdArgs := []string{
		"pull",
		dockerImage,
	}

	shell.StartAndWait(dockerBinary, cmdArgs)
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

		environmentVariables = append(environmentVariables, "DEBUG_MODE=true")
	}

	if memory != "" {
		environmentVariables = append(environmentVariables, "JVM_TUNING_MEMORY="+memory)
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

	dockerClient := getDockerClient()

	out, err := dockerClient.ImagePull(
		context.Background(), image.GetFullyQualifiedName(), types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	containerCreationResponse, err := dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        image.GetFullyQualifiedName(),
			Env:          environmentVariables,
			ExposedPorts: exposedPorts,
			Labels: map[string]string{
				"lpn": "",
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
