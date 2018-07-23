package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)

	subcommands := []*cobra.Command{checkContainerCommerceCmd, checkContainerNightlyCmd, checkContainerReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		checkCmd.AddCommand(subcommand)
	}
}

var checkCmd = &cobra.Command{
	Use:   "checkContainer",
	Short: "Checks if there is a container created by lpn",
	Long: `Checks if there is a container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a container with name "lpn" plus image type created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var checkContainerCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Checks if there is a Commerce container created by lpn",
	Long: `Checks if there is a Commerce container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Commerce container with name [lpn-commerce] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		checkDockerContainerExists(commerce)
	},
}

var checkContainerNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Checks if there is a Nightly Build container created by lpn",
	Long: `Checks if there is a Nightly Build container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Nightly Build container with name [lpn-nightly] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		checkDockerContainerExists(nightly)
	},
}

var checkContainerReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Checks if there is a Release container created by lpn",
	Long: `Checks if there is a Release container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Release container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		checkDockerContainerExists(release)
	},
}

// checkDockerContainerExists removes the running container
func checkDockerContainerExists(image liferay.Image) {
	exists := docker.CheckDockerContainerExists(image)

	if !exists {
		log.Fatalln("The container [" + image.GetContainerName() + "] is NOT running.")
	}

	log.Println("The container [" + image.GetContainerName() + "] is running.")
}
