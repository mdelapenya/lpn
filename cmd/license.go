package cmd

import (
	"fmt"

	license "github.com/mdelapenya/lpn/assets/license"
	log "github.com/sirupsen/logrus"

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
			log.WithFields(log.Fields{
				"error":   err,
				"command": cmd.Use,
			}).Error("Error executing Command")
			return
		}

		fmt.Println(string(licenseFile))
	},
}
