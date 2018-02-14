package main

import "fmt"
import docker "liferay-gowl/docker"

func main() {
	fmt.Print("Enter the Image Tag you want to use for [" + docker.DockerImage + "]: ")
	var imageTag string

	fmt.Scanf("%s", &imageTag)

	docker.DownloadDockerImage(getDockerImage(imageTag))
}

func getDockerImage(imageTag string) string {
	return docker.DockerImage + ":" + imageTag
}
