package cmd

import (
	"log"
	"strings"

	"github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	deployCmd.AddCommand(deployRelease)

	deployRelease.Flags().StringVarP(
		&filePath, "files", "f", "",
		`The file or files to deploy. A comma-separated list of files is accepted to deploy
							multiple files at the same time`)
	deployRelease.Flags().StringVarP(
		&directoryPath, "dir", "d", "",
		`The directory to deploy its content. Only first-level files will be deployed, so no
							recursive deployment will happen`)
}

var deployRelease = &cobra.Command{
	Use:   "release",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long: `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn.
	The appropriate tag is calculated from the image the container was build from.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArguments()

		imageName, err := docker.GetDockerImageFromRunningContainer()
		if err != nil {
			log.Fatalln("The container [" + docker.DockerContainerName + "] is NOT running.")
		}

		index := strings.LastIndex(imageName, ":")

		tag := imageName[index+1 : len(imageName)-2]

		release := liferay.Release{Tag: tag}

		if filePath != "" {
			deployFiles(release, filePath)
		}

		if directoryPath != "" {
			deployDirectory(release, directoryPath)
		}
	},
}
