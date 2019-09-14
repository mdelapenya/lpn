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
