package cmd

import (
	"errors"
	date "lpn/date"
	docker "lpn/docker"

	"github.com/spf13/cobra"
)

var httpPort int

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&httpPort, "httpPort", "p", 8080, "HTTP Port")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Liferay Portal nightly instance",
	Long: `Runs a Liferay Portal nightly instance, obtained from ` + docker.DockerImage + `.
	If no image tag is passed to the command, the tag representing the current date [` + date.CurrentDate + `]
	will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
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

		docker.RunDockerImage(docker.DockerImage+":"+tag, httpPort)
	},
}
