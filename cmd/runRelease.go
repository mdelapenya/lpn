package cmd

import (
	"errors"

	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	runCmd.AddCommand(runReleaseCmd)

	runReleaseCmd.Flags().IntVarP(&httpPort, "httpPort", "p", 8080, "Sets the HTTP port of Liferay Portal's bundle.")
	runReleaseCmd.Flags().BoolVarP(&enableDebug, "debug", "d", false, "Enables debug mode. (default false)")
	runReleaseCmd.Flags().IntVarP(&debugPort, "debugPort", "D", 9000, "Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled")
	runReleaseCmd.Flags().IntVarP(&gogoPort, "gogoPort", "g", 11311, "Sets the GoGo Shell port of Liferay Portal's bundle.")
	runReleaseCmd.Flags().StringVarP(&tagToRun, "tag", "t", "latest", "Sets the image tag to run")
}

var runReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Runs a Liferay Portal instance from releases",
	Long: `Runs a Liferay Portal instance, obtained from the unofficial releases repository: ` + liferay.Releases + `.
	If no image tag is passed to the command, the "latest" tag will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{Tag: tagToRun}

		RunDockerImage(release, httpPort, gogoPort, enableDebug, debugPort)
	},
}
