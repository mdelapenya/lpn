package main

import (
	"log"

	"github.com/mdelapenya/lpn/cmd"

	docker "github.com/mdelapenya/lpn/docker"
)

func main() {
	installed := docker.CheckDocker()

	if !installed {
		log.Fatalln(`Docker is not installed. Please visit "https://docs.docker.com/install/#desktop" to install it before using lpn`)
	}

	cmd.Execute()
}
