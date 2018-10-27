package cmd

import (
	"log"

	date "github.com/mdelapenya/lpn/date"
	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

var tagToRemove string

func init() {
	rootCmd.AddCommand(rmiCmd)

	subcommands := []*cobra.Command{rmiCommerceCmd, rmiNightlyCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		rmiCmd.AddCommand(subcommand)

		subcommand.Flags().StringVarP(&tagToRemove, "tag", "t", date.CurrentDate, "Sets the image tag to remove")
	}

	rmiReleaseCmd.Flags().StringVarP(&tagToRemove, "tag", "t", "latest", "Sets the image tag to remove")
	rmiCmd.AddCommand(rmiReleaseCmd)
}

var rmiCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Removes the Liferay Portal image",
	Long:  `Removes the Liferay Portal image related to the lpn instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var rmiCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Removes the Liferay Portal Commerce image",
	Long:  `Removes the Liferay Portal Commerce image, identified by ["` + liferay.CommercesRepository + `"].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{Tag: tagToRemove}

		removeDockerImage(commerce)
	},
}

var rmiNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Removes the Liferay Portal Nightly Build image",
	Long:  `Removes the Liferay Portal Nightly Build image, identified by ["` + liferay.NightliesRepository + `"].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{Tag: tagToRemove}

		removeDockerImage(nightly)
	},
}

var rmiReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Removes the Liferay Portal Release image",
	Long:  `Removes the Liferay Portal Release image, identified by ["` + liferay.ReleasesRepository + `"].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{Tag: tagToRemove}

		removeDockerImage(release)
	},
}

// removeDockerImage removes the running container
func removeDockerImage(image liferay.Image) {
	err := docker.RemoveDockerImage(image.GetFullyQualifiedName())
	if err != nil {
		log.Fatalln("Impossible to remove the image [" + image.GetFullyQualifiedName() + "]")
	}
}
