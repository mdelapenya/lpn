package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

func getLpnPath(t *testing.T) string {
	// Try to find lpn in PATH first
	lpnPath, err := exec.LookPath("lpn")
	if err == nil {
		return lpnPath
	}

	// Otherwise look for it in ../bin/lpn relative to current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)
	
	// If we're in the e2e directory, go up one level
	if filepath.Base(cwd) == "e2e" {
		lpnPath = filepath.Join(filepath.Dir(cwd), "bin", "lpn")
	} else {
		lpnPath = filepath.Join(cwd, "bin", "lpn")
	}
	
	if _, err := os.Stat(lpnPath); os.IsNotExist(err) {
		t.Skip("lpn binary not found, run 'go build -o ./bin/lpn' first")
	}
	
	// Return absolute path
	absPath, err := filepath.Abs(lpnPath)
	require.NoError(t, err)
	return absPath
}

func TestCheckContainerNotRunning(t *testing.T) {
	lpnPath := getLpnPath(t)

	testCases := []string{"ce", "commerce", "dxp", "nightly", "release"}

	for _, containerType := range testCases {
		t.Run(containerType, func(t *testing.T) {
			opts := termtest.Options{
				CmdName:       lpnPath,
				Args:          []string{"checkc", containerType},
				RetainWorkDir: false,
			}

			tt, err := termtest.New(opts)
			require.NoError(t, err)
			defer tt.Close()

			// Expect container does NOT exist
			_, err = tt.Expect("Container does NOT exist in the system", 10*time.Second)
			require.NoError(t, err)
			
			_, err = tt.Expect("container=lpn-"+containerType, 5*time.Second)
			require.NoError(t, err)

			// Verify exit code
			_, err = tt.ExpectExitCode(0, 5*time.Second)
			require.NoError(t, err)
		})
	}
}
