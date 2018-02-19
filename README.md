# lpn (Liferay Portal Nook)

This Golang CLI makes it easier to run Liferay Portal's Docker images.

It wraps Docker commands so you just have to run this tool, and pass the specific tag you want to run.

## Install

Install this tool downloading it from our [stable release channel](https://dl.equinox.io/mdelapenya/lpn/stable).

### Brew (Mac Users)

```shell
$ brew tap eqnxio/mdelapenya
$ brew install lpn
```

## Requirements

You have to install Docker on your machine first. Check [this guide](https://docs.docker.com/install).

## What does it do?

```shell
A Fast and Flexible CLI for managing Liferay Portal's Docker images
				built with love by mdelapenya and friends in Go.

Usage:
  lpn [flags]
  lpn [command]

Available Commands:
  checkContainer Checks if there is a container created by lpn (Liferay Portal Nook)
  checkImage     Checks if the proper Liferay Portal image has been pulled by lpn (Liferay Portal Nook)
  help           Help about any command
  log            Displays logs for the Liferay Portal instance
  pull           Pulls a Liferay Portal Docker image
  rm             Removes the Liferay Portal instance
  run            Runs a Liferay Portal instance
  update         Updates lpn (Liferay Portal Nook) to the latest version
  version        Print the version number of lpn (Liferay Portal Nook)

Flags:
  -h, --help   help for lpn

Use "lpn [command] --help" for more information about a command.
```

## Commands

### checkContainer

Checks if there is a container created by lpn (Liferay Portal Nook).

Uses `docker container inspect` to check if there is a container with name `liferay-portal-nook` created by lpn (Liferay Portal Nook).

### checkImage

```shell
$ lpn checkImage
Checks if the proper Liferay Portal image has been pulled by lpn.
	Uses "docker image inspect" to check if the proper Liferay Portal image has
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.

Usage:
  lpn checkImage [flags]
  lpn checkImage [command]

Available Commands:
  nightly     Check if the proper Liferay Portal Nightly Build image has been pulled by lpn
  release     Check if the proper Liferay Portal release image has been pulled by lpn

Flags:
  -h, --help   help for checkImage

Use "lpn checkImage [command] --help" for more information about a command.
```

#### checkImage nightly

```shell
$ lpn checkImage nightly
Checks if the proper Liferay Portal Nightly Build image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal image has
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.

Usage:
  lpn checkImage nightly [flags]

Flags:
  -h, --help         help for nightly
  -t, --tag string   Sets the image tag to check (default "latest")
```

#### checkImage release

```shell
Check if the proper Liferay Portal release image has been pulled by lpn.
	Uses docker image inspect to check if the proper Liferay Portal image has
	been pulled by lpn (Liferay Portal Nook). If no image tag is passed to the command,
	the tag "latest" will be used.

Usage:
  lpn checkImage release [flags]

Flags:
  -h, --help         help for release
  -t, --tag string   Sets the image tag to check (default "latest")
```

### help

Help about any command.

### log

Displays logs for the Liferay Portal instance, identified by [`liferay-portal-nook`].

### pull

```shell
$ lpn pull -h
Pulls a Liferay Portal Docker image from the unofficial repositories "mdelapenya/liferay-portal" and "mdelapenya/liferay-portal-nightlies".
	For that, please run this command adding "release" or "nightly" subcommands.
	If no image tag is passed to the command, the tag representing the current date [yyyyMMdd] will be used.

Usage:
  lpn pull [flags]
  lpn pull [command]

Available Commands:
  nightly        Pulls a Liferay Portal Docker image from Nightly Builds
  release        Pulls a Liferay Portal Docker image from releases

Flags:
  -h, --help   help for pull

Use "lpn pull [command] --help" for more information about a command.
```

#### pull nightly

```shell
$ lpn pull nightly -h
Pulls a Liferay Portal Docker image from the Nighlty Builds repository: "mdelapenya/liferay-portal-nightlies".
 If no image tag is passed to the command, the tag representing the current date [20180219] will be used.

Usage:
  lpn pull nightly [flags]

Flags:
  -h, --help   help for nightly
```

#### pull release

```shell
$ lpn pull release -h
Pulls a Liferay Portal instance, obtained from the unofficial releases repository: "mdelapenya/liferay-portal".
	If no image tag is passed to the command, the "latest" tag will be used.

Usage:
  lpn pull release [flags]

Flags:
  -h, --help   help for release
```

### rm

Removes the Liferay Portal nook instance, identified by [`liferay-portal-nook`].

### run

```shell
$ lpn run -h
Runs a Liferay Portal instance, obtained from the unofficial repositories: `mdelapenya/liferay-portal` or `mdelapenya/liferay-portal-nightlies`. For that, please run this command adding `release` or `nightly` subcommands. If no image tag is passed to the subcommand, the tag representing the current date [`yyyyMMdd`] will be used.

Usage:
  lpn run [flags]
  lpn run [command]

Available Commands:
  nightly     Runs a Liferay Portal instance from Nightly Builds
  release     Runs a Liferay Portal instance from releases

Flags:
  -h, --help            help for run
```

#### run nightly

```shell
$ lpn run nightly -h
Runs a Liferay Portal instance, obtained from Nightly Builds repository: mdelapenya/liferay-portal-nightlies.
	If no image tag is passed to the command, the tag representing the current date [yyyyMMdd] will be used.

Usage:
  lpn run nightly [flags]

Flags:
  -d, --debug           Enables debug mode. (default false)
  -D, --debugPort int   Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled (default 9000)
  -h, --help            help for nightly
  -p, --httpPort int    Sets the HTTP port of Liferay Portal's bundle. (default 8080)
```

#### run release

```shell
$ lpn run release -h
Runs a Liferay Portal instance, obtained from the unofficial releases repository: mdelapenya/liferay-portal.
	If no image tag is passed to the command, the "latest" tag will be used.

Usage:
  lpn run release [flags]

Flags:
  -d, --debug           Enables debug mode. (default false)
  -D, --debugPort int   Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled (default 9000)
  -h, --help            help for release
  -p, --httpPort int    Sets the HTTP port of Liferay Portal's bundle. (default 8080)
```

To achieve that:

- The tool will ask you to type an image tag from Liferay Portal's Docker images (check available tags [here](https://hub.docker.com/r/mdelapenya/liferay-portal/tags/) for releases, and [here](https://hub.docker.com/r/mdelapenya/liferay-portal-nightlies/tags/) for nightly builds.
  - If no tag is provided, then it will use current date as tag, i.e. `20180214`.
- It downloads the Docker image to the local engine.
- It checks whether the Docker container this tool spins up is running. In that case, the tool deletes it.
- It spins up a Docker container, using the port configured for Tomcat, and 11311 for OSGi console. The name of the container will be `liferay-portal-nook`. Once started, please open a web browser in [https://localhost:8080](http://localhost:8080) to check the portal.

### update

Updates lpn (Liferay Portal Nook) to the latest version on stable channel.

### version

All software has versions. This is lpn (Liferay Portal Nook).