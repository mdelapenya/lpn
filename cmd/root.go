package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lpn",
	Short: "lpn (Liferay Portal Nook) makes it easier to run Liferay Portal's Docker images.",
	Long: `A Fast and Flexible CLI for managing Liferay Portal's Docker images
				built with love by mdelapenya and friends in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Error executing lpn. Aborting")
	}
}

// SubCommandInfo Shows a message for subcommands
func SubCommandInfo() {
	// delegate to subcommands
	log.Warn(
		"Please run this command adding 'ce', 'commerce', 'dxp', 'nightly' or 'release' " +
			"subcommands.")
}
