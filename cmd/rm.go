package cmd

import (
	"log"

	docker "github.com/mdelapenya/lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes the Liferay Portal nook instance",
	Long:  `Removes the Liferay Portal nook instance, identified by [` + docker.DockerContainerName + `].`,
	Run: func(cmd *cobra.Command, args []string) {
		err := docker.RemoveDockerContainer()
		if err != nil {
			log.Fatalln("Impossible to remove the container")
		}
	},
}
