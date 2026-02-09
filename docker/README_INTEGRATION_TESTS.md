# Testcontainers-go Integration Tests

This directory contains integration tests for database container management using [testcontainers-go](https://github.com/testcontainers/testcontainers-go).

## Overview

LPN now uses testcontainers-go for database container management with **label-based container identification** instead of hardcoded names. This approach:
- **Eliminates container name conflicts** - testcontainers auto-generates unique names
- **Uses labels for identification** - containers are found via the `lpn-container-name` label
- **Simplifies testing** - no cleanup code needed
- **Supports production** - disable Ryuk to prevent automatic cleanup

## Files

### Production Implementation (`docker.go`)

The `RunDatabaseDockerImage()` function uses testcontainers-go for database management:
- Containers identified by `lpn-container-name` label
- Auto-generated container names (no conflicts)
- Data persistence with bind mounts
- Built-in wait strategies ensure readiness

### Testing (`database_integration_test.go`)

Tests focusing on our code, not the library:
- **Unit Tests**: Test database selection, struct methods, aliases, constants
- **Integration Tests**: Test `RunDatabaseDockerImage()` with MySQL and PostgreSQL
- Containers identified by labels
- No cleanup code needed - testcontainers handles lifecycle

## Running the Tests

### Run all tests (recommended for CI/CD)
```bash
go test ./... -short
```
This runs unit tests and skips integration tests.

### Run with integration tests
```bash
go test -v ./docker -timeout 5m
```

### Run only unit tests
```bash
go test -v ./docker -short
```

## Container Identification

LPN uses label-based container identification instead of hardcoded names:

- **Label**: `lpn-container-name` stores the logical container name
- **Auto-generated names**: Testcontainers creates unique container names
- **No conflicts**: Multiple test runs don't conflict
- **Lookup**: All operations use `getContainersByLabel()` helper

## Why Testcontainers-go for Production?

### Benefits for Production Use
- ✅ **Label-based identification**: Eliminates container name conflicts
- ✅ **No automatic cleanup**: Ryuk can be disabled globally
- ✅ **Built-in modules**: Pre-configured for MySQL, PostgreSQL, etc.
- ✅ **Wait strategies**: Ensures services are ready before returning
- ✅ **Unified API**: Same library for testing and production

### Additional Production Configuration

**Disable Ryuk Reaper Globally:**

The Ryuk reaper is a sidecar container that automatically cleans up containers. For production use, you can disable it:

```bash
# Environment variable
export TESTCONTAINERS_RYUK_DISABLED=true
```

Or create `.testcontainers.properties`:
```properties
ryuk.disabled=true
```

This ensures no automatic cleanup happens.

### Benefits for Testing
- ✅ **Automatic cleanup**: Testcontainers handles container lifecycle
- ✅ **Port conflict avoidance**: Dynamic port allocation
- ✅ **Test isolation**: Each test gets fresh containers
- ✅ **Simplified setup**: No manual Docker commands
- ✅ **Label-based lookup**: No cleanup code needed

## Production vs Testing Pattern

### Production Pattern (Persistent Containers)

Containers persist by disabling Ryuk globally:

```bash
# Before running your application
export TESTCONTAINERS_RYUK_DISABLED=true
```

```go
// Container persists - Ryuk is disabled
container, err := mysql.Run(ctx, "mysql:8.0",
    mysql.WithDatabase("lportal"),
    testcontainers.WithLabels(map[string]string{
        "lpn-container-name": "my-db",  // Identified by label
    }),
    testcontainers.WithMounts(
        testcontainers.BindMount(volumePath, mountTarget),
    ),
)
```

### Testing Pattern (Ephemeral Containers)
```go
// Container is automatically cleaned up after test
container, err := mysql.Run(ctx, "mysql:8.0",
    mysql.WithDatabase("testdb"),
    testcontainers.WithLabels(map[string]string{
        "lpn-container-name": "test-db",
    }),
)
require.NoError(t, err)
// Container lifecycle managed by testcontainers
```

## Integration Test Pattern

Following testcontainers-go best practices with label-based identification:

```go
func TestDatabaseIntegration(t *testing.T) {
    ctx := context.Background()
    
    // Start container with label for identification
    container, err := mysql.Run(ctx, "mysql:8.0",
        mysql.WithDatabase("testdb"),
        testcontainers.WithLabels(map[string]string{
            "lpn-container-name": "test-mysql",
        }),
    )
    require.NoError(t, err)
    
    // Testcontainers handles cleanup automatically
    // Use container for testing...
}
```

## Production Container Pattern

```go
func StartProductionDatabase(image DatabaseImage) error {
    ctx := context.Background()
    containerName := image.GetContainerName()
    
    // Start container with label identification
    container, err := mysql.Run(ctx, 
        image.GetFullyQualifiedName(),
        mysql.WithDatabase(DBName),
        testcontainers.WithLabels(map[string]string{
            "lpn-container-name": containerName,  // Identify by label
            "db-type":            image.GetType(),
            "lpn-type":           "production",
        }),
        testcontainers.WithMounts(...),  // Data persists
    )
    
    return err
}
```

The key difference is the Ryuk reaper state (enabled for tests, disabled for production).

## More Information

For implementation details, see:
- `docker/docker.go` - Production implementation with label-based lookup
- `docker/database_integration_test.go` - Integration tests

## Dependencies

- `github.com/testcontainers/testcontainers-go v0.40.0`
- `github.com/testcontainers/testcontainers-go/modules/postgres v0.40.0`
- `github.com/testcontainers/testcontainers-go/modules/mysql v0.40.0`
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/go-sql-driver/mysql` - MySQL driver
