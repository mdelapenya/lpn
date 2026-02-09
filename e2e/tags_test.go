package e2e

import (
	"testing"
	"time"

	"github.com/ActiveState/termtest"
	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	lpnPath := getLpnPath(t)

	testCases := []string{"ce", "commerce", "dxp", "nightly", "release"}

	for _, imageType := range testCases {
		t.Run(imageType, func(t *testing.T) {
			opts := termtest.Options{
				CmdName:       lpnPath,
				Args:          []string{"tags", imageType},
				RetainWorkDir: false,
			}

			tt, err := termtest.New(opts)
			require.NoError(t, err)
			defer tt.Close()

			// Wait for tags output - should show some tag information
			// This command reaches out to Docker Hub, so it may take time
			_, err = tt.ExpectExitCode(0, 30*time.Second)
			require.NoError(t, err)
		})
	}
}
