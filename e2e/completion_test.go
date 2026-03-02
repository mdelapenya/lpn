package e2e

import (
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

func TestCompletion(t *testing.T) {
	lpnPath := getLpnPath(t)

	opts := termtest.Options{
		CmdName:       lpnPath,
		Args:          []string{"completion"},
		RetainWorkDir: false,
	}

	tt, err := termtest.New(opts)
	require.NoError(t, err)
	defer tt.Close()

	// Expect ZSH completion output
	_, err = tt.Expect("#compdef", 5*time.Second)
	require.NoError(t, err)

	// Verify exit code
	_, err = tt.ExpectExitCode(0, 5*time.Second)
	require.NoError(t, err)
}
