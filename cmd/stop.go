package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.AddCommand(stopCommerceCmd)
	stopCmd.AddCommand(stopNightlyCmd)
	stopCmd.AddCommand(stopReleaseCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the Liferay Portal nook instance",
	Long:  `Stops the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
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
	err := docker.StopDockerContainer(image)
	if err != nil {
		log.Fatalln("Impossible to stop the container [" + image.GetContainerName() + "]")
	}
}
