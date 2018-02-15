package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lpn",
	Short: "lpn (Liferay Portal Nightly) makes it easier to run Liferay Portal's nightly builds.",
	Long: `A Fast and Flexible CLI for managing Liferay Portal's nightly builds built with
				love by mdelapenya and friends in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
