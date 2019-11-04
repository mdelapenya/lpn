package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)

	subcommands := []*cobra.Command{
		startCECmd, startCommerceCmd, startDXPCmd, startNightlyCmd, startReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		startCmd.AddCommand(subcommand)

		subcommand.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs commands with Debug log level")
		subcommand.VisitParents(addVerboseFlag)
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Liferay Portal nook instance",
	Long:  `Starts the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		enableDebugLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var startCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Starts the Liferay Portal CE instance",
	Long:  `Starts the Liferay Portal CE instance, identified by [lpn-ce].`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		startDockerContainer(ce)
	},
}

var startCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Starts the Liferay Portal Commerce instance",
	Long:  `Starts the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		startDockerContainer(commerce)
	},
}

var startDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Starts the Liferay DXP instance",
	Long:  `Starts the Liferay DXP instance, identified by [lpn-dxp].`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		startDockerContainer(dxp)
	},
}

var startNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Starts the Liferay Portal Nightly Build instance",
	Long:  `Starts the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		startDockerContainer(nightly)
	},
}

var startReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Starts the Liferay Portal Release instance",
	Long:  `Starts the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		startDockerContainer(release)
	},
}

// startDockerContainer starts the stopped container
func startDockerContainer(image liferay.Image) {
	err := docker.StartDockerContainer(image)
	if err != nil {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
			"error":     err,
		}).Warn("Impossible to start the container")
	}
}
