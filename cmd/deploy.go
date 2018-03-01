package cmd

import (
	"log"
	"os"
	"strings"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var filePath string

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a file to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys a file to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

// deployFiles deploys a file to the running container
func deployFiles(image liferay.Image, path string) {
	paths := strings.Split(path, ",")

	filesChannel := make(chan string, len(paths))
	for i := range paths {
		filesChannel <- paths[i]
	}
	close(filesChannel)

	workers := 8
	if len(paths) < workers {
		workers = len(paths)
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
			log.Println("[" + paths[i] + "] deployed sucessfully to " + image.GetDeployFolder())
		case err := <-errorChannel:
			log.Println("Impossible to deploy the file to the container", err)
		}
	}
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
