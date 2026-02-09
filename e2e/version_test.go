package e2e

import (
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	lpnPath := getLpnPath(t)

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
