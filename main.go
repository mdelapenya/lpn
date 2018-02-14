package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

const dockerImage = "mdelapenya/liferay-portal-nightlies"

func main() {
	fmt.Print("Enter the Image Tag you want to use for [" + dockerImage + "]: ")
	var imageTag string

	fmt.Scanf("%s", &imageTag)

	downloadDockerImage(getDockerImage(imageTag))
}

func getDockerImage(imageTag string) string {
	return dockerImage + ":" + imageTag
}

// CapturingPassThroughWriter is a writer that remembers
// data written to it and passes it to w
type CapturingPassThroughWriter struct {
	buf bytes.Buffer
	w   io.Writer
}

// NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
func NewCapturingPassThroughWriter(w io.Writer) *CapturingPassThroughWriter {
	return &CapturingPassThroughWriter{
		w: w,
	}
}

func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
	w.buf.Write(d)
	return w.w.Write(d)
}

// Bytes returns bytes written to the writer
func (w *CapturingPassThroughWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func downloadDockerImage(dockerImage string) {
	var errStdout, errStderr error
	cmd := exec.Command("docker", "pull", dockerImage)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := NewCapturingPassThroughWriter(os.Stdout)
	stderr := NewCapturingPassThroughWriter(os.Stderr)

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
		log.Fatalf("failed to capture stdout or stderr\n")
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
