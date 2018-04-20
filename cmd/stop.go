package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the Liferay Portal nook instance",
	Long:  `Stops the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

// StopDockerContainer stops the running container
func StopDockerContainer(image liferay.Image) {
	err := docker.StopDockerContainer(image)
	if err != nil {
		log.Fatalln("Impossible to stop the container [" + image.GetContainerName() + "]")
	}
}
