package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)

	subcommands := []*cobra.Command{
		stopCECmd, stopCommerceCmd, stopDXPCmd, stopNightlyCmd, stopReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		stopCmd.AddCommand(subcommand)
	}
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the Liferay Portal nook instance",
	Long:  `Stops the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var stopCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Stops the Liferay Portal CE instance",
	Long:  `Stops the Liferay Portal CE instance, identified by [lpn-cd].`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		stopDockerContainer(ce)
	},
}

var stopCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Stops the Liferay Portal Commerce instance",
	Long:  `Stops the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		stopDockerContainer(commerce)
	},
}

var stopDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Stops the Liferay DXP instance",
	Long:  `Stops the Liferay DXP instance, identified by [lpn-dxp].`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		stopDockerContainer(dxp)
	},
}

var stopNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Stops the Liferay Portal Nightly Build instance",
	Long:  `Stops the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		stopDockerContainer(nightly)
	},
}

var stopReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Stops the Liferay Portal Release instance",
	Long:  `Stops the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		stopDockerContainer(release)
	},
}

// stopDockerContainer stops the running container
func stopDockerContainer(image liferay.Image) {
	err := docker.StopDockerContainer(image.GetContainerName())
	if err != nil {
		log.Fatalln("Impossible to stop the container ["+image.GetContainerName()+"]", err)
	}

	log.Println("[" + image.GetContainerName() + "] stopped")
}
