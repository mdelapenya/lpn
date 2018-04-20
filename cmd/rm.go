package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.AddCommand(rmCommerceCmd)
	rmCmd.AddCommand(rmNightlyCmd)
	rmCmd.AddCommand(rmReleaseCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes the Liferay Portal nook instance",
	Long:  `Removes the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var rmCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Removes the Liferay Portal Commerce instance",
	Long:  `Removes the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		RemoveDockerContainer(commerce)
	},
}

var rmNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Removes the Liferay Portal Nightly Build instance",
	Long:  `Removes the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		RemoveDockerContainer(nightly)
	},
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

// RemoveDockerContainer removes the running container
func RemoveDockerContainer(image liferay.Image) {
	err := docker.RemoveDockerContainer(image)
	if err != nil {
		log.Fatalln("Impossible to remove the container [" + image.GetContainerName() + "]")
	}
}
