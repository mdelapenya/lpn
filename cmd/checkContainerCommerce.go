package cmd

import (
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	checkCmd.AddCommand(checkContainerCommerceCmd)
}

var checkContainerCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Checks if there is a Commerce container created by lpn",
	Long: `Checks if there is a Commerce container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Commerce container with name [lpn-commere] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		CheckDockerContainerExists(commerce)
	},
}
