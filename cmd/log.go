package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays logs for the Liferay Portal nightly instance",
	Long:  `Displays logs for the Liferay Portal nightly instance, identified by [` + docker.DockerContainerName + `].`,
	Run: func(cmd *cobra.Command, args []string) {
		docker.LogDockerContainer()
	},
}
