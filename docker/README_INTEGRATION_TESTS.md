# Testcontainers-go Integration Tests

This directory contains integration tests that demonstrate the use of [testcontainers-go](https://github.com/testcontainers/testcontainers-go) for testing database container functionality.

## Overview

These tests showcase a **hybrid approach** where:
- **Production code** continues to use the Docker client API directly for managing long-lived Liferay containers
- **Integration tests** use testcontainers-go for reliable, isolated, ephemeral test containers

## Running the Tests

### Run all tests (recommended for CI/CD)
```bash
go test ./... -short
```
This skips PostgreSQL tests that may have IPv6 networking issues in some environments.

### Run MySQL integration test
```bash
go test -v ./docker -run TestMySQLContainerIntegration -timeout 5m
```

### Run all integration tests (including PostgreSQL)
```bash
go test -v ./docker -timeout 5m
```

## Test Files

### `database_integration_test.go`

Contains three integration tests:

1. **TestMySQLContainerIntegration** ✅
   - Demonstrates MySQL container integration
   - Always passes in all environments
   - Shows automatic lifecycle management

2. **TestPostgreSQLContainerIntegration** ⚠️
   - Demonstrates PostgreSQL container integration
   - Skipped in short mode due to IPv6 networking issues
   - Container starts successfully, connection may fail in some environments

3. **TestPostgreSQLSnapshot**
   - Demonstrates advanced snapshot/restore features
   - Skipped in short mode
   - Shows unique testcontainers-go capabilities

## Why Testcontainers-go?

### Benefits for Testing
- ✅ **Automatic cleanup**: No leftover containers
- ✅ **Port conflict avoidance**: Dynamic port allocation
- ✅ **Test isolation**: Each test gets fresh containers
- ✅ **Simplified setup**: No manual Docker commands
- ✅ **Advanced features**: Snapshots, wait strategies, etc.

### Why Not for Production Code?
- ❌ **Designed for tests**: Automatic cleanup is counter to LPN's purpose
- ❌ **Ephemeral focus**: LPN manages persistent containers
- ❌ **API mismatch**: Test-oriented API vs. production container management

## Integration Test Pattern

Following testcontainers-go best practices:

```go
func TestDatabaseIntegration(t *testing.T) {
    ctx := context.Background()
    
    // Start container
    container, err := module.Run(ctx, "image:tag", options...)
    
    // CRITICAL: Register cleanup BEFORE error check
    // This prevents resource leaks
    testcontainers.CleanupContainer(t, container)
    
    // Then check for errors
    require.NoError(t, err)
    
    // Use container for testing
    // ...
}
```

## More Information

See [TESTCONTAINERS_EVALUATION.md](../docs/TESTCONTAINERS_EVALUATION.md) for:
- Complete evaluation rationale
- Architectural analysis
- Code comparisons
- Detailed recommendations

## Dependencies

- `github.com/testcontainers/testcontainers-go v0.40.0`
- `github.com/testcontainers/testcontainers-go/modules/postgres v0.40.0`
- `github.com/testcontainers/testcontainers-go/modules/mysql v0.40.0`
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/go-sql-driver/mysql` - MySQL driver

## Known Issues

### PostgreSQL IPv6 Networking
Some test environments have IPv6 localhost connection issues with PostgreSQL containers. The tests handle this by:
- Skipping in short mode (`go test -short`)
- Passing in environments with proper IPv6 configuration
- Not affecting production code at all

To run despite this issue:
```bash
# Skip problematic tests
go test ./... -short

# Or run specific working tests
go test -v ./docker -run TestMySQLContainerIntegration
```
