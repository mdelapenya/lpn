package cmd

import (
	"errors"
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var tagToPull string
var forceRemoval bool

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a Liferay Portal Docker image",
	Long: `Pulls a Liferay Portal Docker image from the unofficial repositories "` + liferay.ReleasesRepository + `" and "` + liferay.NightliesRepository + `".
	For that, please run this command adding "release" or "nightly" subcommands.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

// PullDockerImage uses the image interface to pull it from Docker Hub, removing the cached on if
func PullDockerImage(image liferay.Image, forceRemoval bool) {
	if forceRemoval {
		err := docker.RemoveDockerImage(image.GetFullyQualifiedName())
		if err != nil {
			log.Println(
				"The image " + image.GetFullyQualifiedName() +
					" was not found in th local cache. Skipping removal")
		}
	}

	docker.PullDockerImage(image.GetFullyQualifiedName())
}
