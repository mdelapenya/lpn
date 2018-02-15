package cmd

import (
	"errors"
	"log"
	docker "lpn/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkImageCmd)
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
		var tag string

		if len(args) == 0 {
			tag = "latest"
		} else {
			tag = args[0]
		}

		dockerImage := docker.DockerImage + ":" + tag

		if docker.CheckDockerImageExists(dockerImage) {
			log.Println("The image [" + dockerImage + "] has been pulled from Docker Hub.")
		} else {
			log.Println("The image [" + dockerImage + "] has NOT been pulled from Docker Hub.")
		}
	},
}
