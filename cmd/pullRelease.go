package cmd

import (
	"errors"

	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	pullCmd.AddCommand(pullRelease)

	pullRelease.Flags().StringVarP(&tagToPull, "tag", "t", "latest", "Sets the image tag to pull")
}

var pullRelease = &cobra.Command{
	Use:   "release",
	Short: "Pulls a Liferay Portal Docker image from releases",
	Long: `Pulls a Liferay Portal instance, obtained from the unofficial releases repository: "` + liferay.Releases + `".
	If no image tag is passed to the command, the "latest" tag will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull nightly requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{Tag: tagToPull}

		PullDockerImage(release)
	},
}
