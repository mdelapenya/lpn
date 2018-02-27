package docker

import (
	"fmt"
	"log"

	liferay "github.com/mdelapenya/lpn/liferay"
	shell "github.com/mdelapenya/lpn/shell"
)

// DockerContainerName represents the name of the container to be run
const DockerContainerName = "liferay-portal-nook"

// dockerBinary represents the name of the binary to execute Docker commands
const dockerBinary = "docker"

// CheckDockerContainerExists checks if the container is running
func CheckDockerContainerExists() bool {
	cmdArgs := []string{
		"container",
		"inspect",
		DockerContainerName,
	}

	err := shell.Run(dockerBinary, cmdArgs)
	if err != nil {
		log.Println("The container [" + DockerContainerName + "] is NOT running.")
		return false
	}

	log.Println("The container [" + DockerContainerName + "] is running.")
	return true
}

// CheckDockerImageExists checks if the image is already present
func CheckDockerImageExists(dockerImage string) error {
	cmdArgs := []string{
		"image",
		"inspect",
		dockerImage,
	}

	return shell.Run(dockerBinary, cmdArgs)
}

// CopyFileToContainer copies a file to the running container
func CopyFileToContainer(image liferay.Image, path string) error {
	cmdArgs := []string{
		"cp",
		path,
		DockerContainerName + ":" + image.GetDeployFolder() + "/",
	}

	log.Println("Deploying [" + path + "] to " + image.GetDeployFolder())

	err := shell.Run(dockerBinary, cmdArgs)
	if err != nil {
		return err
	}

	return nil
}

// GetDockerImageFromRunningContainer gets the image name of the container
func GetDockerImageFromRunningContainer() (string, error) {
	cmdArgs := []string{
		"inspect",
		"--format='{{.Config.Image}}'",
		DockerContainerName,
	}

	return shell.Command(dockerBinary, cmdArgs)
}

// LogDockerContainer downloads the image
func LogDockerContainer() {
	cmdArgs := []string{
		"logs",
		"-f",
		DockerContainerName,
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
func RemoveDockerContainer() error {
	cmdArgs := []string{
		"rm",
		"-fv",
		DockerContainerName,
	}

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}

// RunDockerImage runs the image, setting the HTTP port for bundle and debug mode if needed
func RunDockerImage(dockerImage string, httpPort int, enableDebug bool, debugPort int) error {
	if CheckDockerContainerExists() {
		_ = RemoveDockerContainer()
	}

	port := fmt.Sprintf("%d", httpPort)

	cmdArgs := []string{
		"run",
		"-d",
		"-p", port + ":8080",
		"--name", DockerContainerName,
	}

	if enableDebug {
		debugPortFlag := fmt.Sprintf("%d", debugPort) + ":9000"
		cmdArgs = append(cmdArgs, "-e", "DEBUG_MODE=true", "-p", debugPortFlag)
	}

	cmdArgs = append(cmdArgs, dockerImage)

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}

// StopDockerContainer stops the running container
func StopDockerContainer() error {
	cmdArgs := []string{
		"stop",
		DockerContainerName,
	}

	return shell.CombinedOutput(dockerBinary, cmdArgs)
}
