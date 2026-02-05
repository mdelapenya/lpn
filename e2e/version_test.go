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

func TestVersion(t *testing.T) {
	// Try to find lpn in PATH first
	lpnPath, err := exec.LookPath("lpn")
	if err != nil {
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
		
		// Get absolute path
		lpnPath, err = filepath.Abs(lpnPath)
		require.NoError(t, err)
	}

	// Create a new terminal session
	opts := termtest.Options{
		CmdName:       lpnPath,
		Args:          []string{"version"},
		RetainWorkDir: false,
	}

	tt, err := termtest.New(opts)
	require.NoError(t, err)
	defer tt.Close()

	// Wait for expected output
	_, err = tt.Expect("0.14.0", 5*time.Second)
	require.NoError(t, err)
	
	_, err = tt.Expect("dockerClient=", 5*time.Second)
	require.NoError(t, err)
	
	_, err = tt.Expect("dockerServer=", 5*time.Second)
	require.NoError(t, err)
	
	_, err = tt.Expect("golang=", 5*time.Second)
	require.NoError(t, err)
	
	_, err = tt.Expect("lpn (Liferay Portal Nook) v", 5*time.Second)
	require.NoError(t, err)

	// Wait for command to complete and verify exit code
	_, err = tt.ExpectExitCode(0, 5*time.Second)
	require.NoError(t, err)
}
