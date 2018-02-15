package shell

import "os/exec"

// Run runs a shell program with arguments
func Run(cmdBinary string, cmdArgs []string) bool {
	cmd := exec.Command(cmdBinary, cmdArgs...)

	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
