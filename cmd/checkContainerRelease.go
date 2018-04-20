package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	checkCmd.AddCommand(checkContainerReleaseCmd)
}

var checkContainerReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Checks if there is a Release container created by lpn",
	Long: `Checks if there is a Release container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Release container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		CheckDockerContainerExists(release)
	},
}
