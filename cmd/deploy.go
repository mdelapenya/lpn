package cmd

import (
	"strings"
	"time"

	"github.com/mdelapenya/lpn/docker"
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

// deployFile deploys a file to the running container
func deployFile(image liferay.Image, path string) {
	paths := strings.Split(path, ",")

	for i := range paths {
		go docker.CopyFileToContainer(image, paths[i])
		time.Sleep(1 * time.Second)
	}
}
