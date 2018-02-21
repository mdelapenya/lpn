package cmd

import (
	"errors"

	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	checkImageCmd.AddCommand(checkImageRelease)

	checkImageRelease.Flags().StringVarP(&tagToCheck, "tag", "t", "latest", "Sets the image tag to check")
}

var checkImageRelease = &cobra.Command{
	Use:   "release",
	Short: "Checks if the proper Liferay Portal release image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal release image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage release requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		CheckImage(release, tagToCheck)
	},
}
