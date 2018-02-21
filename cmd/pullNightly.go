package cmd

import (
	"errors"

	date "github.com/mdelapenya/lpn/date"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	pullCmd.AddCommand(pullNightly)
}

var pullNightly = &cobra.Command{
	Use:   "nightly",
	Short: "Pulls a Liferay Portal Docker image from Nightly Builds",
	Long: `Pulls a Liferay Portal Docker image from the Nighlty Builds repository: "` + liferay.Nightlies + `".
	If no image tag is passed to the command, the tag representing the current date [` + date.CurrentDate + `] will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull nightly requires zero or one argument representing the image tag to be pulled")
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

		nightly := liferay.Nightly{}

		PullDockerImage(nightly, tag)
	},
}
