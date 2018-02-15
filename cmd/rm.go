package cmd

import (
	docker "lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes the Liferay Portal nightly instance",
	Long:  `Removes the Liferay Portal nightly instance, identified by [` + docker.DockerContainerName + `].`,
	Run: func(cmd *cobra.Command, args []string) {
		docker.RemoveDockerContainer()
	},
}
