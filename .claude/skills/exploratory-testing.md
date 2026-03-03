# Exploratory Testing Skill for lpn (Liferay Portal Nook)

You are an exploratory tester for **lpn** (Liferay Portal Nook), a Go CLI that simplifies running and managing Liferay Portal Docker containers. Your goal is to demonstrate that the software is not broken by exercising all commands and observing behavior, rather than mechanically asserting expected values.

You approach testing like a curious, experienced QA engineer: you read the docs, set up the environment, run through scenarios, adapt when things behave unexpectedly, and document discrepancies without panicking.

---

## Documentation Sources

Read these sources before testing to understand what the software is supposed to do:

1. **README.md** (project root) — Primary user-facing documentation. Covers installation, configuration, all commands with examples, flags, and usage patterns.
2. **docker/README_INTEGRATION_TESTS.md** — Explains the label-based container identification system, testcontainers-go integration, and database orchestration patterns.
3. **CLI help output** — Run `lpn --help` and `lpn <command> --help` for each command to see the built-in documentation. Compare this against README.md for discrepancies.
4. **Source code in `cmd/`** — Each file corresponds to a cobra command. Read these to understand flag definitions, default values, and error handling paths.
5. **Configuration** — The tool uses `~/.lpn/config.yaml` (auto-created on first run). Default image tags from config: CE=`7.0.6-ga7`, DXP=`7.0.10.8`, Nightly=`master` (but `run`/`pull` commands override to current date `YYYYMMDD` when no tag is specified), Release=`latest`, Commerce=`1.1.1`. Database images: MySQL=`docker.io/mdelapenya/mysql-utf8:5.7`, PostgreSQL=`postgres:9.6-alpine`.
6. **Environment variables** — `LPN_LOG_LEVEL` controls the log level (`TRACE`, `DEBUG`, `INFO`, `WARNING`, `ERROR`, `FATAL`, `PANIC`; default `INFO`). `LPN_LOG_INCLUDE_TIMESTAMP` adds timestamps to log output when set to `TRUE`. These are read at startup in `internal/config.go`.
7. **Feature files in `features/`** — BDD-style scenarios that document the expected behavior of each command. Review these for concrete examples of expected inputs and outputs.
8. **E2E tests in `e2e/`** — Go-based end-to-end tests using `termtest` that exercise the built binary. These define the exact output strings the tool produces.

---

## Environment Setup

Run these steps before testing. If something fails, troubleshoot before moving on.

### 1. Verify Prerequisites

```bash
# Check Docker is available and running
docker info > /dev/null 2>&1 && echo "Docker is running" || echo "ERROR: Docker is not running"

# Check Go is available
go version

# Check current directory is the project root
ls main.go cmd/ e2e/ docker/ liferay/
```

### 2. Install Dependencies

```bash
go get -v -t -d ./...
```

### 3. Build the Binary

```bash
go build -v -o ./bin/lpn .
```

### 4. Verify the Binary

```bash
./bin/lpn version
```

If the binary does not exist at `./bin/lpn`, the build failed. Check the Go build output for errors.

### 5. Set Up PATH (optional)

```bash
export PATH="$(pwd)/bin:$PATH"
```

After this, you can use `lpn` directly instead of `./bin/lpn`. All test scenarios below assume `lpn` is available in PATH or use `./bin/lpn`.

---

## Test Plan

Work through these scenarios in order. For each scenario, run the command, observe the output, and note whether the behavior matches expectations. If something is unexpected, note the discrepancy and continue — do not stop testing.

Use the `release` image type for scenarios requiring a running container, as `mdelapenya/liferay-portal:latest` is the most likely to be available on Docker Hub.

### Image Type Reference

The following image types and their Docker Hub repositories are used throughout:

| Type | Repository | Default Tag |
|------|-----------|------------|
| ce | `liferay/portal` | `7.0.6-ga7` (from config) |
| commerce | `liferay/commerce` | `1.1.1` (from config) |
| dxp | `liferay/dxp` | `7.0.10.8` (from config) |
| nightly | `mdelapenya/portal-snapshot` | Current date `YYYYMMDD` |
| release | `mdelapenya/liferay-portal` | `latest` |

### Container Name Reference

Each image type maps to a container name: `lpn-ce`, `lpn-commerce`, `lpn-dxp`, `lpn-nightly`, `lpn-release`.

Database container names: `db-ce`, `db-commerce`, `db-dxp`, `db-nightly`, `db-release`.

---

### 1. Version and Build Info

**Goal:** Confirm the binary runs and reports version information.
**Covers:** `e2e/version_test.go:TestVersion`, `features/version.feature`

```bash
lpn version
```

**Observe:** Output should contain these strings in order:
- `lpn (Liferay Portal Nook) v`
- `dockerClient=`
- `dockerServer=`
- `golang=`

If Docker is not running, the Docker version fields may be empty — note this but continue. Exit code should be 0.

### 2. Help Output and Command Discovery

**Goal:** Verify all commands are registered and help text is coherent.

```bash
# Root help
lpn --help

# Verify each top-level command exists in help output
lpn help
```

**Observe:** The help output should list these commands: `checkc`, `checki`, `completion`, `deploy`, `license`, `log`, `open`, `prune`, `pull`, `rm`, `rmi`, `run`, `start`, `stop`, `tags`, `update`, `version`. If any command is missing, note it.

#### 2a. Subcommand Help for Commands with Subcommands

For each command that has subcommands, verify subcommand help:

```bash
lpn checkc --help
lpn checki --help
lpn deploy --help
lpn log --help
lpn open --help
lpn pull --help
lpn rm --help
lpn rmi --help
lpn run --help
lpn start --help
lpn stop --help
lpn tags --help
```

**Observe:** Each should show `Available Commands:` with subcommands `ce`, `commerce`, `dxp`, `nightly`, `release`. Commands without subcommands (`completion`, `license`, `prune`, `update`, `version`) should show their own help without listing subcommands.

#### 2b. Pull Subcommand Help

**Covers:** `e2e/pull_test.go:TestPullCommandHelp`, `e2e/pull_test.go:TestPullCommandSubcommands`

```bash
# Pull parent help — should show description and Available Commands
lpn pull --help
```

**Observe:** Should contain `Pulls a Liferay Portal Docker image` and `Available Commands:`.

```bash
# Pull subcommand help — each should show "Pulls a Liferay" in its description
lpn pull ce --help
lpn pull commerce --help
lpn pull dxp --help
lpn pull nightly --help
lpn pull release --help
```

**Observe:** Each pull subcommand help should contain `Pulls a Liferay` in its description. Each should list flags `-f`/`--forceRemoval`, `-t`/`--tag`, and `-V`/`--verbose`.

#### 2c. Deploy Subcommand Help

**Covers:** `e2e/deploy_test.go:TestDeployCommandHelp`, `e2e/deploy_test.go:TestDeployCommandSubcommands`

```bash
lpn deploy --help
```

**Observe:** Should contain `Deploys files or a directory to Liferay Portal's deploy folder` and `Available Commands:`.

```bash
lpn deploy ce --help
lpn deploy commerce --help
lpn deploy dxp --help
lpn deploy nightly --help
lpn deploy release --help
```

**Observe:** Each subcommand help should contain `Deploys files or a directory` and list the `-f`/`--files` and `-d`/`--dir` flags.

#### 2d. Run Subcommand Help

```bash
lpn run --help
```

**Observe:** Should contain `Runs a Liferay Portal instance` and `Available Commands:`.

```bash
lpn run ce --help
lpn run commerce --help
lpn run dxp --help
lpn run nightly --help
lpn run release --help
```

**Observe:** Each should list flags: `-p`/`--httpPort`, `-d`/`--debug`, `-D`/`--debugPort`, `-g`/`--gogoPort`, `-s`/`--datastore`, `-t`/`--tag`, `-m`/`--memory`, and `-V`/`--verbose`.

#### 2e. Tags Subcommand Help

```bash
lpn tags --help
```

**Observe:** Should contain `Lists all tags for Liferay Portal Docker image` and `Available Commands:`.

```bash
lpn tags ce --help
lpn tags commerce --help
lpn tags dxp --help
lpn tags nightly --help
lpn tags release --help
```

**Observe:** Each should list flags: `-p`/`--page`, `-s`/`--size`, and `-V`/`--verbose`.

#### 2f. Other Subcommand Help

Verify help for remaining commands with subcommands:

```bash
lpn checkc ce --help
lpn checki ce --help
lpn log ce --help
lpn open ce --help
lpn rm ce --help
lpn rmi ce --help
lpn start ce --help
lpn stop ce --help
```

**Observe:** Each should display a meaningful description. `checki` and `rmi` subcommands should show a `-t`/`--tag` flag. All should show `-V`/`--verbose`.

#### 2g. Standalone Command Help

```bash
lpn completion --help
lpn license --help
lpn prune --help
lpn update --help
lpn version --help
```

**Observe:** Each should show its own description without `Available Commands:`.

### 3. License Display

**Goal:** Verify the license command outputs license text.
**Covers:** `e2e/license_test.go:TestLicense`, `features/license.feature`

```bash
lpn license
```

**Observe:** Should print `GNU Lesser General Public License` text. If it shows an error about missing assets, the binary was not built with embedded assets. Exit code should be 0.

### 4. Shell Completion

**Goal:** Verify completion script generation.
**Covers:** `e2e/completion_test.go:TestCompletion`

```bash
lpn completion
```

**Observe:** Should output a ZSH completion script starting with `#compdef`. The output should be valid shell script syntax. Exit code should be 0.

### 5. Tag Listing (Docker Hub Query)

**Goal:** Verify the tool can query Docker Hub for available image tags.
**Covers:** `e2e/tags_test.go:TestTags`, `features/tags.feature`

```bash
# List CE tags (default page size 25)
lpn tags ce

# List for other image types
lpn tags dxp
lpn tags release
lpn tags nightly
lpn tags commerce
```

**Observe:** Each should display a table with `Image:Tag` and `Size` columns, plus a footer showing pagination (e.g., `1 of N`). Exit code should be 0. If Docker Hub is unreachable or rate-limited, note the error but continue. Commerce tags may return empty results if the repository no longer exists.

#### 5b. Tags with Custom Pagination

**Covers:** `features/tags.feature` (pagination scenarios)

```bash
# List with custom page and size
lpn tags ce --page 1 --size 5

# List with only custom size
lpn tags ce --size 10
```

**Observe:** With `--size 5`, should display 5 tags. With `--page 1 --size 5`, should display page 1 with 5 results.

#### 5c. Tags with Invalid Pagination

**Covers:** `features/tags.feature` (not found resource error)

```bash
# Use a very high page number that doesn't exist
lpn tags ce --page 99999 --size 1
```

**Observe:** Should print `There are no available tags for that pagination. Please use --page and --size arguments to filter properly`. Exit code should be 0 (this is an informational message, not a crash).

### 6. Check Container Status (No Container Running)

**Goal:** Verify graceful handling when no containers exist.
**Covers:** `e2e/checkc_test.go:TestCheckContainerNotRunning`, `features/checkc.feature`

```bash
lpn checkc ce
lpn checkc commerce
lpn checkc dxp
lpn checkc nightly
lpn checkc release
```

**Observe:** Each should report `Container does NOT exist in the system` with `container=lpn-<type>` (e.g., `container=lpn-ce`, `container=lpn-release`). Exit code should be 0 (the absence of a container is not an error).

### 7. Check Image Status (No Image Pulled)

**Goal:** Verify graceful handling when images are not locally available.
**Covers:** `e2e/checki_test.go:TestCheckImageWhenNotPresent`, `features/checki.feature`

```bash
lpn checki ce
lpn checki commerce
lpn checki dxp
lpn checki nightly
lpn checki release
```

**Observe:** Each should report `Image has NOT been pulled from Docker Hub` with the fully qualified image name. Exit code should be 0.

#### 7b. Check Image with Specific Tag

**Covers:** `features/checki.feature` (tag-specific scenarios)

```bash
lpn checki ce -t 7.0.6-ga7
lpn checki dxp -t 7.0.10.8
lpn checki release -t latest
lpn checki nightly -t master
lpn checki commerce -t 1.1.1
```

**Observe:** Each should report the image status for the specific tag. If not pulled, shows `Image has NOT been pulled from Docker Hub`. Exit code should be 0.

### 8. Pull an Image

**Goal:** Verify the tool can pull Docker images from Docker Hub.
**Covers:** `features/pull.feature`

```bash
# Pull a small/common image to test the pull mechanism
lpn pull release -t latest
```

**Observe:** The image should be pulled successfully. If the pull takes a long time, wait patiently — large images can take minutes. If the image is not found, note the error and try another image type.

After pulling, verify:

```bash
lpn checki release -t latest
```

**Observe:** Should now report `Image has been pulled from Docker Hub`.

### 8b. Pull with Force Removal (Image Cached)

**Goal:** Verify the `--forceRemoval/-f` flag removes the cached local image before pulling.
**Covers:** `features/pull.feature` (force removal scenarios)

```bash
# First ensure an image is present (use the one pulled in scenario 8)
lpn checki release -t latest

# Pull again with force removal — should remove the local cache and re-pull
lpn pull release -t latest -f

# Verify the image is still present after re-pull
lpn checki release -t latest
```

**Observe:** With `-f`/`--forceRemoval`, the tool should remove the local image first and then re-pull it. The image should be present after the operation.

### 8c. Pull with Force Removal (Image Not Cached)

**Goal:** Verify force removal handles missing images gracefully.
**Covers:** `features/pull.feature` (force removal of non-present image)

```bash
# Try to force-pull an image that isn't cached locally
lpn pull ce -t 7.0.6-ga7 -f
```

**Observe:** Should log `The image was not found in the local cache. Skipping removal` and proceed to pull. This is an info message, not an error.

### 9. Pull with Non-Existent Tag (Error Handling)

**Goal:** Verify error handling for invalid image tags.
**Covers:** `e2e/pull_test.go:TestPullNonExistentImage`, `features/pull.feature`

```bash
lpn pull ce -t "nonexistent-tag-12345"
```

**Observe:** Should display `The image could not be pulled` with the image reference `liferay/portal`. Exit code should be 1.

Repeat for all image types to verify consistent error handling with the correct repository names:

```bash
lpn pull commerce -t "nonexistent-tag-12345"
```

**Observe:** Error should reference `liferay/commerce`.

```bash
lpn pull dxp -t "nonexistent-tag-12345"
```

**Observe:** Error should reference `liferay/dxp`.

```bash
lpn pull nightly -t "nonexistent-tag-12345"
```

**Observe:** Error should reference `mdelapenya/portal-snapshot`.

```bash
lpn pull release -t "nonexistent-tag-12345"
```

**Observe:** Error should reference `mdelapenya/liferay-portal`.

### 10. Container Lifecycle: Run, Check, Stop, Start, Remove

**Goal:** Exercise the full container lifecycle. This is the core workflow of the tool.
**Covers:** `features/run.feature`, `features/stop.feature`, `features/start.feature`, `features/rm.feature`, `features/checkc.feature`

**Note:** This requires a previously pulled image. Use the image pulled in scenario 8, or pull one now.

```bash
# Run a container (use release type with latest tag)
lpn run release -t latest -p 8080
```

**Observe:** Should show `The container has been run successfully` with `container=lpn-release`, `image=mdelapenya/liferay-portal:latest`, `datastore=hsql`, `httpPort=8080`, `gogoPort=11311`, `debug=false`, `debugPort=9000`, `memory=`.

```bash
# Check if the container is running
lpn checkc release
```

**Observe:** Should report `Container DOES exist in the system` with `container=lpn-release`.

```bash
# Check logs (briefly — this tails logs, so interrupt after a few seconds with Ctrl+C or timeout)
timeout 10 lpn log release || true
```

**Observe:** Should display container log output. If the container is still starting, logs may show startup messages.

```bash
# Stop the container
lpn stop release
```

**Observe:** Should succeed without errors. Container transitions to stopped state.

```bash
# Verify it's stopped — checkc may still show it exists but stopped
lpn checkc release
```

```bash
# Start the stopped container
lpn start release
```

**Observe:** Should succeed, restarting the previously stopped container.

```bash
# Verify it's running again
lpn checkc release
```

```bash
# Remove the container
lpn rm release
```

**Observe:** Should succeed without errors.

```bash
# Verify it's gone
lpn checkc release
```

**Observe:** Should report `Container does NOT exist in the system`.

If the container takes time to start (Liferay Portal is large), wait up to 2-3 minutes. If it fails to start, check Docker logs directly: `docker logs lpn-release`.

### 10b. Start an Already-Running Container

**Goal:** Verify that starting an already-running container is handled gracefully.
**Covers:** `features/start.feature` (started container scenario)

```bash
# Run a container first
lpn run release -t latest -p 8080

# Try starting it when it's already running
lpn start release

# Clean up
lpn rm release
```

**Observe:** Starting an already-running container should not crash or produce unexpected errors. It may show the container is already running.

### 11. Deploy Command Validation

**Goal:** Verify deploy command input validation, error handling, and help output.
**Covers:** `e2e/deploy_test.go:TestDeployCommandWithoutFile`

```bash
# Deploy without specifying files or dir (should error)
lpn deploy ce
```

**Observe:** Should print `Please pass a valid path to a file or to a directory as argument` and exit with code 1.

### 11a. Deploy with Non-Existent File

**Covers:** `e2e/deploy_test.go:TestDeployCommandWithNonExistentFile`

```bash
lpn deploy ce -f /tmp/nonexistent-file-12345.jar
```

**Observe:** Should fail with exit code 1.

### 11b. Deploy Valid File Without Running Container

**Goal:** Verify deploy fails gracefully when deploying a valid file to a non-running container for all flavors.
**Covers:** `e2e/deploy_test.go:TestDeployCommandWithoutContainer`, `features/deploy.feature`

```bash
# Create a temporary test file
echo "test jar content" > /tmp/test-deploy.jar

# Deploy the valid file to each flavor — all should fail since no container is running
lpn deploy ce -f /tmp/test-deploy.jar
lpn deploy commerce -f /tmp/test-deploy.jar
lpn deploy dxp -f /tmp/test-deploy.jar
lpn deploy nightly -f /tmp/test-deploy.jar
lpn deploy release -f /tmp/test-deploy.jar

# Clean up
rm /tmp/test-deploy.jar
```

**Observe:** Each should fail with exit code 1 because no container is running. The error should be about the container not being found, not about the file being invalid. This confirms the tool validates the file exists before attempting to deploy.

### 11c. Deploy with Directory Flag

**Goal:** Verify the `--dir/-d` flag for directory deployment.
**Covers:** `features/deploy.feature` (directory deployment scenarios)

```bash
# Create a temporary directory with test files
mkdir -p /tmp/test-deploy-dir
echo "module A" > /tmp/test-deploy-dir/moduleA.jar
echo "module B" > /tmp/test-deploy-dir/moduleB.jar

# Deploy a directory without a running container (should fail)
lpn deploy release -d /tmp/test-deploy-dir

# Deploy with invalid directory (should error)
lpn deploy release -d /tmp/nonexistent-dir-12345

# Clean up
rm -rf /tmp/test-deploy-dir
```

**Observe:** Deploying a directory to a non-running container should fail with exit code 1. Deploying with an invalid directory should print `The directory is not valid` and exit with code 1.

### 11d. Deploy Directory Skipping Subdirectories

**Goal:** Verify that directory deployment only deploys first-level files and skips subdirectories.
**Covers:** `features/deploy.feature` (subdirectory skipping)

```bash
# This scenario requires a running container.
# First, run a container:
lpn run release -t latest -p 8080

# Create a directory structure with subdirectories
mkdir -p /tmp/test-deploy-nested/subdir
echo "top level file" > /tmp/test-deploy-nested/moduleA.jar
echo "nested file" > /tmp/test-deploy-nested/subdir/nested.jar

# Deploy the directory
lpn deploy release -d /tmp/test-deploy-nested

# Clean up
lpn rm release
rm -rf /tmp/test-deploy-nested
```

**Observe:** Only `moduleA.jar` should be deployed. The `subdir/nested.jar` should be skipped (the tool only deploys first-level files, no recursive deployment). Look for `File deployed successfully to deploy dir` messages — should only reference `moduleA.jar`.

### 11e. Deploy Multiple Files

**Goal:** Verify deploying multiple comma-separated files with the `-f` flag.
**Covers:** `features/deploy.feature` (multiple file deployment)

```bash
# This scenario requires a running container.
lpn run release -t latest -p 8080

# Create multiple test files
echo "file A" > /tmp/test-deploy-a.jar
echo "file B" > /tmp/test-deploy-b.jar

# Deploy multiple files (comma-separated)
lpn deploy release -f /tmp/test-deploy-a.jar,/tmp/test-deploy-b.jar

# Clean up
lpn rm release
rm /tmp/test-deploy-a.jar /tmp/test-deploy-b.jar
```

**Observe:** Both files should be deployed. Should see `File deployed successfully to deploy dir` for each file.

### 12. Deploy to a Running Container

**Goal:** Verify file deployment works when a container is running.
**Covers:** `features/deploy.feature` (single file deployment)

```bash
# First, ensure a container is running
lpn run release -t latest -p 8080

# Create a test file
echo "test content" > /tmp/test-deploy.txt

# Deploy the file
lpn deploy release -f /tmp/test-deploy.txt

# Clean up
lpn rm release
rm /tmp/test-deploy.txt
```

**Observe:** The deploy command should copy the file to the container's deploy folder. Should show `File deployed successfully to deploy dir`.

### 13. Remove Image

**Goal:** Verify image removal.
**Covers:** `features/rmi.feature`

```bash
# Remove a previously pulled image
lpn rmi release -t latest

# Verify it's gone
lpn checki release -t latest
```

**Observe:** The image should be removed. After removal, `checki` should report `Image has NOT been pulled from Docker Hub`.

#### 13b. Remove Non-Existent Image

**Covers:** `features/rmi.feature` (image does not exist scenario)

```bash
# Try to remove an image that doesn't exist
lpn rmi ce -t nonexistent-tag-12345
```

**Observe:** Should handle the missing image gracefully without crashing.

### 14. Run with Database Orchestration

**Goal:** Verify database container orchestration with the `--datastore` flag.
**Covers:** `features/run.feature`, `features/start.feature`, `features/stop.feature`, `features/rm.feature` (database service scenarios)

```bash
# Run with MySQL datastore
lpn run release -t latest -p 8080 --datastore mysql
```

**Observe:** Should show `The stack has been run successfully` (note: "stack" not "container" when using a database). Two containers should be created: one for the portal (`lpn-release`) and one for the database. The success log should include `datastore=mysql`.

```bash
# Check if the portal container is running
lpn checkc release

# List Docker containers to see both portal and database containers
docker ps --filter "label=lpn-container-name"
```

**Observe:** Both portal and database containers should be visible.

#### 14b. Stop and Start with Database Service

**Covers:** `features/stop.feature`, `features/start.feature` (service scenarios)

```bash
# Stop the container (should also handle the database service)
lpn stop release

# Start it again
lpn start release

# Clean up
lpn rm release
```

**Observe:** Stop and start should handle the database service containers. The `rm` command should clean up both portal and database containers.

#### 14c. Run with PostgreSQL

```bash
# Run with PostgreSQL datastore
lpn run release -t latest -p 8081 --datastore postgresql

# Verify
lpn checkc release

# Clean up
lpn rm release
```

**Observe:** Same behavior as MySQL but with PostgreSQL. If testcontainers-go Ryuk is running, a third container (ryuk) may appear — this is normal.

### 15. Run with Custom Flags

**Goal:** Verify all `run` flags are accepted and applied.
**Covers:** `features/run.feature` (flag scenarios)

```bash
# Run with debug mode, custom ports, and memory setting
lpn run release -t latest -p 9090 -g 22222 --debug -D 8787 -m "-Xmx4096m"
```

**Observe:** The container should start with the specified settings. The success log should show `httpPort=9090`, `gogoPort=22222`, `debug=true`, `debugPort=8787`, and `memory=-Xmx4096m`.

```bash
# Verify the container is running
lpn checkc release

# Clean up
lpn rm release
```

Run flags reference:
- `-p`/`--httpPort` (int, default 8080) — HTTP port
- `-d`/`--debug` (bool, default false) — Enable debug mode
- `-D`/`--debugPort` (int, default 9000) — Debug port (only applies when debug enabled)
- `-g`/`--gogoPort` (int, default 11311) — GoGo Shell port
- `-s`/`--datastore` (string, default "hsql") — Database service: hsql, mysql, or postgresql
- `-t`/`--tag` (string) — Image tag to run
- `-m`/`--memory` (string) — JVM memory configuration (e.g., `-Xmx2048m`)

### 15b. Run with Debug Mode Verification

**Goal:** Verify debug mode sets the correct environment variables.
**Covers:** `features/run.feature` (debug enabled scenario)

```bash
# Run with debug enabled
lpn run release -t latest -p 8080 --debug -D 9000

# Inspect the container to verify environment variables
docker inspect lpn-release --format '{{.Config.Env}}' | grep -E "LIFERAY_JPDA_ENABLED|DEBUG_MODE"

# Clean up
lpn rm release
```

**Observe:** The container should have `LIFERAY_JPDA_ENABLED=true` environment variable set. The success log should show `debug=true`.

### 15c. Run with Memory Configuration

**Goal:** Verify memory configuration is passed to the container.
**Covers:** `features/run.feature` (memory enabled scenario)

```bash
# Run with custom memory
lpn run release -t latest -p 8080 -m "-Xmx2048m"

# Inspect the container to verify environment variable
docker inspect lpn-release --format '{{.Config.Env}}' | grep "LIFERAY_JVM_OPTS"

# Clean up
lpn rm release
```

**Observe:** The container should have `LIFERAY_JVM_OPTS=-Xmx2048m` environment variable set. The success log should show `memory=-Xmx2048m`.

### 15d. Run with Port Conflict (Error Handling)

**Goal:** Verify graceful error handling when a port is already in use.
**Covers:** `features/run.feature` (failure scenario)

```bash
# Start something on a port first (e.g., run a container on port 9999)
docker run -d --name port-blocker -p 9999:80 nginx:alpine 2>/dev/null || true

# Try to run lpn on the same port
lpn run release -t latest -p 9999

# Clean up
docker rm -f port-blocker 2>/dev/null || true
lpn rm release 2>/dev/null || true
```

**Observe:** Should show `Impossible to run the container` when the port is already in use. The error should mention the port conflict.

### 16. Prune Command

**Goal:** Verify the prune command cleans up all LPN resources.

```bash
lpn prune
```

**Observe:** Should remove all LPN containers and images. Should print `LPN state has been pruned!` on success.

### 17. Update Command

**Goal:** Verify the update command behavior.

```bash
lpn update
```

**Observe:** Should print a message about Equinox updates being disabled and provide a URL to download releases: `Updates from Equinox are disabled. Please go to 'https://mdelapenya.github.io/lpn/releases.html' to download your release`. Exit code should be 1 (since updates are disabled).

### 18. Parent Commands Without Subcommands

**Goal:** Verify that running a parent command without a subcommand shows a helpful message.

```bash
lpn run
lpn stop
lpn start
lpn rm
lpn deploy
lpn pull
lpn tags
lpn checkc
lpn checki
lpn log
lpn open
lpn rmi
```

**Observe:** Each should print a warning: `Please run this command adding 'ce', 'commerce', 'dxp', 'nightly' or 'release' subcommands.`

### 19. Verbose Flag

**Goal:** Verify the verbose flag enables debug logging.

```bash
lpn checkc ce -V
lpn version -V 2>&1 | head -20
```

**Observe:** With `-V`, output should include debug-level log messages. Compare with non-verbose output to verify additional information is shown.

### 19b. Log Level Environment Variable

**Goal:** Verify the `LPN_LOG_LEVEL` environment variable controls log output independently of the `-V` flag.

```bash
# Run with DEBUG log level
LPN_LOG_LEVEL=DEBUG lpn checkc ce 2>&1

# Run with ERROR log level (should suppress info and warning messages)
LPN_LOG_LEVEL=ERROR lpn checkc ce 2>&1

# Test timestamp inclusion
LPN_LOG_INCLUDE_TIMESTAMP=TRUE lpn checkc ce 2>&1
```

**Observe:** With `LPN_LOG_LEVEL=DEBUG`, output should include additional debug-level messages. With `LPN_LOG_LEVEL=ERROR`, only error messages should appear (the "Container does NOT exist" warning should be suppressed). With `LPN_LOG_INCLUDE_TIMESTAMP=TRUE`, each log line should include a timestamp. Valid log levels are: `TRACE`, `DEBUG`, `INFO`, `WARNING`, `ERROR`, `FATAL`, `PANIC`.

### 20. Stop/Start/Remove Without Running Container

**Goal:** Verify graceful handling of operations on non-existent containers.
**Covers:** `features/stop.feature`, `features/start.feature`, `features/rm.feature` (container does not exist scenarios)

```bash
lpn stop ce
lpn start ce
lpn rm ce
```

**Observe:** Each should handle the missing container gracefully. `stop` should show `Impossible to stop the container`, `start` should show `Impossible to start the container`, and `rm` should show `Impossible to remove the container`. These are warnings, not crashes. Exit code should be 0.

### 21. Open Command (Headless Environment)

**Goal:** Test the open command behavior in a headless environment.

```bash
lpn open ce
```

**Observe:** In a headless environment (CI/server), this will fail because there's no browser or display. It should fail with `Failed to open browser` or an error about `xdg-open` not being available rather than hanging. On a desktop, it would open `http://localhost:<port>` in the default browser. On Linux it uses `xdg-open`, on macOS it uses `open`, on Windows it uses `rundll32`.

### 22. Run All Image Types

**Goal:** Verify that all image types can be run (not just release).
**Covers:** `features/run.feature` (all image types)

For each image type that has a pullable image, verify the run command works:

```bash
# CE with its default tag
lpn run ce -t 7.0.6-ga7 -p 8080
lpn checkc ce
lpn rm ce

# Commerce (may fail if image no longer exists on Docker Hub)
lpn run commerce -t 1.1.1 -p 8080
lpn checkc commerce
lpn rm commerce

# DXP (may require authentication for Docker Hub)
lpn run dxp -t 7.0.10.8 -p 8080
lpn checkc dxp
lpn rm dxp
```

**Observe:** Each should follow the same lifecycle pattern. Some may fail if the images are no longer available on Docker Hub — note these discrepancies but continue.

---

## Troubleshooting Guide

### Docker Daemon Not Running

**Symptoms:** Commands fail with errors about connecting to Docker socket, `Cannot connect to the Docker daemon`, or `docker: command not found`.

**Resolution:**
1. Check if Docker is installed: `which docker`
2. Check if Docker daemon is running: `docker info`
3. Start Docker:
   - Linux: `sudo systemctl start docker`
   - macOS: Open Docker Desktop
   - Windows: Start Docker Desktop
4. If using rootless Docker, ensure the socket path is correct: `export DOCKER_HOST=unix:///run/user/$(id -u)/docker.sock`

### Image Not Found on Docker Hub

**Symptoms:** `The image could not be pulled` error when using `lpn pull` or `lpn run`.

**Resolution:**
1. Verify the image exists: `docker pull <image>:<tag>` directly
2. Check if Docker Hub is rate-limiting you: `docker login` to authenticate
3. The repository may have been renamed or removed:
   - CE: `liferay/portal` (official)
   - DXP: `liferay/dxp` (official, may require authentication)
   - Nightly: `mdelapenya/portal-snapshot` (unofficial, may not exist anymore)
   - Release: `mdelapenya/liferay-portal` (unofficial)
   - Commerce: `liferay/commerce` (may be discontinued)
4. Try a different tag: `lpn tags <type>` to see available tags
5. Check configuration: `cat ~/.lpn/config.yaml` for custom repository URLs

### Port Conflicts

**Symptoms:** Container fails to start with error about port already in use: `Bind for 0.0.0.0:8080 failed: port is already allocated`.

**Resolution:**
1. Find what's using the port: `lsof -i :8080` or `ss -tlnp | grep 8080`
2. Use a different port: `lpn run release -p 9090`
3. Stop the conflicting service or container: `docker stop <container_id>`
4. Kill the process using the port: `kill <pid>` (use with caution)

### Container Name Conflicts

**Symptoms:** Error about container name already in use when running `lpn run`.

**Resolution:**
1. Remove the existing container: `lpn rm <type>` (e.g., `lpn rm release`)
2. If that fails, remove directly: `docker rm -f lpn-release`
3. Use `lpn prune` to clean up all LPN resources
4. Note: LPN uses label-based identification (`lpn-container-name` label), so name conflicts are rare with the current architecture

### Binary Not Found

**Symptoms:** `lpn: command not found` or `./bin/lpn: No such file or directory`.

**Resolution:**
1. Build the binary: `go build -v -o ./bin/lpn .`
2. Check it was built: `ls -la ./bin/lpn`
3. Add to PATH: `export PATH="$(pwd)/bin:$PATH"`
4. If Go build fails, check Go version: `go version` (requires Go 1.25+)
5. Install dependencies: `go get -v -t -d ./...`

### Configuration File Issues

**Symptoms:** Unexpected default values, wrong image tags, or config-related errors.

**Resolution:**
1. Check config location: `cat ~/.lpn/config.yaml`
2. Delete to regenerate defaults: `rm ~/.lpn/config.yaml` then run any `lpn` command
3. Verify YAML syntax if manually edited
4. Default config tags: CE=`7.0.6-ga7`, DXP=`7.0.10.8`, Release=`latest`, Nightly=`master` (but `run`/`pull` override to current date), Commerce=`1.1.1`

### Testcontainers Ryuk Cleanup Issues

**Symptoms:** Database containers are being removed unexpectedly, or Ryuk container appears and removes test containers.

**Resolution:**
1. For production/persistent containers, disable Ryuk: `export TESTCONTAINERS_RYUK_DISABLED=true`
2. Or create `~/.testcontainers.properties` with `ryuk.disabled=true`
3. For testing, Ryuk cleanup is expected and normal behavior

### Database Container Fails to Start

**Symptoms:** `lpn run <type> --datastore mysql` fails, or database container exits immediately.

**Resolution:**
1. Check Docker has enough memory allocated (MySQL needs ~512MB, PostgreSQL ~256MB)
2. Check disk space: `df -h`
3. Check Docker logs: `docker logs <db-container-name>`
4. Try with PostgreSQL instead of MySQL (lighter): `lpn run release --datastore postgresql`
5. Verify the database image can be pulled: `docker pull docker.io/mdelapenya/mysql-utf8:5.7` or `docker pull postgres:9.6-alpine` (these are the default database images configured in `internal/config.go`)

### Slow Container Startup

**Symptoms:** `lpn run` seems to hang or takes very long.

**Resolution:**
1. Liferay Portal containers are large (1-2GB) and take 1-5 minutes to start. This is normal.
2. Check progress with logs: `lpn log <type>` or `docker logs -f lpn-<type>`
3. Look for `Server startup in` in the logs to confirm it's ready
4. If it's the first run, the image needs to be pulled first, which adds time
5. Allocate more memory to Docker if the container is OOM-killed

### E2E Test Failures

**Symptoms:** E2E tests fail when running `go test -v ./e2e/... -timeout 30m`.

**Resolution:**
1. Ensure the binary is built: `go build -v -o ./bin/lpn .`
2. Ensure Docker is running: `docker info`
3. Check timeout — some tests need up to 30 seconds per case
4. Run tests individually to isolate failures: `go test -v ./e2e/ -run TestVersion`
5. Check if previous test containers are still running: `docker ps | grep lpn`
6. Clean up: `lpn prune` or `docker rm -f $(docker ps -aq --filter "label=lpn-container-name")`

### Deploy Fails with Running Container

**Symptoms:** `lpn deploy <type> -f <file>` fails even though a container is running.

**Resolution:**
1. Verify the container is running: `lpn checkc <type>`
2. Verify the file exists: `ls -la <file-path>`
3. Check the container has a deploy folder: `docker exec lpn-<type> ls /opt/liferay/deploy/` (or the appropriate deploy path for the image type)
4. Ensure the file path doesn't contain special characters
5. If deploying multiple files, ensure the comma-separated list has no spaces: `-f file1.jar,file2.jar`
