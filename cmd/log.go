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
	docker "github.com/liferay/lpn/docker"
	liferay "github.com/liferay/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)

	subcommands := []*cobra.Command{
		logCECmd, logCommerceCmd, logDXPCmd, logNightlyCmd, logReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		logCmd.AddCommand(subcommand)
	}
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays logs for the Liferay Portal instance",
	Long:  `Displays logs for the Liferay Portal instance, identified by [lpn] plus image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var logCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Displays logs for the Liferay Portal CE instance",
	Long:  `Displays logs for the Liferay Portal CE instance, identified by [lpn-ce].`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		logContainer(ce)
	},
}

var logCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Displays logs for the Liferay Portal Commerce instance",
	Long:  `Displays logs for the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		logContainer(commerce)
	},
}

var logDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Displays logs for the Liferay DXP instance",
	Long:  `Displays logs for the Liferay DXP instance, identified by [lpn-dxp].`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		logContainer(dxp)
	},
}

var logNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Displays logs for the Liferay Portal Nightly Build instance",
	Long:  `Displays logs for the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		logContainer(nightly)
	},
}

var logReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Displays logs for the Liferay Portal Release instance",
	Long:  `Displays logs for the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		logContainer(release)
	},
}

// logContainer show the logs for the running container of the specified type
func logContainer(image liferay.Image) {
	docker.LogContainer(image)
}
