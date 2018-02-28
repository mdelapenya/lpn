# Changelog

<a name="0.3.0"></a>
### 0.3.0
#### :rocket: Enhancements
* Add :sparkles:`deploy`:sparkles: command, which deploys files to Liferay Portal's deploy directory
* Add `dep` as dependency manager
* Add `stop` command
* Add `codecov` support
* Add functional tests for the CLI :tada:

<a name="0.2.0"></a>
### 0.2.0
#### :rocket: Enhancements
* Support running releases of Liferay Portal

<a name="0.1.2"></a>
### 0.1.2
#### :rocket: Enhancements
* Add a flag for `run` command to set up the `debugPort` to use when debugging.
* Support for :sparkles: **debug mode** :sparkles:
* Add a flag for `checkImage` command to set up the `tag` to be checked.
* Include Windows support.

<a name="0.1.1"></a>
### 0.1.1
#### :rocket: Enhancements
* Add a flag for `run` command to set up the `httpPort` to be used.
* Minor improvement: reorganise code to be more reusable.
* Modularise using `Cobra` as command manager.
    * Add `update` command.
    * Add `log` command.
    * Add `rm` command.
    * Add `pull` command.
    * Add `run` command.
    * Add `version` command.
    * Add `checkImage` command.
    * Add `checkContainer` command.

<a name="0.1.0"></a>
### 0.1.0
* Run a Docker container from Liferay Portal nightly builds.