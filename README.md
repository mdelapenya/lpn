# lpn (Liferay Portal Nightly)

This Golang CLI makes it easier to run Liferay Portal's nightly builds.

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
A Fast and Flexible CLI for managing Liferay Portal's nightly builds
				built with love by mdelapenya and friends in Go.

Usage:
  lpn [flags]
  lpn [command]

Available Commands:
  checkContainer Checks if there is a container created by lpn (Liferay Portal Nightly)
  checkImage     Checks if the proper Liferay Portal nightly image has been pulled by lpn (Liferay Portal Nightly)
  help           Help about any command
  log            Displays logs for the Liferay Portal nightly instance
  pull           Pulls a Liferay Portal nightly Docker image
  rm             Removes the Liferay Portal nightly instance
  run            Runs a Liferay Portal nightly instance
  update         Updates lpn (Liferay Portal Nightly) to the latest version
  version        Print the version number of lpn (Liferay Portal Nightly)

Flags:
  -h, --help   help for lpn

Use "lpn [command] --help" for more information about a command.
```

## Commands

### checkContainer

Uses `docker container inspect` to check if there is a container with name `liferay-portal-nightly` created by lpn (Liferay Portal Nightly).

### checkImage

Checks if the proper Liferay Portal nightly image has been pulled by lpn (Liferay Portal Nightly).

Uses `docker image inspect` to check if the proper Liferay Portal nightly image has been pulled by lpn (Liferay Portal Nightly). If no image tag is passed to the command, the tag `latest` will be used.

This command accepts following flags:

| Flag | Description | Default value |
|------|-------------|---------------|
|`-t`, `--tag`| Sets the image tag to check.| latest|

### help

Help about any command.

### log

Displays logs for the Liferay Portal nightly instance, identified by [`liferay-portal-nightly`].

### pull

Pulls a Liferay Portal nightly Docker image from `mdelapenya/liferay-portal-nightlies` repository. If no image tag is passed to the command, the tag representing the current date will be used.

### rm

Removes the Liferay Portal nightly instance, identified by [`liferay-portal-nightly`].

### run

```shell
lpn run -h
Runs a Liferay Portal nightly instance, obtained from mdelapenya/liferay-portal-nightlies.
	If no image tag is passed to the command, the tag representing the current date [20180215]
	will be used.

Usage:
  lpn run [flags]

Flags:
  -d, --debug           Enables debug mode on port 9000
  -D, --debugPort int   Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled (default 9000)
  -h, --help            help for run
  -p, --httpPort int    Sets the HTTP port of Liferay Portal's bundle. (default 8080)
```

Runs a Liferay Portal nightly instance, obtained from `mdelapenya/liferay-portal-nightlies`. If no image tag is passed to the command, the tag representing the current date [`liferay-portal-nightly`] will be used.

This command accepts following flags:

| Flag | Description | Default value |
|------|-------------|---------------|
|`-d`, `--debug`| Enables debug mode.| false |
|`-D`, `--debugPort` | Sets the debug port of Liferay Portal's bundle. It only applies if debug mode is enabled.| 9000 |
|`-p`, `--httpPort` | Sets the HTTP port of Liferay Portal's bundle.| 8080|

To achieve that:

- The tool will ask you to type an image tag from Liferay Portal's nightly builds (check available tags [here](https://hub.docker.com/r/mdelapenya/liferay-portal-nightlies/tags/).
  - If no tag is provided, then it will use current date as tag, i.e. `20180214`.
- It downloads the Docker image to the local engine.
- It checks whether the Docker container this tool spins up is running. In that case, the tool deletes it.
- It spins up a Docker container, using the port configured for Tomcat, and 11311 for OSGi console. The name of the container will be `liferay-portal-nightly`. Once started, please open a web browser in [https://localhost:8080](http://localhost:8080) to check the portal.

### update

Updates lpn (Liferay Portal Nightly) to the latest version on stable channel.

### version

All software has versions. This is lpn (Liferay Portal Nightly).