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
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	docker "github.com/liferay/lpn/docker"
	liferay "github.com/liferay/lpn/liferay"

	"github.com/spf13/cobra"
)

var filePath string
var directoryPath string

func init() {
	rootCmd.AddCommand(deployCmd)

	subcommands := []*cobra.Command{
		deployCE, deployCommerce, deployDXP, deployNightly, deployRelease}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		deployCmd.AddCommand(subcommand)

		subcommand.Flags().StringVarP(
			&filePath, "files", "f", "",
			`The file or files to deploy. A comma-separated list of files is accepted to deploy
							multiple files at the same time`)

		subcommand.Flags().StringVarP(
			&directoryPath, "dir", "d", "",
			`The directory to deploy its content. Only first-level files will be deployed, so no
							recursive deployment will happen`)
	}
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommandInfo()
	},
}

var deployCE = &cobra.Command{
	Use:   "ce",
	Short: "Deploys files or a directory to Liferay CE's deploy folder in the container run by lpn",
	Long: `Deploys files or a directory to Liferay CE's deploy folder in the container run by lpn.
	The appropriate tag is calculated from the image the container was build from.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArguments()

		ce := liferay.CE{}

		ce.Tag = getTag(ce)

		doDeploy(ce)
	},
}

var deployCommerce = &cobra.Command{
	Use:   "commerce",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArguments()

		commerce := liferay.Commerce{}

		commerce.Tag = getTag(commerce)

		doDeploy(commerce)
	},
}

var deployDXP = &cobra.Command{
	Use:   "dxp",
	Short: "Deploys files or a directory to Liferay DXP's deploy folder in the container run by lpn",
	Long: `Deploys files or a directory to Liferay DXP's deploy folder in the container run by lpn.
	The appropriate tag is calculated from the image the container was build from.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArguments()

		dxp := liferay.DXP{}

		dxp.Tag = getTag(dxp)

		doDeploy(dxp)
	},
}

var deployNightly = &cobra.Command{
	Use:   "nightly",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long:  `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArguments()

		nightly := liferay.Nightly{}

		nightly.Tag = getTag(nightly)

		doDeploy(nightly)
	},
}

var deployRelease = &cobra.Command{
	Use:   "release",
	Short: "Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn",
	Long: `Deploys files or a directory to Liferay Portal's deploy folder in the container run by lpn.
	The appropriate tag is calculated from the image the container was build from.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateArguments()

		release := liferay.Release{}

		release.Tag = getTag(release)

		doDeploy(release)
	},
}

// deployDirectory deploys a directory's content to the running container
func deployDirectory(image liferay.Image, dirPath string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalln("The directory is not valid", err)
	}

	var onlyFilePaths []string

	for _, f := range files {
		if !f.Mode().IsDir() {
			onlyFilePaths = append(onlyFilePaths, path.Join(dirPath, f.Name()))
		}
	}

	deployPaths(image, onlyFilePaths)
}

// deployFiles deploys files to the running container
func deployFiles(image liferay.Image, path string) {
	paths := strings.Split(path, ",")

	deployPaths(image, paths)
}

func deployFile(
	file string, image liferay.Image, resultChannel chan bool,
	errorChannel chan error) {

	if _, err := os.Stat(file); os.IsNotExist(err) {
		select {
		case errorChannel <- err:
			// will break parent goroutine out of loop
		default:
			// don't care, first error wins
		}
		return
	}

	err := docker.CopyFileToContainer(image, file)
	if err != nil {
		select {
		case errorChannel <- err:
			// will break parent goroutine out of loop
		default:
			// don't care, first error wins
		}
		return
	}

	resultChannel <- true
}

// deployPaths deploys files to the running container
func deployPaths(image liferay.Image, paths []string) {
	filesChannel := make(chan string, len(paths))
	for i := range paths {
		filesChannel <- paths[i]
	}
	close(filesChannel)

	workers := 8
	if len(filesChannel) < workers {
		workers = len(filesChannel)
	}

	errorChannel := make(chan error, 1)
	resultChannel := make(chan bool, len(paths))

	for i := 0; i < workers; i++ {
		// Consume work from filesChannel. Loop will end when no more work.
		for file := range filesChannel {
			go deployFile(file, image, resultChannel, errorChannel)
		}
	}

	// Collect results from workers

	for i := 0; i < len(paths); i++ {
		select {
		case <-resultChannel:
			log.Println("[" + paths[i] + "] deployed successfully to " + image.GetDeployFolder())
		case err := <-errorChannel:
			log.Println("Impossible to deploy the file to the container", err)
		}
	}
}

func doDeploy(image liferay.Image) {
	if filePath != "" {
		deployFiles(image, filePath)
	}

	if directoryPath != "" {
		deployDirectory(image, directoryPath)
	}
}

func getTag(image liferay.Image) string {
	imageName, err := docker.GetDockerImageFromRunningContainer(image)
	if err != nil {
		log.Fatalln("The container [" + image.GetContainerName() + "] is NOT running.")
	}

	index := strings.LastIndex(imageName, ":")

	return imageName[index+1 : utf8.RuneCountInString(imageName)]
}

func validateArguments() {
	if filePath == "" && directoryPath == "" {
		log.Fatalln("Please pass a valid path to a file or to a directory as argument")
	}
}
