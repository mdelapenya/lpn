package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)

	subcommands := []*cobra.Command{logCommerceCmd, logNightlyCmd, logReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		logCmd.AddCommand(subcommand)
	}
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays logs for the Liferay Portal instance",
	Long:  `Displays logs for the Liferay Portal instance, identified by [lpn] plus image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var logCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Displays logs for the Liferay Portal Commerce instance",
	Long:  `Displays logs for the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		logDockerContainer(commerce)
	},
}

var logNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Displays logs for the Liferay Portal Nightly Build instance",
	Long:  `Displays logs for the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		logDockerContainer(nightly)
	},
}

var logReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Displays logs for the Liferay Portal Release instance",
	Long:  `Displays logs for the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		logDockerContainer(release)
	},
}

// logDockerContainer show the logs for the running container of the specified type
func logDockerContainer(image liferay.Image) {
	docker.LogDockerContainer(image)
}
