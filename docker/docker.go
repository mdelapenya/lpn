package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	types "github.com/docker/docker/api/types"
	filters "github.com/docker/docker/api/types/filters"
	client "github.com/docker/docker/client"
	liferay "github.com/mdelapenya/lpn/liferay"
	shell "github.com/mdelapenya/lpn/shell"
)

// dockerBinary represents the name of the binary to execute Docker commands
const dockerBinary = "docker"

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

// LogDockerContainer downloads the image
func LogDockerContainer(image liferay.Image) {
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

	cmdArgs := []string{
		"run",
		"-d",
		"--label", "lpn",
		"-p", port + ":8080",
		"-p", gogoPort + ":11311",
		"--name", image.GetContainerName(),
	}

	if enableDebug {
		debugPortFlag := fmt.Sprintf("%d", debugPort) + ":9000"
		cmdArgs = append(cmdArgs, "-e", "DEBUG_MODE=true", "-p", debugPortFlag)
	}

	if memory != "" {
		jvmMemory := "JVM_TUNING_MEMORY=" + memory
		cmdArgs = append(cmdArgs, "-e", jvmMemory)
	}

	if properties != "" {
		log.Println("Mounting " + properties + " as configuration file")

		portalProperties := filepath.FromSlash(
			properties + ":" + image.GetLiferayHome() + "/portal-ext.properties")

		cmdArgs = append(cmdArgs, "-v", portalProperties)
	}

	cmdArgs = append(cmdArgs, image.GetFullyQualifiedName())

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}

// StopDockerContainer stops the running container
func StopDockerContainer(containerName string) error {
	cmdArgs := []string{
		"stop",
		containerName,
	}

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}

// ContainerInstance simple model for a container
type ContainerInstance struct {
	ID     string `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Status string `json:"status" binding:"required"`
}
