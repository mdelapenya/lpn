# lpn (Liferay Portal Nook)

[![Build Status](https://travis-ci.org/mdelapenya/lpn.svg?branch=master)](https://travis-ci.org/mdelapenya/lpn)
[![Codecov Coverage](https://codecov.io/gh/mdelapenya/lpn/branch/master/graph/badge.svg)](https://codecov.io/gh/mdelapenya/lpn)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdelapenya/lpn)](https://goreportcard.com/report/github.com/mdelapenya/lpn)

This Golang CLI makes it easier to run Liferay Portal's Docker images.

It wraps Docker commands so you just have to run this tool, and pass the specific tag you want to run.

## Requirements

You have to install Docker on your machine first. Check [this guide](https://docs.docker.com/install).

## Installation

For current stable version of `lpn`, please visit [downloads page](https://lpn.lfr.io/releases.html), and select the target platform, based on O.S. and architecture.

If you are in the Golang world, you could install this tool from source code.

For that reason you need to:

- [Install Golang runtime](https://golang.org/doc/install)
- Make sure that you have `GOPATH` environment variable defined to the location you want to have your Golang projects
- Add `$GOPATH\bin` to your `$PATH`, like this: `export PATH=${PATH}:${GOPATH//://bin:}/bin`
- Clone this repo under `$GOPATH`: `git clone https://github.com/mdelapenya/lpn $GOPATH/src/github.com/mdelapenya/lpn`
- From inside the project, run `go install`

Now you can use `lpn` from your command line.

## What does lpn do?

With `lpn` you'll be able to:

- Run Liferay Portal containers using the released version you prefer:
  - Liferay Portal official images, obtained from [here for CE](https://hub.docker.com/r/liferay/portal/tags/) and [here for DXP](https://hub.docker.com/r/liferay/dxp/tags/).
  - Liferay Portal nightly builds, obtained from [here](https://hub.docker.com/r/liferay/portal-snapshot/tags/).
  - Liferay Portal releases, obtained from [here](https://hub.docker.com/r/mdelapenya/liferay-portal/tags/).
  - Liferay Portal including specific products, like Commerce.
- Deploy applications to the running containers:
  - Imagine a developer gives you a non-released JAR/WAR file to test it. You could deploy and test it in seconds!
- Check logs of the running containers
- And more!

## Usage

The available capabilities present in the tool are the following:

- Run a container from the desired Liferay Portal/DXP image.
- Configure a Liferay Portal/DXP container to be run using a custom portal-ext configuration file.
- Deploy a file or the content of a directory to the deploy folder of a Liferay Portal/DXP running container.
- Display logs of a Liferay Portal/DXP running container.
- Pull Liferay Portal/DXP images on demand.
- List the available tags to pull from the Docker Hub repository of a Liferay Portal/DXP image.
- Check if a Liferay Portal/DXP image was already pulled.
- Check if a container of the desired Liferay Portal/DXP image is running.
- Stop a Liferay Portal/DXP running container.
- Remove a Liferay Portal/DXP running container.
- Remove a Liferay Portal/DXP image from your local Docker installation.
- Open a Liferay Portal/DXP running container in the default browser.

### Which are the available commands?

Execute `lpn help` to see the list of available commands. Each command could have subcommands, so append the `help` command to each subcommand and you'll get a list of options for each one.

This could be an example of how the `help` subcommand shows:

```shell
$ lpn help
A Fast and Flexible CLI for managing Liferay Portal's Docker images
				built with love by mdelapenya and friends in Go.

Usage:
  lpn [flags]
  lpn [command]

Available Commands:
  checkContainer Checks if there is a container created by lpn (Liferay Portal Nook)
  checkImage     Checks if the proper Liferay Portal image has been pulled by lpn (Liferay Portal Nook)
  deploy         Deploys a file to Liferay Portal's deploy folder in the container run by lpn
  help           Help about any command
  license        Print the license of lpn
  log            Displays logs for the Liferay Portal instance
  pull           Pulls a Liferay Portal Docker image
  rm             Removes the Liferay Portal instance
  run            Runs a Liferay Portal instance
  start          Starts the Liferay Portal instance
  stop           Stops the Liferay Portal instance
  update         Updates lpn (Liferay Portal Nook) to the latest version
  version        Print the version number of lpn (Liferay Portal Nook)

Flags:
  -h, --help   help for lpn

Use "lpn [command] --help" for more information about a command.
```

Once you have typed the proper command, to specify with which image type you want to execute the command, there are the following subcommands:

- ce
- dxp
- release
- nightly
- commerce

So any command needs the combination of one of the subcommands above. So to run a DXP image, you would need to execute `lpn run dxp`.

## Running a container from a Liferay Portal/DXP image

It will run the desired image, pulling it first if it does not exist in your local Docker installation. To specify which image type you want to run, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

You will be able to configure in which state you want to run the image, using the following flags:

| Flag | Description |
|:-|:-|
| ` -d, --debug` | Enables debug mode. (default false) |
| ` -D, --debugPort` | Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled (default 9000) |
| ` -g, --gogoPort` | Sets the GoGo Shell port of Liferay Portal's bundle. (default 11311) |
| ` -p, --httpPort` | Sets the HTTP port of Liferay Portal's bundle. (default 8080) |
| ` -m, --memory` | Sets the memory for the JVM memory configuration of Liferay Portal's bundle. (default "-Xmx2048m" in the CE and DXP images, and "2048m" in the rest) |
| ` -P, --properties` | Sets the location of a portal-ext properties files to configure the running instance of Liferay Portal's bundle. |
| ` -t, --tag` | Sets the image tag to run |

Examples:
```shell
$ lpn run ce -t "7.1.1-ga2"
$ lpn run dxp --properties "/tmp/portal-ext.properties"
$ lpn run nightly
$ lpn run commerce --debug --httpPort 8081 ---memory "Xmx8g"
```

## Copying files to the deploy folder

It will deploy a file, or the first-level content of a directory, to the deploy folder of a running container, pulling it first if it does not exist in your local Docker installation. To specify to which image type you want to deploy, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

You will be able to configure which file or directory you want to deploy to the running container, using the following flags:

| Flag | Description |
|:-|:-|
| ` -d, --dir` | The directory to deploy its content. Only first-level files will be deployed, so no recursive deployment will happen |
| ` -f, --files` | The file or files to deploy. A comma-separated list of files is accepted to deploy multiple files at the same time |

Examples:
```shell
$ lpn deploy ce --dir /tmp/modules-from-my-dev-team
$ lpn deploy nightly --files /tmp/moduleA.jar
$ lpn deploy commerce --files /tmp/moduleA.jar,/tmp/themeB.war
```

## Displaying logs

It will display the logs of a running container, reading each log line in a _tail_ mode. In this case this log corresponds to the Tomcat's log file. To specify to which image type you want to show logs, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

This command does not accept any flag to configure its execution.

Examples:
```shell
$ lpn log ce
$ lpn log dxp
$ lpn log release
$ lpn log nightly
$ lpn log commerce
```

## Pulling Liferay images

It will pull a desired image from Docker Hub to your local Docker installation. To specify to which image type you want to pull, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

You will be able to configure which image you are going to pull using the following flags:

| Flag | Description |
|:-|:-|
| ` -f, --forceRemoval` | Removes the cached, local image, if exists |
| ` -t, --tag` | Sets the image tag to pull |

Depending on the image type, the default value for `--tag` flag would be:

- For CE: `7.0.6-ga7`
- For DXP: `7.0.10.8`
- For Releases: `latest`
- For Nightly Builds and Commerce: Current date in the `20181128` format.

Examples:
```shell
$ lpn pull ce
$ lpn pull release --forceRemoval
$ lpn pull commerce --tag "20181026"
```

## Listing the available Liferay images

It will list the existing tags on Docker Hub for the desired images type. To specify to which image type you want to list its tags, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

You will be able to browse the list of available tags using the following flags:

| Flag | Description |
|:-|:-|
| ` -p, --page` | Sets the page element where tags exist (default 1) |
| ` -size, --size` | Sets the number of tags to retrieve per page (default 25) |

Examples:
```shell
$ lpn tags ce
$ lpn tags dxp --page 2 --size 5
$ lpn tags release
$ lpn tags nightly -p 2 -s 5
$ lpn tags commerce
```

## Checking if an image is present in the local Docker installation

It will check if the desired images type exists in your local Docker installation. To specify to which image type you want to check, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

You will be able to configure which image you are going to check using the following flags:

| Flag | Description |
|:-|:-|
| ` -t, --tag` | Sets the image tag to check |

Depending on the image type, the default value for `--tag` flag would be:

- For CE: `7.0.6-ga7`
- For DXP: `7.0.10.8`
- For Releases: `latest`
- For Nightly Builds and Commerce: Current date in the `20181128` format.

Examples:
```shell
$ lpn checkImage ce --tag "6.1.2-ga3"
$ lpn checkImage dxp --tag "7.0.10.8"
$ lpn checkImage release
$ lpn checkImage nightly --tag "20181026"
$ lpn checkImage commerce
```

## Checking if a container is running

It will check if there is a running container for the desired images type. To specify to which image type you want to check its container, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

This command does not accept any flag to configure its execution.

Examples:
```shell
$ lpn checkContainer ce
$ lpn checkContainer dxp
$ lpn checkContainer release
$ lpn checkContainer nightly
$ lpn checkContainer commerce
```

## Starting a stopped container

It will start an already stopped container, if it exists. To specify to which image type you want to start its container, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

This command does not accept any flag to configure its execution.

Examples:
```shell
$ lpn start ce
$ lpn start dxp
$ lpn start release
$ lpn start nightly
$ lpn start commerce
```

## Stopping a running container

It will stop a running container, if it exists. To specify to which image type you want to stop its container, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

This command does not accept any flag to configure its execution.

Examples:
```shell
$ lpn stop ce
$ lpn stop dxp
$ lpn stop release
$ lpn stop nightly
$ lpn stop commerce
```

## Removing a running container

It will remove a running container, if it exists. To specify to which image type you want to remove its container, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

This command does not accept any flag to configure its execution.

Examples:
```shell
$ lpn rm ce
$ lpn rm dxp
$ lpn rm release
$ lpn rm nightly
$ lpn rm commerce
```

## Removing an image from your local Docker installation

It will remove the desired images type from your local Docker installation. To specify to which image type you want to remove, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

You will be able to configure which image you are going to remove using the following flags:

| Flag | Description |
|:-|:-|
| ` -t, --tag` | Sets the image tag to remove |

Depending on the image type, the default value for `--tag` flag would be:

- For CE: `7.0.6-ga7`
- For DXP: `7.0.10.8`
- For Releases: `latest`
- For Nightly Builds and Commerce: Current date in the `20181128` format.

Examples:
```shell
$ lpn rmi ce --tag "6.1.2-ga3"
$ lpn rmi dxp --tag "7.0.10.8"
$ lpn rmi release
$ lpn rmi nightly --tag "20181026"
$ lpn rmi commerce
```

## Showing the license

It will display the license of the tool. It's using BSD-3 license, but we are in the process of deciding which one to use.

```shell
$ lpn license
```

## Showing tool's current version

It will display the version of the tool, including the version of the required runtime dependencies (Docker).

```shell
$ lpn version
lpn (Liferay Portal Nook) v0.7.1 -- HEAD
Docker version:
Client version: 1.25
Server version: 18.06.1-ce
Go version: go1.10.3
```

It shows `Go version` because is used by Docker.

## Opening the running instance from a browser

It will open the O.S. default browser with the home page of the running instance of the desired images type. To specify to which image type you want to open in the browser, please select it adding the `ce`, `dxp`, `release`, `nightly`, `commerce` subcommands.

This command does not accept any flag to configure its execution.

Examples:
```shell
$ lpn open ce
$ lpn open dxp
$ lpn open release
$ lpn open nightly
$ lpn open commerce
```
