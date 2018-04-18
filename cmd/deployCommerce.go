package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	deployCmd.AddCommand(deployCommerce)

	deployCommerce.Flags().StringVarP(
		&filePath, "files", "f", "",
		`The file or files to deploy. A comma-separated list of files is accepted to deploy
							multiple files at the same time`)
	deployCommerce.Flags().StringVarP(
		&directoryPath, "dir", "d", "",
		`The directory to deploy its content. Only first-level files will be deployed, so no
							recursive deployment will happen`)
}

var deployCommerce = &cobra.Command{
	Use:   "commerce",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		validateArguments()

		if filePath != "" {
			deployFiles(commerce, filePath)
		}

		if directoryPath != "" {
			deployDirectory(commerce, directoryPath)
		}
	},
}
