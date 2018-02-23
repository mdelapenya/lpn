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

// CheckImage uses the image interface to check if it exists
func CheckImage(image liferay.Image) {
	err := docker.CheckDockerImageExists(image.GetFullyQualifiedName())
	if err == nil {
		log.Println("The image [" + image.GetFullyQualifiedName() + "] has been pulled from Docker Hub.")
	} else {
		log.Println("The image [" + image.GetFullyQualifiedName() + "] has NOT been pulled from Docker Hub.")
	}
}
