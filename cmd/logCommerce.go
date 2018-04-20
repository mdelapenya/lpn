package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	logCmd.AddCommand(logCommerceCmd)
}

var logCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Displays logs for the Liferay Portal Commerce instance",
	Long:  `Displays logs for the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		LogDockerContainer(commerce)
	},
}
