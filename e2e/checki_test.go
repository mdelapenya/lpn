package e2e

import (
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

func TestCheckImageWhenNotPresent(t *testing.T) {
	lpnPath := getLpnPath(t)

	testCases := []string{"ce", "commerce", "dxp", "nightly", "release"}

	for _, imageType := range testCases {
		t.Run(imageType, func(t *testing.T) {
			opts := termtest.Options{
				CmdName:       lpnPath,
				Args:          []string{"checki", imageType},
				RetainWorkDir: false,
			}

			tt, err := termtest.New(opts)
			require.NoError(t, err)
			defer tt.Close()

			// The command may show the image doesn't exist or list images
			// Verify exit code
			_, err = tt.ExpectExitCode(0, 10*time.Second)
			require.NoError(t, err)
		})
	}
}
