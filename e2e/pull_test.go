package e2e

import (
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

// TestPullCommandHelp tests the pull command help output
func TestPullCommandHelp(t *testing.T) {
	lpnPath := getLpnPath(t)

	opts := termtest.Options{
		CmdName:       lpnPath,
		Args:          []string{"pull", "--help"},
		RetainWorkDir: false,
	}

	tt, err := termtest.New(opts)
	require.NoError(t, err)
	defer tt.Close()

	// Check for expected help text
	_, err = tt.Expect("Pulls a Liferay Portal Docker image", 5*time.Second)
	require.NoError(t, err)

	// Check for subcommands
	_, err = tt.Expect("Available Commands:", 5*time.Second)
	require.NoError(t, err)

	// Verify exit code
	_, err = tt.ExpectExitCode(0, 5*time.Second)
	require.NoError(t, err)
}

// TestPullCommandSubcommands tests that all pull subcommands are available
func TestPullCommandSubcommands(t *testing.T) {
	lpnPath := getLpnPath(t)

	subcommands := []string{"ce", "commerce", "dxp", "nightly", "release"}

	for _, subcommand := range subcommands {
		t.Run(subcommand, func(t *testing.T) {
			opts := termtest.Options{
				CmdName:       lpnPath,
				Args:          []string{"pull", subcommand, "--help"},
				RetainWorkDir: false,
			}

			tt, err := termtest.New(opts)
			require.NoError(t, err)
			defer tt.Close()

			// Check for expected help text
			_, err = tt.Expect("Pulls a Liferay", 5*time.Second)
			require.NoError(t, err)

			// Verify exit code
			_, err = tt.ExpectExitCode(0, 5*time.Second)
			require.NoError(t, err)
		})
	}
}

// TestPullNonExistentImage tests pulling a non-existent image tag
func TestPullNonExistentImage(t *testing.T) {
	lpnPath := getLpnPath(t)

	testCases := []struct {
		imageType  string
		repository string
	}{
		{"ce", "liferay/portal"},
		{"commerce", "liferay/commerce"},
		{"dxp", "liferay/dxp"},
		{"nightly", "mdelapenya/portal-snapshot"},
		{"release", "mdelapenya/liferay-portal"},
	}

	for _, tc := range testCases {
		t.Run(tc.imageType, func(t *testing.T) {
			opts := termtest.Options{
				CmdName:       lpnPath,
				Args:          []string{"pull", tc.imageType, "-t", "nonexistent-tag-12345"},
				RetainWorkDir: false,
			}

			tt, err := termtest.New(opts)
			require.NoError(t, err)
			defer tt.Close()

			// Expect error message about image not being pulled
			_, err = tt.Expect("The image could not be pulled", 30*time.Second)
			require.NoError(t, err)

			// Expect the docker image reference in the error output
			_, err = tt.Expect(tc.repository, 5*time.Second)
			require.NoError(t, err)

			// Verify exit code is non-zero (failure)
			_, err = tt.ExpectExitCode(1, 5*time.Second)
			require.NoError(t, err)
		})
	}
}
