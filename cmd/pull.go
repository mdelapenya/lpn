package cmd

import (
	"errors"
	"fmt"

	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a Liferay Portal Docker image",
	Long: `Pulls a Liferay Portal Docker image from the unofficial repositories "` + liferay.GetReleasesRepository() + `" and "` + liferay.GetNightlyBuildsRepository() + `".
	For that, please run this command adding "release" or "nightly" subcommands.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// delegate to subcommands
		fmt.Println("Please run this command adding 'nightly' or 'release' subcommands.")
	},
}
