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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lpn",
	Short: "lpn (Liferay Portal Nook) makes it easier to run Liferay Portal's Docker images.",
	Long: `A Fast and Flexible CLI for managing Liferay Portal's Docker images
				built with love by mdelapenya and friends in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// SubCommandInfo Shows a message for subcommands
func SubCommandInfo() {
	// delegate to subcommands
	fmt.Println(
		"Please run this command adding 'ce', 'commerce', 'dxp', 'nightly' or 'release' " +
			"subcommands.")
}
