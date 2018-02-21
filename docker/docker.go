package docker

import (
	"fmt"
	"log"
	"os"

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

	return shell.RunCheck(dockerBinary, cmdArgs)
}

// CheckDockerImageExists checks if the image is already present
func CheckDockerImageExists(dockerImage string) bool {
	cmdArgs := []string{
		"image",
		"inspect",
		dockerImage,
	}

	result := shell.RunCheck(dockerBinary, cmdArgs)
	if result {
		log.Println("The image [" + dockerImage + "] has been pulled from Docker Hub.")
	} else {
		log.Println("The image [" + dockerImage + "] has NOT been pulled from Docker Hub.")
	}

	return result
}

// CopyFileToContainer copies a file to the running container
func CopyFileToContainer(image liferay.Image, path string) {
	if !CheckDockerContainerExists() {
		log.Println("The container [" + DockerContainerName + "] is NOT running.")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalln("The file [" + path + "] does NOT exist.")
		return
	}

	cmdArgs := []string{
		"cp",
		path,
		DockerContainerName + ":" + image.GetDeployFolder() + "/",
	}

	log.Println("Deploying [" + path + "] to " + image.GetDeployFolder())

	err := shell.Run(dockerBinary, cmdArgs)
	if err != nil {
		log.Fatal("Impossible to deploy the file to the container")
	}
}

// GetDockerImageFromRunningContainer gets the image name of the container
func GetDockerImageFromRunningContainer() string {
	cmdArgs := []string{
		"inspect",
		"--format='{{.Config.Image}}'",
		DockerContainerName,
	}

	return shell.Command(dockerBinary, cmdArgs)
}

// LogDockerContainer downloads the image
func LogDockerContainer() {
	if !CheckDockerContainerExists() {
		log.Println("The container [" + DockerContainerName + "] is NOT running.")
		return
	}

	cmdArgs := []string{
		"logs",
		"-f",
		DockerContainerName,
	}

	shell.StartAndWait(dockerBinary, cmdArgs)
}

// PullDockerImage downloads the image
func PullDockerImage(dockerImage string) {
	if CheckDockerImageExists(dockerImage) {
		log.Println("Skipping pulling [" + dockerImage + "] as it's already present locally.")
		return
	}

	log.Println("Pulling [" + dockerImage + "].")

	cmdArgs := []string{
		"pull",
		dockerImage,
	}

	shell.StartAndWait(dockerBinary, cmdArgs)
}

// RemoveDockerContainer removes the running container
func RemoveDockerContainer() {
	cmdArgs := []string{
		"rm",
		"-fv",
		DockerContainerName,
	}

	err := shell.CombinedOutput(dockerBinary, cmdArgs)
	if err != nil {
		log.Fatal("Impossible to remove the container")
	}
}

// RunDockerImage runs the image, setting the HTTP port for bundle and debug mode if needed
func RunDockerImage(dockerImage string, httpPort int, enableDebug bool, debugPort int) {
	PullDockerImage(dockerImage)

	if CheckDockerContainerExists() {
		RemoveDockerContainer()
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

	err := shell.CombinedOutput(dockerBinary, cmdArgs)
	if err != nil {
		log.Fatal("Impossible to run the container")
	}
}

// StopDockerContainer stops the running container
func StopDockerContainer() {
	cmdArgs := []string{
		"stop",
		DockerContainerName,
	}

	err := shell.CombinedOutput(dockerBinary, cmdArgs)
	if err != nil {
		log.Fatal("Impossible to stop the container")
	}
}
