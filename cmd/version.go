package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of lpn",
	Long:  `All software has versions. This is lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := ioutil.ReadFile("VERSION.txt")

		if err != nil {
			fmt.Print(err)
		}

		fmt.Println("lpn (Liferay Portal Nook) v" + string(version) + " -- HEAD")
	},
}
