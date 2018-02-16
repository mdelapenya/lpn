package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of lpn (Liferay Portal Nightly)",
	Long:  `All software has versions. This is lpn (Liferay Portal Nightly)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lpn (Liferay Portal Nightly) v0.1.2 -- HEAD")
	},
}
