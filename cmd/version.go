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
