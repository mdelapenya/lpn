package cmd

import (
	"errors"

	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	checkImageCmd.AddCommand(checkImageCommerce)

	checkImageCommerce.Flags().StringVarP(&tagToCheck, "tag", "t", "latest", "Sets the image tag to check")
}

var checkImageCommerce = &cobra.Command{
	Use:   "commerce",
	Short: "Checks if the proper Liferay Portal with Commerce Build image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal with Commerce Build image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal with Commerce image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage commerce requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{Tag: tagToCheck}

		CheckImage(commerce)
	},
}
