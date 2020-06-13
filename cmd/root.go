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
	internal "github.com/mdelapenya/lpn/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "lpn",
	Short: "lpn (Liferay Portal Nook) makes it easier to run Liferay Portal's Docker images.",
	Long: `A Fast and Flexible CLI for managing Liferay Portal's Docker images
				built with ❤️ by mdelapenya and friends in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		enableDebugLevel()
	},
}

func addVerboseFlag(c *cobra.Command) {
	if c.Flag("verbose") == nil {
		c.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Runs command with Debug log level")
	}
}

func enableDebugLevel() {
	if verbose {
		internal.ConfigureLogger("DEBUG")
	}
	log.Debug("Debug logger activated")
}

// Execute execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// SubCommandInfo Shows a message for subcommands
func SubCommandInfo() {
	// delegate to subcommands
	log.Warn(
		"Please run this command adding 'ce', 'commerce', 'dxp', 'nightly' or 'release' " +
			"subcommands.")
}
