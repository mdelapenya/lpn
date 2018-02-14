# Liferay Gowl

This Golang CLI makes it easier to run Liferay Portal's nightly builds.

It wraps Docker commands so you just have to run this tool, and pass the specific tag you want to run.

## Requirements

You have to install Docker on your machine first. Check [this guide](https://docs.docker.com/install).

## What does it do?

- The tool asks you to type an image tag from Liferay Portal's nightly builds (check available tags [here](https://hub.docker.com/r/mdelapenya/liferay-portal-nightlies/tags/).
  - If no tag is provided, the it will use current date as tag, i.e. `20180214`.
- It downloads the Docker image to the local engine.
- It checks whether the Docker container this tool spins up is running. In that case, the tool deletes it.
- It spins up a Docker container, using port 8080 for Tomcat, and 11311 for OSGi console. The name of the container will be `liferay-portal-nightly`. Once started, please open a web browser in [https://localhost:8080](http://localhost:8080) to check the portal.