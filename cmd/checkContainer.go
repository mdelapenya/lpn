package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)

	subcommands := []*cobra.Command{
		checkContainerCECmd, checkContainerCommerceCmd, checkContainerDXPCmd,
		checkContainerNightlyCmd, checkContainerReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		subcommand.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs commands with Debug log level")
		subcommand.VisitParents(addVerboseFlag)

		checkCmd.AddCommand(subcommand)
	}
}

var checkCmd = &cobra.Command{
	Use:   "checkContainer",
	Short: "Checks if there is a container created by lpn",
	Long: `Checks if there is a container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a container with name "lpn" plus image type created by lpn (Liferay Portal Nook)`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		enableDebugLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var checkContainerCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Checks if there is a CE container created by lpn",
	Long: `Checks if there is a CE container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a CE container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		checkDockerContainerExists(ce)
	},
}

var checkContainerCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Checks if there is a Commerce container created by lpn",
	Long: `Checks if there is a Commerce container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Commerce container with name [lpn-commerce] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		checkDockerContainerExists(commerce)
	},
}

var checkContainerDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Checks if there is a DXP container created by lpn",
	Long: `Checks if there is a DXP container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a DXP container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		checkDockerContainerExists(dxp)
	},
}

var checkContainerNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Checks if there is a Nightly Build container created by lpn",
	Long: `Checks if there is a Nightly Build container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Nightly Build container with name [lpn-nightly] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		checkDockerContainerExists(nightly)
	},
}

var checkContainerReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Checks if there is a Release container created by lpn",
	Long: `Checks if there is a Release container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Release container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		checkDockerContainerExists(release)
	},
}

// checkDockerContainerExists removes the running container
func checkDockerContainerExists(image liferay.Image) {
	exists := docker.CheckDockerContainerExists(image.GetContainerName())

	if !exists {
		log.WithFields(log.Fields{
			"container": image.GetContainerName(),
		}).Warn("Container does NOT exist in the system.")
		return
	}

	log.WithFields(log.Fields{
		"container": image.GetContainerName(),
	}).Info("Container DOES exist in the system.")
}
