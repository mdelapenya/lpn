package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the Liferay Portal nook instance",
	Long:  `Stops the Liferay Portal nook instance, identified by [` + docker.DockerContainerName + `].`,
	Run: func(cmd *cobra.Command, args []string) {
		err := docker.StopDockerContainer()
		if err != nil {
			log.Fatalln("Impossible to stop the container")
		}
	},
}
