package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	stopCmd.AddCommand(stopReleaseCmd)
}

var stopReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Stops the Liferay Portal Release instance",
	Long:  `Stops the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		StopDockerContainer(release)
	},
}
