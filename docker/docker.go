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

// DownloadDockerImage downloads the image
func downloadDockerImage(dockerImage string) {
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

// RunDockerImage runs the image
func RunDockerImage(dockerImage string) {
	cmd := exec.Command("docker", "run", "-d", dockerImage)
	cmd.Run()
}
