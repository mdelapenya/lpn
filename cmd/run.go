package cmd

import (
	"errors"
	"log"

	date "github.com/mdelapenya/lpn/date"
	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var enableDebug bool
var datastore string
var debugPort int
var gogoPort int
var httpPort int
var memory string
var tagToRun string

func init() {
	rootCmd.AddCommand(runCmd)

	subcommands := []*cobra.Command{
		runCECmd, runCommerceCmd, runDXPCmd, runNightlyCmd, runReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		runCmd.AddCommand(subcommand)

		subcommand.Flags().IntVarP(&httpPort, "httpPort", "p", 8080, "Sets the HTTP port of Liferay Portal's bundle.")
		subcommand.Flags().BoolVarP(&enableDebug, "debug", "d", false, "Enables debug mode. (default false)")
		subcommand.Flags().IntVarP(&debugPort, "debugPort", "D", 9000, "Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled")
		subcommand.Flags().IntVarP(&gogoPort, "gogoPort", "g", 11311, "Sets the GoGo Shell port of Liferay Portal's bundle.")
		subcommand.Flags().StringVarP(&datastore, "datastore", "s", "hsql", "Creates a database service for the running instance. Supported values are [hsql|mysql|postgresql] (default HSQL)")
		subcommand.Flags().StringVarP(&tagToRun, "tag", "t", "", "Sets the image tag to run")
	}

	runCECmd.Flags().StringVarP(&memory, "memory", "m", "-Xmx2048m", "Sets the memory for the JVM memory configuration of Liferay Portal's bundle.")
	runCommerceCmd.Flags().StringVarP(&memory, "memory", "m", "-Xmx2048m", "Sets the memory for the Xmx and Xms JVM memory configuration of Liferay Portal's bundle.")
	runDXPCmd.Flags().StringVarP(&memory, "memory", "m", "-Xmx2048m", "Sets the memory for the JVM memory configuration of Liferay Portal's bundle.")
	runNightlyCmd.Flags().StringVarP(&memory, "memory", "m", "-Xmx2048m", "Sets the memory for the Xmx and Xms JVM memory configuration of Liferay Portal's bundle.")
	runReleaseCmd.Flags().StringVarP(&memory, "memory", "m", "2048m", "Sets the memory for the Xmx and Xms JVM memory configuration of Liferay Portal's bundle.")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Liferay Portal instance",
	Long: `Runs a Liferay Portal instance, obtained from the Official repositories (see configuration file).
		For non-official Docker images, the tool runs images obtained from the unofficial repositories (see configuration file).
	For that, please run this command adding "ce", "commerce", "dxp", "release" or "nightly" subcommands.`,
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

var runCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Runs a Liferay Portal CE instance",
	Long: `Runs a Liferay Portal CE instance, obtained from the Official repository.
	If no image tag is passed to the command, the "` + liferay.CEDefaultTag + `" tag will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRun == "" {
			tagToRun = liferay.CEDefaultTag
		}

		ce := liferay.CE{Tag: tagToRun}

		runLiferayDockerImage(ce, datastore, httpPort, gogoPort, enableDebug, debugPort, memory)
	},
}

var runCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Runs a Liferay Portal with Commerce instance from Commerce Builds",
	Long: `Runs a Liferay Portal with Commerce instance, obtained from Commerce Builds repository.
	If no image tag is passed to the command, the "` + liferay.CommerceDefaultTag + `" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRun == "" {
			tagToRun = date.CurrentDate
		}

		commerce := liferay.Commerce{Tag: tagToRun}

		runLiferayDockerImage(
			commerce, datastore, httpPort, gogoPort, enableDebug, debugPort, memory)
	},
}

var runDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Runs a Liferay DXP instance",
	Long: `Runs a Liferay DXP instance, obtained from the Official repository, including a 30-day activation key.
	If no image tag is passed to the command, the "` + liferay.DXPDefaultTag + `" tag will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRun == "" {
			tagToRun = liferay.DXPDefaultTag
		}

		dxp := liferay.DXP{Tag: tagToRun}

		runLiferayDockerImage(dxp, datastore, httpPort, gogoPort, enableDebug, debugPort, memory)
	},
}

var runNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Runs a Liferay Portal instance from Nightly Builds",
	Long: `Runs a Liferay Portal instance, obtained from Nightly Builds repository.
	If no image tag is passed to the command, the tag representing the current date [` + date.CurrentDate + `] will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRun == "" {
			tagToRun = date.CurrentDate
		}

		nightly := liferay.Nightly{Tag: tagToRun}

		runLiferayDockerImage(
			nightly, datastore, httpPort, gogoPort, enableDebug, debugPort, memory)
	},
}

var runReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Runs a Liferay Portal instance from releases",
	Long: `Runs a Liferay Portal instance, obtained from the unofficial releases repository.
	If no image tag is passed to the command, the "latest" tag will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("run requires zero or one argument representing the image tag to be run")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRun == "" {
			tagToRun = "latest"
		}

		release := liferay.Release{Tag: tagToRun}

		runLiferayDockerImage(
			release, datastore, httpPort, gogoPort, enableDebug, debugPort, memory)
	},
}

// runLiferayDockerImage runs the Liferay image, potentially with a datastore
func runLiferayDockerImage(
	image liferay.Image, datastore string, httpPort int, gogoPort int, enableDebug bool,
	debugPort int, memory string) {

	if datastore != "hsql" {
		database := docker.GetDatabase(image, datastore)

		err := docker.RunLiferayDockerImage(
			image, database, httpPort, gogoPort, enableDebug, debugPort, memory)

		if err != nil {
			log.Fatalln("Impossible to run the stack for [" + image.GetContainerName() + "]")
		}

		log.Println("The stack for [" + image.GetContainerName() + "] has been run successfully")
	} else {
		err := docker.RunLiferayDockerImage(
			image, nil, httpPort, gogoPort, enableDebug, debugPort, memory)

		if err != nil {
			log.Fatalln("Impossible to run the container [" + image.GetContainerName() + "]")
		}

		log.Println("The container [" + image.GetContainerName() + "] has been run successfully")
	}
}
