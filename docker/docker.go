package docker

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// DockerImage represents the base namespace for the Docker image
const DockerImage = "mdelapenya/liferay-portal-nightlies"

// dockerContainerName represents the name of the container to be run
const dockerContainerName = "liferay-portal-nightly"

// checkDockerContainerExists checks if the container is running
func checkDockerContainerExists() bool {
	cmdName := "docker"
	cmdArgs := []string{"container", "inspect", dockerContainerName}

	cmd := exec.Command(cmdName, cmdArgs...)

	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

// checkDockerImageExists checks if the image is already present
func checkDockerImageExists(dockerImage string) bool {
	cmdName := "docker"
	cmdArgs := []string{"image", "inspect", dockerImage}

	cmd := exec.Command(cmdName, cmdArgs...)

	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

// DownloadDockerImage downloads the image
func downloadDockerImage(dockerImage string) {
	if checkDockerImageExists(dockerImage) {
		log.Println("Skipping pulling [" + dockerImage + "] as it's already present locally.")
		return
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("docker", "pull", dockerImage)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("%s", errStr)
	fmt.Printf("%s", outStr)
}

// removeDockerContainer removes the running container
func removeDockerContainer() {
	cmd := exec.Command("docker", "rm", "-fv", dockerContainerName)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", stdoutStderr)
}

// RunDockerImage runs the image
func RunDockerImage(dockerImage string) {
	downloadDockerImage(dockerImage)

	if checkDockerContainerExists() {
		removeDockerContainer()
	}

	cmd := exec.Command("docker", "run", "-d", "--name", dockerContainerName, dockerImage)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", stdoutStderr)
}
