package cmd

import (
	"errors"

	date "github.com/mdelapenya/lpn/date"
	docker "github.com/mdelapenya/lpn/docker"

	"github.com/spf13/cobra"
)

var enableDebug bool
var debugPort int
var httpPort int

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&httpPort, "httpPort", "p", 8080, "Sets the HTTP port of Liferay Portal's bundle.")
	runCmd.Flags().BoolVarP(&enableDebug, "debug", "d", false, "Enables debug mode. (default false)")
	runCmd.Flags().IntVarP(&debugPort, "debugPort", "D", 9000, "Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled")
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

		docker.RunDockerImage(docker.DockerImage+":"+tag, httpPort, enableDebug, debugPort)
	},
}
