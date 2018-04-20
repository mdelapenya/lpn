package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	stopCmd.AddCommand(stopCommerceCmd)
}

var stopCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Stops the Liferay Portal Commerce instance",
	Long:  `Stops the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		StopDockerContainer(commerce)
	},
}
