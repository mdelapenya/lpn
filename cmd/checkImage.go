package cmd

import (
	"errors"
	"log"

	docker "github.com/mdelapenya/lpn/docker"

	"github.com/spf13/cobra"
)

var tagToCheck string

func init() {
	rootCmd.AddCommand(checkImageCmd)

	checkImageCmd.Flags().StringVarP(&tagToCheck, "tag", "t", "latest", "Image tag to check")
}

var checkImageCmd = &cobra.Command{
	Use:   "checkImage",
	Short: "Check if the proper Liferay Portal nightly image has been pulled by lpn (Liferay Portal Nightly)",
	Long: `Uses docker image inspect to check if the proper Liferay Portal nightly image has 
	been pulled by lpn (Liferay Portal Nightly). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		dockerImage := docker.DockerImage + ":" + tagToCheck

		if docker.CheckDockerImageExists(dockerImage) {
			log.Println("The image [" + dockerImage + "] has been pulled from Docker Hub.")
		} else {
			log.Println("The image [" + dockerImage + "] has NOT been pulled from Docker Hub.")
		}
	},
}
