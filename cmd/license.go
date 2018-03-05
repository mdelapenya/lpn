package cmd

import (
	"fmt"

	license "github.com/mdelapenya/lpn/assets/license"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(licenseCmd)
}

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Print the license of lpn",
	Long:  `All software has a license. This is lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		licenseFile, err := license.Asset("LICENSE.txt")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(string(licenseFile))
	},
}
