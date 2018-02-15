package shell

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// CombinedOutput runs a shell program with arguments
func CombinedOutput(cmdBinary string, cmdArgs []string) error {
	cmd := exec.Command(cmdBinary, cmdArgs...)
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	fmt.Printf("%s\n", stdoutStderr)

	return nil
}

// Run runs a shell program with arguments
func Run(cmdBinary string, cmdArgs []string) bool {
	cmd := exec.Command(cmdBinary, cmdArgs...)

	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

// StartAndWait starts a program with arguments and waits for its output
func StartAndWait(cmdBinary string, cmdArgs []string) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(cmdBinary, cmdArgs...)

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
