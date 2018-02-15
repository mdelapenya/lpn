package docker

import (
	"log"
	shell "lpn/shell"
)

// DockerImage represents the base namespace for the Docker image
const DockerImage = "mdelapenya/liferay-portal-nightlies"

// DockerContainerName represents the name of the container to be run
const DockerContainerName = "liferay-portal-nightly"

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

	return shell.RunCheck(dockerBinary, cmdArgs)
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

// RunDockerImage runs the image
func RunDockerImage(dockerImage string) {
	PullDockerImage(dockerImage)

	if CheckDockerContainerExists() {
		RemoveDockerContainer()
	}

	cmdArgs := []string{
		"run",
		"-d",
		"-p", "8080:8080",
		"-p", "11311:11311",
		"--name", DockerContainerName,
		dockerImage,
	}

	err := shell.CombinedOutput(dockerBinary, cmdArgs)
	if err != nil {
		log.Fatal("Impossible to run the container")
	}
}
