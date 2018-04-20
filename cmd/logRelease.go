package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	logCmd.AddCommand(logReleaseCmd)
}

var logReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Displays logs for the Liferay Portal Release instance",
	Long:  `Displays logs for the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		LogDockerContainer(release)
	},
}
