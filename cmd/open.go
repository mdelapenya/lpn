package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(openCmd)

	subcommands := []*cobra.Command{
		openCECmd, openCommerceCmd, openDXPCmd, openNightlyCmd, openReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		subcommand.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs commands with Debug log level")
		subcommand.VisitParents(addVerboseFlag)

		openCmd.AddCommand(subcommand)
	}
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Opens a browser with the Liferay Portal nook instance",
	Long:  `Opens a browser with the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		enableDebugLevel()
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var openCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Opens a browser with  the Liferay Portal CE instance",
	Long:  `Opens a browser with  the Liferay Portal CE instance, identified by [lpn-ce].`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		openBrowser(ce)
	},
}

var openCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Opens a browser with  the Liferay Portal Commerce instance",
	Long:  `Opens a browser with  the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		openBrowser(commerce)
	},
}

var openDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Opens a browser with  the Liferay DXP instance",
	Long:  `Opens a browser with  the Liferay DXP instance, identified by [lpn-dxp].`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		openBrowser(dxp)
	},
}

var openNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Opens a browser with  the Liferay Portal Nightly Build instance",
	Long:  `Opens a browser with  the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		openBrowser(nightly)
	},
}

var openReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Opens a browser with the Liferay Portal Release instance",
	Long:  `Opens a browser with  the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		openBrowser(release)
	},
}

// openBrowser opens a browser the running container
func openBrowser(image liferay.Image) {
	url := "http://localhost:" + docker.GetTomcatPort(image)

	log.WithFields(log.Fields{
		"url": url,
	}).Debug("Opening portal URL")

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}
