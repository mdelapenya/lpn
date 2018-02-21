package cmd

import (
	"errors"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	checkImageCmd.AddCommand(checkImageNightly)

	checkImageNightly.Flags().StringVarP(&tagToCheck, "tag", "t", "latest", "Sets the image tag to check")
}

var checkImageNightly = &cobra.Command{
	Use:   "nightly",
	Short: "Checks if the proper Liferay Portal Nightly Build image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal Nightly Build image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage nightly requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		dockerImage := nightly.GetRepository() + ":" + tagToCheck

		docker.CheckDockerImageExists(dockerImage)
	},
}
