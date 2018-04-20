package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes the Liferay Portal nook instance",
	Long:  `Removes the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

// RemoveDockerContainer removes the running container
func RemoveDockerContainer(image liferay.Image) {
	err := docker.RemoveDockerContainer(image)
	if err != nil {
		log.Fatalln("Impossible to remove the container [" + image.GetContainerName() + "]")
	}
}
