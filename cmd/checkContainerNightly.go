package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	checkCmd.AddCommand(checkContainerNightlyCmd)
}

var checkContainerNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Checks if there is a Nightly Build container created by lpn",
	Long: `Checks if there is a Nightly Build container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Nightly Build container with name [lpn-nightly] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		CheckDockerContainerExists(nightly)
	},
}
