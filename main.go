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

package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/mdelapenya/lpn/cmd"

	docker "github.com/mdelapenya/lpn/docker"
	internal "github.com/mdelapenya/lpn/internal"
)

func init() {
	checkWorkspace()

	installed := docker.CheckDocker()

	if !installed {
		log.Fatalln(`Docker is not installed. Please visit "https://docs.docker.com/install/#desktop" to install it before using lpn`)
	}
}

// checkWorkspace creates this tool workspace under user's home, in a hidden directory named ".wt"
func checkWorkspace() {
	usr, _ := user.Current()

	w := filepath.Join(usr.HomeDir, ".lpn")

	if _, err := os.Stat(w); os.IsNotExist(err) {
		err = os.MkdirAll(w, 0755)
		if err != nil {
			log.Fatalf("Cannot create workdir for LPN at "+w, err)
		}

		log.Println("lpn workdir created at " + w)
	}

	internal.LpnWorkspace = w

	internal.LpnConfig = internal.NewConfig(w)
}

func main() {
	cmd.Execute()
}
