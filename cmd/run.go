package cmd

import (
	"errors"
	"log"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var enableDebug bool
var debugPort int
var gogoPort int
var httpPort int
var tagToRun string

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Liferay Portal instance",
	Long: `Runs a Liferay Portal instance, obtained from the unofficial repositories: ` + liferay.Releases +
		` or ` + liferay.Nightlies + `.
		For that, please run this command adding "release" or "nightly" subcommands.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

// RunDockerImage runs the image
func RunDockerImage(
	image liferay.Image, httpPort int, gogoPort int, enableDebug bool, debugPort int) {

	err := docker.RunDockerImage(
		image.GetFullyQualifiedName(), httpPort, gogoPort, enableDebug, debugPort)

	if err != nil {
		log.Fatalln("Impossible to run the container")
	}

	log.Println("The container [" + docker.DockerContainerName + "] has been run sucessfully")
}
