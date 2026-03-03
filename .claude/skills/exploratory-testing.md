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

### 1. Version and Build Info

**Goal:** Confirm the binary runs and reports version information.

```bash
lpn version
```

**Observe:** Output should contain `lpn (Liferay Portal Nook) v`, `dockerClient=`, `dockerServer=`, and `golang=`. If Docker is not running, the Docker version fields may be empty — note this but continue.

### 2. Help Output and Command Discovery

**Goal:** Verify all commands are registered and help text is coherent.

```bash
# Root help
lpn --help

# Verify each top-level command exists in help output
lpn help
```

**Observe:** The help output should list these commands: `checkc`, `checki`, `completion`, `deploy`, `license`, `log`, `open`, `prune`, `pull`, `rm`, `rmi`, `run`, `start`, `stop`, `tags`, `update`, `version`. If any command is missing, note it.

Then, for each command that has subcommands, verify subcommand help:

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

Then verify individual subcommand help for key commands:

```bash
# Pull subcommand help — each should show "Pulls a Liferay" in its description
lpn pull ce --help
lpn pull commerce --help
lpn pull dxp --help
lpn pull nightly --help
lpn pull release --help
```

**Observe:** `lpn pull --help` should show `Pulls a Liferay Portal Docker image` and `Available Commands:`. Each pull subcommand help should contain `Pulls a Liferay` in its description.

### 3. License Display

**Goal:** Verify the license command outputs license text.

```bash
lpn license
```

**Observe:** Should print `GNU Lesser General Public License` text. If it shows an error about missing assets, the binary was not built with embedded assets.

### 4. Shell Completion

**Goal:** Verify completion script generation.

```bash
lpn completion
```

**Observe:** Should output a ZSH completion script starting with `#compdef`. The output should be valid shell script syntax.

### 5. Tag Listing (Docker Hub Query)

**Goal:** Verify the tool can query Docker Hub for available image tags.

```bash
# List CE tags (default page size 25)
lpn tags ce

# List with custom pagination
lpn tags ce --page 1 --size 5

# List for other image types
lpn tags dxp
lpn tags release
lpn tags nightly
lpn tags commerce
```

**Observe:** Each should display a table with `Image:Tag` and `Size` columns, plus a footer showing pagination (e.g., `1 of N`). If Docker Hub is unreachable or rate-limited, note the error but continue. Commerce tags may return empty results if the repository no longer exists.

### 6. Check Container Status (No Container Running)

**Goal:** Verify graceful handling when no containers exist.

```bash
lpn checkc ce
lpn checkc commerce
lpn checkc dxp
lpn checkc nightly
lpn checkc release
```

**Observe:** Each should report `Container does NOT exist in the system` with `container=lpn-<type>`. Exit code should be 0 (the absence of a container is not an error).

### 7. Check Image Status (No Image Pulled)

**Goal:** Verify graceful handling when images are not locally available.

```bash
lpn checki ce
lpn checki commerce
lpn checki dxp
lpn checki nightly
lpn checki release
```

**Observe:** Each should report `Image has NOT been pulled from Docker Hub` with the fully qualified image name. Exit code should be 0.

### 8. Pull an Image

**Goal:** Verify the tool can pull Docker images from Docker Hub.

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

### 8b. Pull with Force Removal

**Goal:** Verify the `--forceRemoval/-f` flag removes the cached local image before pulling.

```bash
# First ensure an image is present (use the one pulled in scenario 8)
lpn checki release -t latest

# Pull again with force removal — should remove the local cache and re-pull
lpn pull release -t latest -f

# Verify the image is still present after re-pull
lpn checki release -t latest
```

**Observe:** With `-f`/`--forceRemoval`, the tool should attempt to remove the local image first. If the image was not cached, it should log `The image was not found in the local cache. Skipping removal` and proceed to pull. If the image was cached, it removes it then pulls again.

### 9. Pull with Non-Existent Tag (Error Handling)

**Goal:** Verify error handling for invalid image tags.

```bash
lpn pull ce -t "nonexistent-tag-12345"
```

**Observe:** Should display `The image could not be pulled` with the image reference (`liferay/portal`). Exit code should be 1.

Repeat for other image types to verify consistent error handling:

```bash
lpn pull dxp -t "nonexistent-tag-12345"
lpn pull nightly -t "nonexistent-tag-12345"
lpn pull release -t "nonexistent-tag-12345"
lpn pull commerce -t "nonexistent-tag-12345"
```

### 10. Container Lifecycle: Run, Check, Stop, Start, Remove

**Goal:** Exercise the full container lifecycle. This is the core workflow of the tool.

**Note:** This requires a previously pulled image. Use the image pulled in scenario 8, or pull one now.

```bash
# Run a container (use release type with latest tag)
lpn run release -t latest -p 8080

# Check if the container is running
lpn checkc release

# Check logs (briefly — this tails logs, so interrupt after a few seconds with Ctrl+C or timeout)
timeout 10 lpn log release || true

# Stop the container
lpn stop release

# Verify it's stopped — checkc may still show it exists but stopped
lpn checkc release

# Start the stopped container
lpn start release

# Verify it's running again
lpn checkc release

# Remove the container
lpn rm release

# Verify it's gone
lpn checkc release
```

**Observe:** Each step should produce appropriate log messages. The `run` command should show the container being created. The `stop`/`start`/`rm` commands should succeed without errors. After `rm`, `checkc` should show the container does not exist.

If the container takes time to start (Liferay Portal is large), wait up to 2-3 minutes. If it fails to start, check Docker logs directly: `docker logs lpn-release`.

### 11. Deploy Command Validation

**Goal:** Verify deploy command input validation, error handling, and help output.

```bash
# Deploy without specifying files or dir (should error)
lpn deploy ce
```

**Observe:** Should print `Please pass a valid path to a file or to a directory as argument` and exit with code 1.

```bash
# Deploy with non-existent file (should error)
lpn deploy ce -f /tmp/nonexistent-file-12345.jar
```

**Observe:** Should fail with exit code 1.

```bash
# Deploy help
lpn deploy --help
```

**Observe:** Should show `Deploys files or a directory to Liferay Portal's deploy folder` and `Available Commands:`.

```bash
# Verify all deploy subcommands show help with flags --files/-f and --dir/-d
lpn deploy ce --help
lpn deploy commerce --help
lpn deploy dxp --help
lpn deploy nightly --help
lpn deploy release --help
```

**Observe:** Each subcommand help should contain `Deploys files or a directory` and list the `-f`/`--files` and `-d`/`--dir` flags.

### 11b. Deploy Valid File Without Running Container

**Goal:** Verify deploy fails gracefully when deploying a valid file to a non-running container for all flavors.

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

### 12. Deploy to a Running Container

**Goal:** Verify file deployment works when a container is running.

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

**Observe:** The deploy command should copy the file to the container's deploy folder. If the container is not running, it should fail gracefully.

### 13. Remove Image

**Goal:** Verify image removal.

```bash
# Remove a previously pulled image
lpn rmi release -t latest

# Verify it's gone
lpn checki release -t latest
```

**Observe:** The image should be removed. After removal, `checki` should report the image has not been pulled.

### 14. Run with Database Orchestration

**Goal:** Verify database container orchestration with the `--datastore` flag.

```bash
# Run with MySQL datastore
lpn run release -t latest -p 8080 --datastore mysql

# Check if both portal and database containers are running
lpn checkc release
docker ps --filter "label=lpn-container-name"

# Clean up
lpn rm release
```

**Observe:** Two containers should be created: one for the portal and one for MySQL. The portal should be configured to connect to the MySQL database. If testcontainers-go Ryuk is running, a third container (ryuk) may appear — this is normal.

Repeat with PostgreSQL if time permits:

```bash
lpn run release -t latest -p 8081 --datastore postgresql
lpn rm release
```

### 15. Run with Custom Flags

**Goal:** Verify all `run` flags are accepted and applied.

```bash
# Run with debug mode, custom ports, and memory setting
lpn run release -t latest -p 9090 -g 22222 --debug -D 8787 -m "-Xmx4096m"

# Verify the container is running
lpn checkc release

# Clean up
lpn rm release
```

**Observe:** The container should start with the specified ports and memory setting. The success log should show `httpPort=9090`, `gogoPort=22222`, `debug=true`, `debugPort=8787`, and `memory=-Xmx4096m`. Check `docker inspect` output to verify port mappings if needed.

Run flags reference:
- `-p`/`--httpPort` (int, default 8080) — HTTP port
- `-d`/`--debug` (bool, default false) — Enable debug mode
- `-D`/`--debugPort` (int, default 9000) — Debug port (only applies when debug enabled)
- `-g`/`--gogoPort` (int, default 11311) — GoGo Shell port
- `-s`/`--datastore` (string, default "hsql") — Database service: hsql, mysql, or postgresql
- `-t`/`--tag` (string) — Image tag to run
- `-m`/`--memory` (string) — JVM memory configuration (e.g., `-Xmx2048m`)

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

**Observe:** Should print a message about Equinox updates being disabled and provide a URL to download releases. Exit code should be 1 (since updates are disabled).

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

```bash
lpn stop ce
lpn start ce
lpn rm ce
```

**Observe:** Each should handle the missing container gracefully with a warning message rather than crashing. Exit code should be 0.

### 21. Open Command (Headless Environment)

**Goal:** Test the open command behavior in a headless environment.

```bash
lpn open ce
```

**Observe:** In a headless environment (CI/server), this will fail because there's no browser or display. It should fail with an error rather than hanging. On a desktop, it would open `http://localhost:<port>` in the default browser.

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
