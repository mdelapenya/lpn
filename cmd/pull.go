package cmd

import (
	"errors"

	date "github.com/mdelapenya/lpn/date"
	docker "github.com/mdelapenya/lpn/docker"
	internal "github.com/mdelapenya/lpn/internal"
	liferay "github.com/mdelapenya/lpn/liferay"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tagToPull string
var forceRemoval bool

func init() {
	rootCmd.AddCommand(pullCmd)

	subcommands := []*cobra.Command{pullCE, pullCommerce, pullDXP, pullNightly, pullRelease}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		pullCmd.AddCommand(subcommand)

		subcommand.Flags().BoolVarP(&forceRemoval, "forceRemoval", "f", false, "Removes the cached, local image, if exists")
		subcommand.Flags().StringVarP(&tagToPull, "tag", "t", "", "Sets the image tag to pull")
	}
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a Liferay Portal Docker image",
	Long: `Pulls a Liferay Portal Docker image from one of the Official repositories (see configuration file).
		For non-official Docker images, the tool pulls from the official repositories (see configuration file)
	For that, please run this command adding "ce", "commerce", "dxp", "release" or "nightly" subcommands.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var pullCE = &cobra.Command{
	Use:   "ce",
	Short: "Pulls a Liferay Portal CE Docker image from Official CE repository",
	Long: `Pulls a Liferay Portal instance, obtained from the official CE repository.
	If no image tag is passed to the command, the default tag (see configuration file) will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull ce requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToPull == "" {
			tagToPull = internal.LpnConfig.GetPortalImageTag("ce")
		}

		ce := liferay.CE{Tag: tagToPull}

		pullDockerImage(ce, forceRemoval)
	},
}

var pullCommerce = &cobra.Command{
	Use:   "commerce",
	Short: "Pulls a Liferay Portal Docker image from Commerce Builds",
	Long: `Pulls a Liferay Portal Docker image from the Commerce Builds repository".
	If no image tag is passed to the command, the default tag (see configuration file) will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull commerce requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToPull == "" {
			tagToPull = internal.LpnConfig.GetPortalImageTag("commerce")
		}

		commerce := liferay.Commerce{Tag: tagToPull}

		pullDockerImage(commerce, forceRemoval)
	},
}

var pullDXP = &cobra.Command{
	Use:   "dxp",
	Short: "Pulls a Liferay DXP Docker image from Official DXP repository",
	Long: `Pulls a Liferay DXP instance, obtained from the official DXP repository, including a 30-day activation key.
	If no image tag is passed to the command, the default tag (see configuration file) will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull ce requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToPull == "" {
			tagToPull = internal.LpnConfig.GetPortalImageTag("dxp")
		}

		dxp := liferay.DXP{Tag: tagToPull}

		pullDockerImage(dxp, forceRemoval)
	},
}

var pullNightly = &cobra.Command{
	Use:   "nightly",
	Short: "Pulls a Liferay Portal Docker image from Nightly Builds",
	Long: `Pulls a Liferay Portal Docker image from the Nightly Builds repository.
	If no image tag is passed to the command, the "master" tag will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull nightly requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToPull == "" {
			tagToPull = date.CurrentDate
		}

		nightly := liferay.Nightly{Tag: tagToPull}

		pullDockerImage(nightly, forceRemoval)
	},
}

var pullRelease = &cobra.Command{
	Use:   "release",
	Short: "Pulls a Liferay Portal Docker image from releases",
	Long: `Pulls a Liferay Portal instance, obtained from the unofficial releases repository.
	If no image tag is passed to the command, the "latest" tag will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("pull nightly requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToPull == "" {
			tagToPull = "latest"
		}

		release := liferay.Release{Tag: tagToPull}

		pullDockerImage(release, forceRemoval)
	},
}

// pullDockerImage uses the image interface to pull it from Docker Hub, removing the cached on if
func pullDockerImage(image liferay.Image, forceRemoval bool) {
	if forceRemoval {
		err := docker.RemoveDockerImage(image.GetFullyQualifiedName())
		if err != nil {
			log.WithFields(log.Fields{
				"image": image.GetFullyQualifiedName(),
			}).Info("The image was not found in the local cache. Skipping removal")
		}
	}

	docker.PullDockerImage(image.GetFullyQualifiedName())
}
