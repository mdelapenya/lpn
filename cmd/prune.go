package cmd

import (
	docker "github.com/mdelapenya/lpn/docker"
	internal "github.com/mdelapenya/lpn/internal"
	liferay "github.com/mdelapenya/lpn/liferay"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pruneCmd)

	pruneCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs commands with Debug log level")
	pruneCmd.VisitParents(addVerboseFlag)
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prunes LPN state",
	Long:  `This command prunes LPN state: containers and images`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		enableDebugLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		images := []liferay.Image{
			liferay.CE{Tag: internal.LpnConfig.GetPortalImageTag("ce")},
			liferay.Commerce{Tag: internal.LpnConfig.GetPortalImageTag("commerce")},
			liferay.DXP{Tag: internal.LpnConfig.GetPortalImageTag("dxp")},
			liferay.Nightly{Tag: internal.LpnConfig.GetPortalImageTag("nightly")},
			liferay.Release{Tag: internal.LpnConfig.GetPortalImageTag("release")},
		}

		removeLPNContainers(images)
		removeLPNImages(images)

		log.Info("LPN state has been pruned!")
	},
}

func removeLPNContainers(images []liferay.Image) {
	for _, img := range images {
		docker.RemoveDockerContainer(img)
	}
}

func removeLPNImages(images []liferay.Image) {
	for _, img := range images {
		docker.RemoveDockerImage(img.GetFullyQualifiedName())
	}
}
