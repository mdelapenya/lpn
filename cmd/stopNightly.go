package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	stopCmd.AddCommand(stopNightlyCmd)
}

var stopNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Stops the Liferay Portal Nightly Build instance",
	Long:  `Stops the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		StopDockerContainer(nightly)
	},
}
