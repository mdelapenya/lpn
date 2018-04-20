package cmd

import (
	"errors"
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var tagToCheck string

func init() {
	rootCmd.AddCommand(checkImageCmd)

	checkImageCmd.AddCommand(checkImageCommerce)
	checkImageCommerce.Flags().StringVarP(&tagToCheck, "tag", "t", "latest", "Sets the image tag to check")

	checkImageCmd.AddCommand(checkImageNightly)
	checkImageNightly.Flags().StringVarP(&tagToCheck, "tag", "t", "latest", "Sets the image tag to check")

	checkImageCmd.AddCommand(checkImageRelease)
	checkImageRelease.Flags().StringVarP(&tagToCheck, "tag", "t", "latest", "Sets the image tag to check")
}

var checkImageCmd = &cobra.Command{
	Use:   "checkImage",
	Short: "Checks if the proper Liferay Portal image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal image has been pulled by lpn.
	Uses "docker image inspect" to check if the proper Liferay Portal image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
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

		checkImage(commerce)
	},
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
		nightly := liferay.Nightly{Tag: tagToCheck}

		checkImage(nightly)
	},
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
		release := liferay.Release{Tag: tagToCheck}

		checkImage(release)
	},
}

// checkImage uses the image interface to check if it exists
func checkImage(image liferay.Image) {
	err := docker.CheckDockerImageExists(image.GetFullyQualifiedName())
	if err != nil {
		log.Fatalln("The image [" + image.GetFullyQualifiedName() + "] has NOT been pulled from Docker Hub.")
	} else {
		log.Println("The image [" + image.GetFullyQualifiedName() + "] has been pulled from Docker Hub.")
	}
}
