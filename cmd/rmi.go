package cmd

import (
	date "github.com/mdelapenya/lpn/date"
	docker "github.com/mdelapenya/lpn/docker"
	internal "github.com/mdelapenya/lpn/internal"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var tagToRemove string

func init() {
	rootCmd.AddCommand(rmiCmd)

	subcommands := []*cobra.Command{rmiCECmd, rmiCommerceCmd, rmiDXPCmd, rmiNightlyCmd, rmiReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		subcommand.Flags().StringVarP(&tagToRemove, "tag", "t", "", "Sets the image tag to remove")

		rmiCmd.AddCommand(subcommand)

		subcommand.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs commands with Debug log level")
		subcommand.VisitParents(addVerboseFlag)
	}
}

var rmiCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Removes the Liferay Portal image",
	Long:  `Removes the Liferay Portal image related to the lpn instances.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		enableDebugLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var rmiCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Removes the Official Liferay Portal CE image",
	Long:  `Removes the Official Liferay Portal CE image from the Docker host.`,
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRemove == "" {
			tagToRemove = internal.LpnConfig.GetPortalImageTag("ce")
		}

		ce := liferay.CE{Tag: tagToRemove}

		removeDockerImage(ce)
	},
}

var rmiCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Removes the Liferay Portal Commerce image",
	Long:  `Removes the Liferay Portal Commerce image from the Docker host.`,
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRemove == "" {
			tagToRemove = date.CurrentDate
		}

		commerce := liferay.Commerce{Tag: tagToRemove}

		removeDockerImage(commerce)
	},
}

var rmiDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Removes the Official Liferay DXP image",
	Long:  `Removes the Official Liferay DXP image from the Docker host.`,
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRemove == "" {
			tagToRemove = internal.LpnConfig.GetPortalImageTag("dxp")
		}

		dxp := liferay.DXP{Tag: tagToRemove}

		removeDockerImage(dxp)
	},
}

var rmiNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Removes the Liferay Portal Nightly Build image",
	Long:  `Removes the Liferay Portal Nightly Build image from the Docker host.`,
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRemove == "" {
			tagToRemove = date.CurrentDate
		}

		nightly := liferay.Nightly{Tag: tagToRemove}

		removeDockerImage(nightly)
	},
}

var rmiReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Removes the Liferay Portal Release image",
	Long:  `Removes the Liferay Portal Release image from the Docker host.`,
	Run: func(cmd *cobra.Command, args []string) {
		if tagToRemove == "" {
			tagToRemove = "latest"
		}

		release := liferay.Release{Tag: tagToRemove}

		removeDockerImage(release)
	},
}

// removeDockerImage removes the running container
func removeDockerImage(image liferay.Image) {
	docker.RemoveDockerImage(image.GetFullyQualifiedName())
}
