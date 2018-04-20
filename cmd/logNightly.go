package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	logCmd.AddCommand(logNightlyCmd)
}

var logNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Displays logs for the Liferay Portal Nightly Build instance",
	Long:  `Displays logs for the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		LogDockerContainer(nightly)
	},
}
