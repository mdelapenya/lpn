package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "checkContainer",
	Short: "Checks if there is a container created by lpn",
	Long: `Checks if there is a container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a container with name [` + docker.DockerContainerName + `] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		docker.CheckDockerContainerExists()
	},
}
