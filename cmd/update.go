package cmd

import (
	"fmt"

	"github.com/equinox-io/equinox"
	"github.com/spf13/cobra"
)

const appID = "app_dK5yVpq7ybd"

var publicKey = []byte(`
-----BEGIN ECDSA PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEn7Tuxdoght/89IBx9mphem4LyaI/Wrb3
wbZgy4fLGlagAZsoDK2QtSYRwZTeHf+jRV7adg4IF3DXVkgw3lj92E9HCKrqUKX+
8OIDIF2D2OzuXPJCi/qIFrDWn5jkvhtK
-----END ECDSA PUBLIC KEY-----
`)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates lpn to the latest version",
	Long:  `Updates lpn (Liferay Portal Nook) to the latest version on stable channel`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(
			"Updates from Equinox are disabled. Please go to 'https://releases-lpn.wedeploy.io'" +
				" to download your release")
	},
}

func equinoxUpdate() error {
	var opts equinox.Options
	if err := opts.SetPublicKeyPEM(publicKey); err != nil {
		return err
	}

	// check for the update
	resp, err := equinox.Check(appID, opts)
	switch {
	case err == equinox.NotAvailableErr:
		fmt.Println("No update available, already at the latest version!")
		return nil
	case err != nil:
		fmt.Println("Update failed:", err)
		return err
	}

	// fetch the update and apply it
	err = resp.Apply()
	if err != nil {
		return err
	}

	fmt.Printf("Updated to new version: %s!\n", resp.ReleaseVersion)
	return nil
}
