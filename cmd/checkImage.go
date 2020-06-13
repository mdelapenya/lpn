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
	"errors"
	"log"

	date "github.com/liferay/lpn/date"
	docker "github.com/liferay/lpn/docker"
	internal "github.com/liferay/lpn/internal"
	liferay "github.com/liferay/lpn/liferay"

	"github.com/spf13/cobra"
)

var tagToCheck string

func init() {
	rootCmd.AddCommand(checkImageCmd)

	subcommands := []*cobra.Command{
		checkImageCE, checkImageCommerce, checkImageDXP, checkImageNightly, checkImageRelease}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		checkImageCmd.AddCommand(subcommand)

		subcommand.Flags().StringVarP(&tagToCheck, "tag", "t", "", "Sets the image tag to check")
	}
}

var checkImageCmd = &cobra.Command{
	Use:   "checkImage",
	Short: "Checks if the proper Liferay Portal image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal image has been pulled by lpn.
	Uses "docker image inspect" to check if the proper Liferay Portal image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var checkImageCE = &cobra.Command{
	Use:   "ce",
	Short: "Checks if the proper Liferay Portal CE image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal CE image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the default tag (see configuration file) will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage release requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToCheck == "" {
			tagToCheck = internal.LpnConfig.GetPortalImageTag("ce")
		}

		ce := liferay.CE{Tag: tagToCheck}

		checkImage(ce)
	},
}

var checkImageCommerce = &cobra.Command{
	Use:   "commerce",
	Short: "Checks if the proper Liferay Portal with Commerce Build image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal with Commerce Build image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal with Commerce image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage commerce requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToCheck == "" {
			tagToCheck = internal.LpnConfig.GetPortalImageTag("commerce")
		}

		commerce := liferay.Commerce{Tag: tagToCheck}

		checkImage(commerce)
	},
}

var checkImageDXP = &cobra.Command{
	Use:   "dxp",
	Short: "Checks if the proper Liferay DXP image has been pulled by lpn",
	Long: `Checks if the proper Liferay DXP image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay DXP image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the default tag (see configuration file) will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage release requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToCheck == "" {
			tagToCheck = internal.LpnConfig.GetPortalImageTag("dxp")
		}

		dxp := liferay.DXP{Tag: tagToCheck}

		checkImage(dxp)
	},
}

var checkImageNightly = &cobra.Command{
	Use:   "nightly",
	Short: "Checks if the proper Liferay Portal Nightly Build image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal Nightly Build image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "` + date.CurrentDate + `" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage nightly requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToCheck == "" {
			tagToCheck = date.CurrentDate
		}

		nightly := liferay.Nightly{Tag: tagToCheck}

		checkImage(nightly)
	},
}

var checkImageRelease = &cobra.Command{
	Use:   "release",
	Short: "Checks if the proper Liferay Portal release image has been pulled by lpn",
	Long: `Checks if the proper Liferay Portal release image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal image has 
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("checkImage release requires zero or one argument representing the image tag")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tagToCheck == "" {
			tagToCheck = "latest"
		}

		release := liferay.Release{Tag: tagToCheck}

		checkImage(release)
	},
}

// checkImage uses the image interface to check if it exists
func checkImage(image liferay.Image) {
	exists := docker.CheckDockerImageExists(image.GetFullyQualifiedName())

	if exists == false {
		log.Fatalln("The image [" + image.GetFullyQualifiedName() + "] has NOT been pulled from Docker Hub.")
	}

	log.Println("The image [" + image.GetFullyQualifiedName() + "] has been pulled from Docker Hub.")
}
