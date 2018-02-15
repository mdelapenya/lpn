package cmd

import (
	"bufio"
	"fmt"
	docker "lpn/docker"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vjeantet/jodaTime"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Liferay Portal nightly instance",
	Long:  `Runs a Liferay Portal nightly instance, using current date if no version is passed`,
	Run: func(cmd *cobra.Command, args []string) {
		var currentDate = jodaTime.Format("YYYYMMdd", time.Now())

		fmt.Println("Enter the Image Tag you want to use for [" + docker.DockerImage + "]")
		fmt.Print("If you leave it empty, we will use [" + currentDate + "]: ")

		var imageTag string
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			imageTag = scanner.Text()
		}

		if imageTag == "" {
			imageTag = currentDate
		}

		docker.RunDockerImage(docker.DockerImage + ":" + imageTag)
	},
}
