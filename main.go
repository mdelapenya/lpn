package main

import "fmt"

const dockerImage = "mdelapenya/liferay-portal-nightlies"

func main() {
	fmt.Print("Enter the Image Tag you want to use for [" + dockerImage + "]: ")
	var imageTag string

	fmt.Scanf("%s", &imageTag)

	var image = dockerImage + ":" + imageTag

	fmt.Println(image)
}
