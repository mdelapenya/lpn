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
	"log"

	docker "github.com/liferay/lpn/docker"
	liferay "github.com/liferay/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)

	subcommands := []*cobra.Command{
		startCECmd, startCommerceCmd, startDXPCmd, startNightlyCmd, startReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		startCmd.AddCommand(subcommand)
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Liferay Portal nook instance",
	Long:  `Starts the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var startCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Starts the Liferay Portal CE instance",
	Long:  `Starts the Liferay Portal CE instance, identified by [lpn-ce].`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		startDockerContainer(ce)
	},
}

var startCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Starts the Liferay Portal Commerce instance",
	Long:  `Starts the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		startDockerContainer(commerce)
	},
}

var startDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Starts the Liferay DXP instance",
	Long:  `Starts the Liferay DXP instance, identified by [lpn-dxp].`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		startDockerContainer(dxp)
	},
}

var startNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Starts the Liferay Portal Nightly Build instance",
	Long:  `Starts the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		startDockerContainer(nightly)
	},
}

var startReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Starts the Liferay Portal Release instance",
	Long:  `Starts the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		startDockerContainer(release)
	},
}

// startDockerContainer starts the stopped container
func startDockerContainer(image liferay.Image) {
	err := docker.StartDockerContainer(image)
	if err != nil {
		log.Fatalln("Impossible to start the container ["+image.GetContainerName()+"]", err)
	}
}
