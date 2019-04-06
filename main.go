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

	internal.NewConfig(w)
}

func main() {
	installed := docker.CheckDocker()

	checkWorkspace()

	if !installed {
		log.Fatalln(`Docker is not installed. Please visit "https://docs.docker.com/install/#desktop" to install it before using lpn`)
	}

	cmd.Execute()
}
