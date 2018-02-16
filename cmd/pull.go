package cmd

import (
	"errors"

	date "github.com/mdelapenya/lpn/date"
	docker "github.com/mdelapenya/lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a Liferay Portal Docker image",
	Long: `Pulls a Liferay Portal Docker image from the unofficial repositories ` + docker.DockerImage + ` and ` + docker.DockerImage + `.
	If no image tag is passed to the command, the tag representing the current date [` + date.CurrentDate + `]
	will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var tag string

		if len(args) == 0 {
			tag = date.CurrentDate
		} else {
			tag = args[0]
		}

		docker.PullDockerImage(docker.DockerImage + ":" + tag)
	},
}
