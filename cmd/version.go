package cmd

import (
	"fmt"

	v "github.com/mdelapenya/lpn/assets/version"
	docker "github.com/mdelapenya/lpn/docker"

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
		version, err := v.Asset("VERSION.txt")

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		dockerVersion, _ := docker.GetDockerVersion()

		fmt.Println("lpn (Liferay Portal Nook) v" + string(version) + " -- HEAD")
		fmt.Println("Docker version:")
		fmt.Println(dockerVersion)
	},
}
