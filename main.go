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

	"github.com/mdelapenya/lpn/cmd"

	docker "github.com/mdelapenya/lpn/docker"
	internal "github.com/mdelapenya/lpn/internal"
)

func init() {
	internal.CheckWorkspace()

	installed := docker.CheckDocker()

	if !installed {
		log.Fatalln(`Docker is not installed. Please visit "https://docs.docker.com/install/#desktop" to install it before using lpn`)
	}
}

func main() {
	cmd.Execute()
}
