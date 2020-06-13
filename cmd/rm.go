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

	docker "github.com/mdelapenya/lpn/docker"
	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)

	subcommands := []*cobra.Command{rmCECmd, rmCommerceCmd, rmDXPCmd, rmNightlyCmd, rmReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		rmCmd.AddCommand(subcommand)
	}
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes the Liferay Portal nook instance",
	Long:  `Removes the Liferay Portal nook instance, identified by [lpn] plus each image type.`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var rmCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Removes the Liferay Portal CE instance",
	Long:  `Removes the Liferay Portal CE instance, identified by [lpn-ce].`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		removeDockerContainer(ce)
	},
}

var rmCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Removes the Liferay Portal Commerce instance",
	Long:  `Removes the Liferay Portal Commerce instance, identified by [lpn-commerce].`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		removeDockerContainer(commerce)
	},
}

var rmDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Removes the Liferay DXP instance",
	Long:  `Removes the Liferay DXP instance, identified by [lpn-dxp].`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		removeDockerContainer(dxp)
	},
}

var rmNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Removes the Liferay Portal Nightly Build instance",
	Long:  `Removes the Liferay Portal Nightly Build instance, identified by [lpn-nightly].`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		removeDockerContainer(nightly)
	},
}

var rmReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Removes the Liferay Portal Release instance",
	Long:  `Removes the Liferay Portal Release instance, identified by [lpn-release].`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		removeDockerContainer(release)
	},
}

// removeDockerContainer removes the running container
func removeDockerContainer(image liferay.Image) {
	err := docker.RemoveDockerContainer(image)
	if err != nil {
		log.Fatalln("Impossible to remove the container ["+image.GetContainerName()+"]", err)
	}
}
