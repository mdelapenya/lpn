package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rmCmd.AddCommand(rmReleaseCmd)
}

var rmReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Removes the Liferay Portal Release instance",
	Long:  `Removes the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		RemoveDockerContainer(release)
	},
}
