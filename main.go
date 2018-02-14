package main

import (
	"bufio"
	"fmt"
	docker "liferay-gowl/docker"
	"os"
	"time"

	"github.com/vjeantet/jodaTime"
)

func main() {
	var currentDate = jodaTime.Format("YYYYMMdd", time.Now())

	fmt.Println("Enter the Image Tag you want to use for [" + docker.DockerImage + "]")
	fmt.Print("If you leave it empty, we will use [" + currentDate + "]: ")

	var imageTag string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		imageTag = scanner.Text()
	}

	if imageTag == "" {
		imageTag = currentDate
	}

	docker.RunDockerImage(getDockerImage(imageTag))
}

func getDockerImage(imageTag string) string {
	return docker.DockerImage + ":" + imageTag
}
