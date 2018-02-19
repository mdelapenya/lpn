package cmd

import (
	"errors"
	"fmt"

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
		// delegate to subcommands
		fmt.Println("Please run this command adding 'nightly' or 'release' subcommands.")
	},
}
