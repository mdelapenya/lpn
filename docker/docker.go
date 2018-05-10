package docker

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	types "github.com/docker/docker/api/types"
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
		panic(err)
	}

	for _, container := range containers {
		if image.GetContainerName() == container.Names[0] {
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
	cmdArgs := []string{
		"cp",
		path,
		image.GetContainerName() + ":" + image.GetDeployFolder() + "/",
	}

	log.Println("Deploying [" + path + "] to " + image.GetDeployFolder())

	err := shell.Run(dockerBinary, cmdArgs)
	if err != nil {
		return err
	}

	return nil
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
	cmdArgs := []string{
		"inspect",
		"--format='{{.Config.Image}}'",
		image.GetContainerName(),
	}

	return shell.Command(dockerBinary, cmdArgs)
}

// GetDockerVersion returns the output of Docker version
func GetDockerVersion() (string, error) {
	cmdArgs := []string{
		"version",
	}

	return shell.Command(dockerBinary, cmdArgs)
}

// LogDockerContainer downloads the image
func LogDockerContainer(image liferay.Image) {
	cmdArgs := []string{
		"logs",
		"-f",
		image.GetContainerName(),
	}

	shell.StartAndWait(dockerBinary, cmdArgs)
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

// RemoveDockerContainer removes the running container
func RemoveDockerContainer(image liferay.Image) error {
	cmdArgs := []string{
		"rm",
		"-fv",
		image.GetContainerName(),
	}

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}

// RemoveDockerImage removes a docker image
func RemoveDockerImage(dockerImageName string) error {
	cmdArgs := []string{
		"rmi",
		dockerImageName,
	}

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}

// RunDockerImage runs the image, setting the HTTP and GoGoShell ports for bundle, debug mode, and
// jvmMemory if needed
func RunDockerImage(
	image liferay.Image, httpPort int, gogoShellPort int, enableDebug bool, debugPort int,
	memory string, properties string) error {

	if CheckDockerContainerExists(image) {
		log.Println("The container [" + image.GetContainerName() + "] is running.")
		_ = RemoveDockerContainer(image)
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
func StopDockerContainer(image liferay.Image) error {
	cmdArgs := []string{
		"stop",
		image.GetContainerName(),
	}

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}
