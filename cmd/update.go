// Copyright (c) 2000-present Liferay, Inc. All rights reserved.
//
// This library is free software; you can redistribute it and/or modify it under
// the terms of the GNU Lesser General Public License as published by the Free
// Software Foundation; either version 2.1 of the License, or (at your option)
// any later version.
//
// This library is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more
// details.

package cmd

import (
	"os"

	v "github.com/mdelapenya/lpn/assets/version"

	"github.com/equinox-io/equinox"
	log "github.com/sirupsen/logrus"
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
		current, _ := v.Asset("VERSION.txt")

		log.WithFields(log.Fields{
			"currentVersion": string(current),
		}).Warn(
			"Updates from Equinox are disabled. Please go to 'https://mdelapenya.github.io/lpn/releases.html'" +
				" to download your release")
		os.Exit(1)
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
		log.Debug("No update available, already at the latest version!")
		return nil
	case err != nil:
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Update failed")
		return err
	}

	// fetch the update and apply it
	err = resp.Apply()
	if err != nil {
		return err
	}

	current, err := v.Asset("VERSION.txt")

	log.WithFields(log.Fields{
		"currentVersion": string(current),
		"newVersion":     resp.ReleaseVersion,
	}).Info("Updated to new version")
	return nil
}
