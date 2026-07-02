# Contributing to lpn (Liferay Portal Nook)

First off, thanks for taking the time to contribute! 🎉

`lpn` is a CLI, written in Go, that makes running Liferay Portal's Docker
images easy. This guide explains how to get a development environment up and
running, and how to get your changes merged.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Where do Pull Requests go?](#where-do-pull-requests-go)
- [Reporting Issues](#reporting-issues)
- [Development Setup](#development-setup)
- [Building](#building)
- [Running the Tests](#running-the-tests)
- [Code Style](#code-style)
- [Submitting a Pull Request](#submitting-a-pull-request)

## Code of Conduct

Be kind and respectful. We want `lpn` to be a welcoming project for everybody,
so please assume good intentions and keep the conversation constructive.

## Where do Pull Requests go?

**Please open your Pull Requests against the `develop` branch, not `master` /
`main`.**

`develop` is the integration branch where all day-to-day work lands. Stable
releases are cut from `develop` into the release branch, so targeting `develop`
keeps things flowing and avoids extra rebasing for everybody.

When you open a PR, double-check the base branch selector at the top of the
GitHub "Open a pull request" screen and make sure it reads `develop`. If you
already opened one against the wrong branch, you can simply edit the PR and
change the base branch — no need to close and reopen it.

## Reporting Issues

Found a bug or have an idea? Please [open an issue](https://github.com/mdelapenya/lpn/issues).

When reporting a bug, it helps a lot if you include:

- The version of `lpn` you are using (`lpn version`).
- Your operating system and Docker version.
- The exact command you ran and the full output (setting
  `LPN_LOG_LEVEL=DEBUG` gives more detail).
- What you expected to happen versus what actually happened.

## Development Setup

You will need:

- [Go](https://golang.org/doc/install) (see the `go` directive in
  [`go.mod`](./go.mod) for the required version).
- [Docker](https://docs.docker.com/install), which the CLI wraps and which the
  end-to-end tests rely on.

Then clone the repository and fetch the dependencies:

```shell
git clone https://github.com/mdelapenya/lpn
cd lpn
go mod download
```

If you plan to send a Pull Request, [fork the repository](https://github.com/mdelapenya/lpn/fork)
first and clone your fork instead.

## Building

Build the binary locally with:

```shell
go build -o ./bin/lpn .
```

You can then run the freshly built binary:

```shell
./bin/lpn help
```

## Running the Tests

Run the unit tests (everything except the end-to-end suite):

```shell
go test $(go list ./... | grep -v /e2e)
```

Run the end-to-end tests (these require Docker and can take a while):

```shell
go test ./e2e/... -timeout 30m
```

To generate a coverage report, use the helper script:

```shell
./scripts/coverage.sh
```

Please make sure the test suite passes before opening a Pull Request. The same
tests run in CI on every Pull Request against `develop`.

## Code Style

- Format your code with `gofmt` (or `go fmt ./...`) before committing.
- Run `go vet ./...` to catch common mistakes.
- Follow standard, idiomatic Go conventions.
- Keep functions focused and use meaningful names.
- Add or update tests (`*_test.go`) alongside your changes.

## Submitting a Pull Request

1. Fork the repository and create a topic branch from `develop`:

   ```shell
   git checkout develop
   git pull
   git checkout -b my-awesome-change
   ```

2. Make your changes, adding tests where it makes sense.
3. Format and lint your code (`go fmt ./...` and `go vet ./...`).
4. Make sure the tests pass (`go test ./...`).
5. Write clear commit messages that explain the "why" of your change.
6. Push your branch and open a Pull Request **against the `develop` branch**.
7. Fill in the Pull Request template and link any related issue.

A maintainer will review your Pull Request as soon as possible. Thanks again
for contributing to `lpn`! 🚀
