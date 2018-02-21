package cmd

import (
	"strings"

	"github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	deployCmd.AddCommand(deployRelease)

	deployRelease.Flags().StringVarP(&filePath, "file", "f", "", "The file to deploy")
}

var deployRelease = &cobra.Command{
	Use:   "release",
	Short: "Deploys a file to Liferay Portal's deploy folder in the container run by lpn",
	Long: `Deploys a file to Liferay Portal's deploy folder in the container run by lpn.
	The appropriate tag is calculated from the image the container was build from.`,
	Run: func(cmd *cobra.Command, args []string) {
		imageName := docker.GetDockerImageFromRunningContainer()

		index := strings.LastIndex(imageName, ":")

		tag := imageName[index+1 : len(imageName)-2]

		release := liferay.Release{Tag: tag}

		deployFile(release, filePath)
	},
}
