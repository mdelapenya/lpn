package cmd

import (
	"errors"

	date "github.com/mdelapenya/lpn/date"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	runCmd.AddCommand(runCommerceCmd)

	runCommerceCmd.Flags().IntVarP(&httpPort, "httpPort", "p", 8080, "Sets the HTTP port of Liferay Portal's bundle.")
	runCommerceCmd.Flags().BoolVarP(&enableDebug, "debug", "d", false, "Enables debug mode. (default false)")
	runCommerceCmd.Flags().IntVarP(&debugPort, "debugPort", "D", 9000, "Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled")
	runCommerceCmd.Flags().IntVarP(&gogoPort, "gogoPort", "g", 11311, "Sets the GoGo Shell port of Liferay Portal's bundle.")
	runCommerceCmd.Flags().StringVarP(&tagToRun, "tag", "t", date.CurrentDate, "Sets the image tag to run")
}

var runCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Runs a Liferay Portal with Commerce instance from Commerce Builds",
	Long: `Runs a Liferay Portal with Commerce instance, obtained from Commerce Builds repository: ` + liferay.CommerceRepository + `.
	If no image tag is passed to the command, the tag representing the current date [` + date.CurrentDate + `] will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{Tag: tagToRun}

		RunDockerImage(commerce, httpPort, gogoPort, enableDebug, debugPort)
	},
}
