package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays logs for the Liferay Portal instance",
	Long:  `Displays logs for the Liferay Portal instance, identified by [lpn] plus image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

// LogDockerContainer show the logs for the running container of the specified type
func LogDockerContainer(image liferay.Image) {
	docker.LogDockerContainer(image)
}
