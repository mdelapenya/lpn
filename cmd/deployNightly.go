package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	deployCmd.AddCommand(deployNightly)

	deployNightly.Flags().StringVarP(&filePath, "file", "f", "", "The file to deploy")
}

var deployNightly = &cobra.Command{
	Use:   "nightly",
	Short: "Deploys a file to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys a file to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		deployFiles(nightly, filePath)
	},
}
