package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var filePath string
var directoryPath string

func init() {
	rootCmd.AddCommand(deployCmd)

	subcommands := []*cobra.Command{deployCommerce, deployNightly, deployRelease}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		deployCmd.AddCommand(subcommand)

		subcommand.Flags().StringVarP(
			&filePath, "files", "f", "",
			`The file or files to deploy. A comma-separated list of files is accepted to deploy
							multiple files at the same time`)

		subcommand.Flags().StringVarP(
			&directoryPath, "dir", "d", "",
			`The directory to deploy its content. Only first-level files will be deployed, so no
							recursive deployment will happen`)
	}
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var deployCommerce = &cobra.Command{
	Use:   "commerce",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		validateArguments()

		imageName, err := docker.GetDockerImageFromRunningContainer(commerce)
		if err != nil {
			log.Fatalln("The container [" + commerce.GetContainerName() + "] is NOT running.")
		}

		index := strings.LastIndex(imageName, ":")

		tag := imageName[index+1 : len(imageName)]

		commerce.Tag = tag

		if filePath != "" {
			deployFiles(commerce, filePath)
		}

		if directoryPath != "" {
			deployDirectory(commerce, directoryPath)
		}
	},
}

var deployNightly = &cobra.Command{
	Use:   "nightly",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		validateArguments()

		imageName, err := docker.GetDockerImageFromRunningContainer(nightly)
		if err != nil {
			log.Fatalln("The container [" + nightly.GetContainerName() + "] is NOT running.")
		}

		index := strings.LastIndex(imageName, ":")

		tag := imageName[index+1 : len(imageName)]

		nightly.Tag = tag

		if filePath != "" {
			deployFiles(nightly, filePath)
		}

		if directoryPath != "" {
			deployDirectory(nightly, directoryPath)
		}
	},
}

var deployRelease = &cobra.Command{
	Use:   "release",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long: `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn.
	The appropriate tag is calculated from the image the container was build from.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArguments()

		release := liferay.Release{}

		imageName, err := docker.GetDockerImageFromRunningContainer(release)
		if err != nil {
			log.Fatalln("The container [" + release.GetContainerName() + "] is NOT running.")
		}

		index := strings.LastIndex(imageName, ":")

		tag := imageName[index+1 : len(imageName)]

		release.Tag = tag

		if filePath != "" {
			deployFiles(release, filePath)
		}

		if directoryPath != "" {
			deployDirectory(release, directoryPath)
		}
	},
}

// deployDirectory deploys a directory's content to the running container
func deployDirectory(image liferay.Image, dirPath string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalln("The directory is not valid", err)
	}

	var onlyFilePaths []string

	for _, f := range files {
		if !f.Mode().IsDir() {
			onlyFilePaths = append(onlyFilePaths, path.Join(dirPath, f.Name()))
		}
	}

	deployPaths(image, onlyFilePaths)
}

// deployFiles deploys files to the running container
func deployFiles(image liferay.Image, path string) {
	paths := strings.Split(path, ",")

	deployPaths(image, paths)
}

func deployFile(
	file string, image liferay.Image, resultChannel chan bool,
	errorChannel chan error) {

	if _, err := os.Stat(file); os.IsNotExist(err) {
		select {
		case errorChannel <- err:
			// will break parent goroutine out of loop
		default:
			// don't care, first error wins
		}
		return
	}

	err := docker.CopyFileToContainer(image, file)
	if err != nil {
		select {
		case errorChannel <- err:
			// will break parent goroutine out of loop
		default:
			// don't care, first error wins
		}
		return
	}

	resultChannel <- true
}

// deployPaths deploys files to the running container
func deployPaths(image liferay.Image, paths []string) {
	filesChannel := make(chan string, len(paths))
	for i := range paths {
		filesChannel <- paths[i]
	}
	close(filesChannel)

	workers := 8
	if len(filesChannel) < workers {
		workers = len(filesChannel)
	}

	errorChannel := make(chan error, 1)
	resultChannel := make(chan bool, len(paths))

	for i := 0; i < workers; i++ {
		// Consume work from filesChannel. Loop will end when no more work.
		for file := range filesChannel {
			go deployFile(file, image, resultChannel, errorChannel)
		}
	}

	// Collect results from workers

	for i := 0; i < len(paths); i++ {
		select {
		case <-resultChannel:
			log.Println("[" + paths[i] + "] deployed successfully to " + image.GetDeployFolder())
		case err := <-errorChannel:
			log.Println("Impossible to deploy the file to the container", err)
		}
	}
}

func validateArguments() {
	if filePath == "" && directoryPath == "" {
		log.Fatalln("Please pass a valid path to a file or to a directory as argument")
	}
}
