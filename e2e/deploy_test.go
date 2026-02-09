package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

// TestDeployCommandWithoutContainer tests deploy command fails gracefully when no container is running
func TestDeployCommandWithoutContainer(t *testing.T) {
	lpnPath := getLpnPath(t)

	// Create a temporary test file to deploy
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-deploy.jar")
	err := os.WriteFile(testFile, []byte("test jar content"), 0644)
	require.NoError(t, err)

	// Test all Liferay flavors
	flavors := []string{"ce", "commerce", "dxp", "nightly", "release"}

	for _, flavor := range flavors {
		t.Run(flavor, func(t *testing.T) {
			opts := termtest.Options{
				CmdName:       lpnPath,
				Args:          []string{"deploy", flavor, "-f", testFile},
				RetainWorkDir: false,
			}

			tt, err := termtest.New(opts)
			require.NoError(t, err)
			defer tt.Close()

			// The command should fail because no container is running
			// We expect an error message about container not found or not running
			_, err = tt.ExpectExitCode(1, 10*time.Second)
			require.NoError(t, err, "Deploy should fail when container is not running")
		})
	}
}

// TestDeployCommandWithoutFile tests deploy command validation
func TestDeployCommandWithoutFile(t *testing.T) {
	lpnPath := getLpnPath(t)

	opts := termtest.Options{
		CmdName:       lpnPath,
		Args:          []string{"deploy", "ce"},
		RetainWorkDir: false,
	}

	tt, err := termtest.New(opts)
	require.NoError(t, err)
	defer tt.Close()

	// Expect error message about missing file or directory
	_, err = tt.Expect("Please pass a valid path to a file or to a directory as argument", 5*time.Second)
	require.NoError(t, err)

	// Verify exit code
	_, err = tt.ExpectExitCode(1, 5*time.Second)
	require.NoError(t, err)
}

// TestDeployCommandWithNonExistentFile tests deploy command with a non-existent file
func TestDeployCommandWithNonExistentFile(t *testing.T) {
	lpnPath := getLpnPath(t)

	opts := termtest.Options{
		CmdName:       lpnPath,
		Args:          []string{"deploy", "ce", "-f", "/tmp/nonexistent-file-12345.jar"},
		RetainWorkDir: false,
	}

	tt, err := termtest.New(opts)
	require.NoError(t, err)
	defer tt.Close()

	// The command should fail
	_, err = tt.ExpectExitCode(1, 10*time.Second)
	require.NoError(t, err, "Deploy should fail with non-existent file")
}

// TestDeployCommandHelp tests the deploy command help output
func TestDeployCommandHelp(t *testing.T) {
	lpnPath := getLpnPath(t)

	opts := termtest.Options{
		CmdName:       lpnPath,
		Args:          []string{"deploy", "--help"},
		RetainWorkDir: false,
	}

	tt, err := termtest.New(opts)
	require.NoError(t, err)
	defer tt.Close()

	// Check for expected help text
	_, err = tt.Expect("Deploys files or a directory to Liferay Portal's deploy folder", 5*time.Second)
	require.NoError(t, err)

	// Check for subcommands
	_, err = tt.Expect("Available Commands:", 5*time.Second)
	require.NoError(t, err)

	// Verify exit code
	_, err = tt.ExpectExitCode(0, 5*time.Second)
	require.NoError(t, err)
}

// TestDeployCommandSubcommands tests that all deploy subcommands are available
func TestDeployCommandSubcommands(t *testing.T) {
	lpnPath := getLpnPath(t)

	subcommands := []string{"ce", "commerce", "dxp", "nightly", "release"}

	for _, subcommand := range subcommands {
		t.Run(subcommand, func(t *testing.T) {
			opts := termtest.Options{
				CmdName:       lpnPath,
				Args:          []string{"deploy", subcommand, "--help"},
				RetainWorkDir: false,
			}

			tt, err := termtest.New(opts)
			require.NoError(t, err)
			defer tt.Close()

			// Check for expected help text
			_, err = tt.Expect("Deploys files or a directory", 5*time.Second)
			require.NoError(t, err)

			// Verify exit code
			_, err = tt.ExpectExitCode(0, 5*time.Second)
			require.NoError(t, err)
		})
	}
}
