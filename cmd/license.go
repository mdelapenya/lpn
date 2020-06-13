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

	license "github.com/mdelapenya/lpn/assets/license"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(licenseCmd)
}

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Print the license of lpn",
	Long:  `All software has a license. This is lpn (Liferay Portal Nook)`,
	Run: func(cmd *cobra.Command, args []string) {
		licenseFile, err := license.Asset("LICENSE.txt")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(string(licenseFile))
	},
}
