package cmd

import (
	v "github.com/mdelapenya/lpn/assets/version"
	docker "github.com/mdelapenya/lpn/docker"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of lpn",
	Long:  `All software has versions. This is lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := v.Asset("VERSION.txt")

		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"command": cmd.Use,
			}).Error("Error executing Command")
			return
		}

		dockerClientVersion, dockerServerVersion, _ := docker.GetDockerVersion()

		log.WithFields(log.Fields{
			"lpn":          string(version),
			"dockerClient": dockerClientVersion,
			"dockerServer": dockerServerVersion.Version,
			"golang":       dockerServerVersion.GoVersion,
		}).Infof("lpn (Liferay Portal Nook) v%s -- HEAD", version)
	},
}
