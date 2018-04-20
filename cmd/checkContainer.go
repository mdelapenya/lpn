package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
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

// CheckDockerContainerExists removes the running container
func CheckDockerContainerExists(image liferay.Image) {
	exists := docker.CheckDockerContainerExists(image)

	if !exists {
		log.Fatalln("The container [" + image.GetContainerName() + "] is NOT running.")
	}

	log.Println("The container [" + image.GetContainerName() + "] is running.")
}
