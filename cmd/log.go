package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)

	subcommands := []*cobra.Command{
		logCECmd, logCommerceCmd, logDXPCmd, logNightlyCmd, logReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		logCmd.AddCommand(subcommand)

		subcommand.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs commands with Debug log level")
		subcommand.VisitParents(addVerboseFlag)
	}
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays logs for the Liferay Portal instance",
	Long:  `Displays logs for the Liferay Portal instance, identified by [lpn] plus image type.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		enableDebugLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var logCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Displays logs for the Liferay Portal CE instance",
	Long:  `Displays logs for the Liferay Portal CE instance, identified by [lpn-ce].`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		logContainer(ce)
	},
}

var logCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Displays logs for the Liferay Portal Commerce instance",
	Long:  `Displays logs for the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		logContainer(commerce)
	},
}

var logDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Displays logs for the Liferay DXP instance",
	Long:  `Displays logs for the Liferay DXP instance, identified by [lpn-dxp].`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		logContainer(dxp)
	},
}

var logNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Displays logs for the Liferay Portal Nightly Build instance",
	Long:  `Displays logs for the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		logContainer(nightly)
	},
}

var logReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Displays logs for the Liferay Portal Release instance",
	Long:  `Displays logs for the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		logContainer(release)
	},
}

// logContainer show the logs for the running container of the specified type
func logContainer(image liferay.Image) {
	docker.LogContainer(image)
}
