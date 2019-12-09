# lpn (Liferay Portal Nook)

[![Build Status](https://travis-ci.org/mdelapenya/lpn.svg?branch=portal-snapshots)](https://travis-ci.org/mdelapenya/lpn)

This branch is here just to build the Snapshot for Liferay Portal in a daily manner, so that it's possible to consume the snapshot from the CLI.

## The script
It clones Liferay's official scripts for building the Docker image (see [here](https://github.com/liferay/liferay-docker)), adding the portal snapshot that it's already published on their public servers.
