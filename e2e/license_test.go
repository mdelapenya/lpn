package e2e

import (
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

func TestLicense(t *testing.T) {
	lpnPath := getLpnPath(t)

	opts := termtest.Options{
		CmdName:       lpnPath,
		Args:          []string{"license"},
		RetainWorkDir: false,
	}

	tt, err := termtest.New(opts)
	require.NoError(t, err)
	defer tt.Close()

	// Expect license information
	_, err = tt.Expect("GNU Lesser General Public License", 5*time.Second)
	require.NoError(t, err)

	// Verify exit code
	_, err = tt.ExpectExitCode(0, 5*time.Second)
	require.NoError(t, err)
}
