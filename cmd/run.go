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
var debugPort int
var gogoPort int
var httpPort int
var memory string
var properties string
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
		subcommand.Flags().StringVarP(&properties, "properties", "P", "", "Sets the location of a portal-ext properties files to configure the running instance of Liferay Portal's bundle.")
		subcommand.Flags().StringVarP(&tagToRun, "tag", "t", "", "Sets the image tag to run")
	}

	runCECmd.Flags().StringVarP(&memory, "memory", "m", "-Xmx2048m", "Sets the memory for the JVM memory configuration of Liferay Portal's bundle.")
	runCommerceCmd.Flags().StringVarP(&memory, "memory", "m", "2048m", "Sets the memory for the Xmx and Xms JVM memory configuration of Liferay Portal's bundle.")
	runDXPCmd.Flags().StringVarP(&memory, "memory", "m", "-Xmx2048m", "Sets the memory for the JVM memory configuration of Liferay Portal's bundle.")
	runNightlyCmd.Flags().StringVarP(&memory, "memory", "m", "2048m", "Sets the memory for the Xmx and Xms JVM memory configuration of Liferay Portal's bundle.")
	runReleaseCmd.Flags().StringVarP(&memory, "memory", "m", "2048m", "Sets the memory for the Xmx and Xms JVM memory configuration of Liferay Portal's bundle.")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Liferay Portal instance",
	Long: `Runs a Liferay Portal instance, obtained from the Official repositories:
		- ` + liferay.CommercesRepository + ` (private),
		- ` + liferay.CERepository + `, and
		- ` + liferay.DXPRepository + `.
		For non-official Docker images, the tool runs images obtained from the unofficial repositories:
		- ` + liferay.NightliesRepository + `, and
		- ` + liferay.ReleasesRepository + `.
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
	Long: `Runs a Liferay Portal CE instance, obtained from the Official repository: ` + liferay.CERepository + `.
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

		runDockerImage(ce, httpPort, gogoPort, enableDebug, debugPort, memory, properties)
	},
}

var runCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Runs a Liferay Portal with Commerce instance from Commerce Builds",
	Long: `Runs a Liferay Portal with Commerce instance, obtained from Commerce Builds repository: ` + liferay.CommercesRepository + `.
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

		commerce := liferay.Commerce{Tag: tagToRun}

		runDockerImage(commerce, httpPort, gogoPort, enableDebug, debugPort, memory, properties)
	},
}

var runDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Runs a Liferay DXP instance",
	Long: `Runs a Liferay DXP instance, obtained from the Official repository: ` + liferay.DXPRepository + `,
	including a 30-day activation key.
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

		runDockerImage(dxp, httpPort, gogoPort, enableDebug, debugPort, memory, properties)
	},
}

var runNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Runs a Liferay Portal instance from Nightly Builds",
	Long: `Runs a Liferay Portal instance, obtained from Nightly Builds repository: ` + liferay.NightliesRepository + `.
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

		runDockerImage(nightly, httpPort, gogoPort, enableDebug, debugPort, memory, properties)
	},
}

var runReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Runs a Liferay Portal instance from releases",
	Long: `Runs a Liferay Portal instance, obtained from the unofficial releases repository: ` + liferay.ReleasesRepository + `.
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

		runDockerImage(release, httpPort, gogoPort, enableDebug, debugPort, memory, properties)
	},
}

// runDockerImage runs the image
func runDockerImage(
	image liferay.Image, httpPort int, gogoPort int, enableDebug bool, debugPort int, memory string,
	properties string) {

	err := docker.RunDockerImage(
		image, httpPort, gogoPort, enableDebug, debugPort, memory, properties)

	if err != nil {
		log.Fatalln("Impossible to run the container [" + image.GetContainerName() + "]")
	}

	log.Println("The container [" + image.GetContainerName() + "] has been run successfully")
}
