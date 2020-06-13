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
	rootCmd.AddCommand(checkCmd)

	subcommands := []*cobra.Command{
		checkContainerCECmd, checkContainerCommerceCmd, checkContainerDXPCmd,
		checkContainerNightlyCmd, checkContainerReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		checkCmd.AddCommand(subcommand)
	}
}

var checkCmd = &cobra.Command{
	Use:   "checkContainer",
	Short: "Checks if there is a container created by lpn",
	Long: `Checks if there is a container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a container with name "lpn" plus image type created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var checkContainerCECmd = &cobra.Command{
	Use:   "ce",
	Short: "Checks if there is a CE container created by lpn",
	Long: `Checks if there is a CE container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a CE container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		ce := liferay.CE{}

		checkDockerContainerExists(ce)
	},
}

var checkContainerCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Checks if there is a Commerce container created by lpn",
	Long: `Checks if there is a Commerce container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Commerce container with name [lpn-commerce] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		checkDockerContainerExists(commerce)
	},
}

var checkContainerDXPCmd = &cobra.Command{
	Use:   "dxp",
	Short: "Checks if there is a DXP container created by lpn",
	Long: `Checks if there is a DXP container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a DXP container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		dxp := liferay.DXP{}

		checkDockerContainerExists(dxp)
	},
}

var checkContainerNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Checks if there is a Nightly Build container created by lpn",
	Long: `Checks if there is a Nightly Build container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Nightly Build container with name [lpn-nightly] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		checkDockerContainerExists(nightly)
	},
}

var checkContainerReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Checks if there is a Release container created by lpn",
	Long: `Checks if there is a Release container created by lpn (Liferay Portal Nook).
	Uses docker container inspect to check if there is a Release container with name [lpn-release] created by lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		checkDockerContainerExists(release)
	},
}

// checkDockerContainerExists removes the running container
func checkDockerContainerExists(image liferay.Image) {
	exists := docker.CheckDockerContainerExists(image.GetContainerName())

	if !exists {
		log.Fatalln("The container [" + image.GetContainerName() + "] does NOT exist in the system.")
	}

	log.Println("The container [" + image.GetContainerName() + "] DOES exist in the system.")
}
