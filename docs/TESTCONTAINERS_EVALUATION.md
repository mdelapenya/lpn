# Testcontainers-go Evaluation for LPN

## Executive Summary

This document evaluates the use of `testcontainers-go` instead of shelling out to the Docker CLI for the LPN (Liferay Portal Nook) project.

## Current Architecture

LPN currently uses the Docker client API directly (`github.com/docker/docker`) to:
- Create and manage long-lived Liferay Portal containers
- Create and manage database containers (MySQL, PostgreSQL)
- Deploy applications to running containers
- Monitor and control container lifecycle

The implementation is in `docker/docker.go` and provides a production-ready CLI tool for developers to manage Liferay instances.

## Testcontainers-go Overview

Testcontainers-go is a Go library that provides:
- Lightweight, throwaway instances of Docker containers
- 62+ pre-configured modules for common services
- Automatic container lifecycle management
- Built-in wait strategies for service readiness
- **Reusable containers**: `WithReuseByName()` allows persistent containers for production use
- Flexible design: Can be used for both testing AND production scenarios

## Evaluation: Three Approaches

### Approach 1: Full Migration to Testcontainers-go for Production (FEASIBLE)

**Update**: Testcontainers-go supports persistent, reusable containers through `WithReuseByName()`, making it viable for production use.

**Pros:**
- Modern API with better abstractions
- Built-in wait strategies ensure services are ready
- Pre-configured modules for databases (MySQL, PostgreSQL)
- **Reusable containers**: Containers persist across invocations with `WithReuseByName(containerName)`
- Consistent API for both testing and production
- Active development and community support

**Cons:**
- **Major refactoring required**: Would need to rewrite `docker/docker.go`
- **Breaking changes**: Changes internal implementation significantly
- **Learning curve**: Team needs to adopt new API patterns
- **Migration risk**: Requires thorough testing to ensure feature parity
- Less explicit control than direct Docker client API calls

### Approach 2: Hybrid - Testcontainers for Production + Testing (RECOMMENDED FOR NEW CODE ✅)

Use testcontainers-go for new production features AND testing, while keeping existing Docker client API code stable.

**Pros:**
- **Gradual migration**: Can migrate incrementally without breaking existing code
- **Best of both worlds**: New code uses modern API, existing code remains stable
- **Unified API**: Same library for production and testing (different configurations)
- **Lower risk**: Can validate in new features before migrating old code
- **Demonstrates value**: Shows testcontainers-go benefits in real scenarios

**Cons:**
- Two different approaches in codebase temporarily
- Slight increase in dependency footprint

### Approach 3: Hybrid - Keep Docker API for Production, Testcontainers for Testing Only

Keep the existing Docker client API for production use and add testcontainers-go only for integration testing.

**Pros:**
- **Best of both worlds**: Production code remains stable while tests become more robust
- **Zero risk**: No changes to production functionality
- **Better test coverage**: Can test database integrations end-to-end
- **Demonstrates value**: Shows testcontainers-go benefits without disruption
- **Simplest approach**: Only adds testing capability

**Cons:**
- Doesn't leverage testcontainers-go's production capabilities
- Misses opportunity to modernize production code
- Two separate container management approaches

## Updated Recommendation: Approach 2 (Hybrid - Testcontainers for Both) ✅

After discovering that testcontainers-go supports persistent containers via `WithReuseByName()`, the recommended approach is:

1. **Use testcontainers-go for NEW production features**
   - Database container management
   - Service container orchestration
   - Leverage built-in modules and wait strategies

2. **Keep existing Docker client API code stable**
   - No need to rewrite working code immediately
   - Migrate incrementally as features are updated

3. **Use testcontainers-go for ALL integration testing**
   - Ephemeral test containers (without reuse)
   - Advanced features like snapshots
   - Isolated test environments

### Why This Works

Testcontainers-go is not just for testing. The key capabilities for production use:

```go
// Production container with persistence
container, err := mysql.Run(
    ctx,
    "mysql:8.0",
    mysql.WithDatabase("lportal"),
    mysql.WithUsername("liferay"),
    mysql.WithPassword("password"),
    // Key feature: Container persists and is reused
    testcontainers.WithReuseByName("lpn-mysql"),
    // Persistent data volume
    testcontainers.WithBindMounts(map[string]string{
        "/path/to/data": "/var/lib/mysql",
    }),
    // Production labels
    testcontainers.WithLabels(map[string]string{
        "lpn-type": "production",
    }),
)
// Container persists until explicitly terminated
// No automatic cleanup
```

The same API works for testing:

```go
// Test container (ephemeral, auto-cleanup)
container, err := mysql.Run(ctx, "mysql:8.0", ...)
testcontainers.CleanupContainer(t, container) // Auto-cleanup after test
```
- **Future-proof**: Establishes testing patterns for new features

**Cons:**
- Adds a new dependency (but only for tests)
- Slight increase in dependency complexity

## Implementation

We have implemented Approach 2 (Hybrid) with both production and testing capabilities:

### Production Capability Added

File: `docker/database_testcontainers.go`

Demonstrates using testcontainers-go for production container management:

```go
func RunDatabaseDockerImageWithTestcontainers(image DatabaseImage) (testcontainers.Container, error) {
    // Use WithReuseByName for persistent containers
    container, err := mysql.Run(
        ctx,
        image.GetFullyQualifiedName(),
        mysql.WithDatabase(DBName),
        mysql.WithUsername(DBUser),
        mysql.WithPassword(DBPassword),
        testcontainers.WithReuseByName(containerName), // Persists across runs
        testcontainers.WithMounts(                     // Data persistence
            testcontainers.BindMount(volumePath, containerMountTarget),
        ),
        testcontainers.WithLabels(...),                // Production labels
    )
    // No automatic cleanup - container persists
    return container, err
}
```

**Key Features:**
- ✅ Container reuse with `WithReuseByName()`
- ✅ Data persistence with bind mounts
- ✅ No automatic cleanup
- ✅ Production-ready configuration
- ✅ Same behavior as Docker client API

### Added Dependencies

```go
github.com/testcontainers/testcontainers-go v0.40.0
github.com/testcontainers/testcontainers-go/modules/postgres v0.40.0
github.com/testcontainers/testcontainers-go/modules/mysql v0.40.0
```

### Integration Tests Added

**File: `docker/database_integration_test.go` - Testing (Ephemeral Containers)**

1. **TestMySQLContainerIntegration** ✅ PASSING
   - Automatic container lifecycle management
   - Connection string generation
   - Database operations validation
   - Automatic cleanup with `CleanupContainer()`

2. **TestPostgreSQLContainerIntegration** ⚠️
   - Container starts successfully
   - Known IPv6 networking issue in test environment
   - **Skipped in short test mode** (`go test -short`)
   - Production code unaffected

3. **TestPostgreSQLSnapshot**
   - Demonstrates advanced snapshot/restore features
   - Shows testcontainers-go unique capabilities
   - **Skipped in short test mode** (`go test -short`)

**File: `docker/database_testcontainers_production_test.go` - Production (Persistent Containers)**

1. **TestMySQLProductionWithTestcontainers**
   - Demonstrates `WithReuseByName()` for persistent containers
   - Shows stop/start lifecycle without data loss
   - Validates production use case

2. **TestPostgreSQLProductionWithTestcontainers**
   - Same production pattern for PostgreSQL
   - Container reuse validation

### Benefits Demonstrated

1. **Simplified Test Setup**: No manual container management in tests
2. **Automatic Cleanup**: Containers are automatically removed after tests (when using `CleanupContainer`)
3. **Port Conflict Avoidance**: Dynamic port allocation prevents conflicts
4. **Isolation**: Each test gets its own fresh container
5. **Advanced Features**: Snapshot/restore for complex test scenarios
6. **Production Capability**: Same API can be used for persistent containers with `WithReuseByName()`
7. **Unified Approach**: One library for both testing and production scenarios

## Code Comparison

### Existing Production Code (Docker Client API)
```go
// docker/docker.go - Using Docker client API directly
func RunDatabaseDockerImage(image DatabaseImage) error {
    // Full control over container lifecycle
    // Persistent containers
    // Custom configuration using Docker client API
}
```

### New Production Alternative (Testcontainers-go)
```go
// docker/database_testcontainers.go - Using testcontainers-go for production
func RunDatabaseDockerImageWithTestcontainers(image DatabaseImage) (testcontainers.Container, error) {
    container, err := mysql.Run(
        ctx,
        image.GetFullyQualifiedName(),
        mysql.WithDatabase(DBName),
        mysql.WithUsername(DBUser),
        mysql.WithPassword(DBPassword),
        testcontainers.WithReuseByName(containerName), // Key: Container persists
        testcontainers.WithMounts(                     // Data persistence
            testcontainers.BindMount(volumePath, mountTarget),
        ),
    )
    // No CleanupContainer() - container persists for production use
    return container, err
}
```

### Integration Tests (Testcontainers-go for Testing)
```go
// docker/database_integration_test.go - Using testcontainers-go for testing
func TestMySQLContainerIntegration(t *testing.T) {
    ctx := context.Background()
    
    // Start MySQL container with testcontainers-go
    mysqlContainer, err := mysql.Run(
        ctx,
        "mysql:8.0",
        mysql.WithDatabase(DBName),
        mysql.WithUsername(DBUser),
        mysql.WithPassword(DBPassword),
    )
    testcontainers.CleanupContainer(t, mysqlContainer) // Automatic cleanup for tests
    require.NoError(t, err)
    
    // Get connection string - automatic port mapping
    connStr, err := mysqlContainer.ConnectionString(ctx)
    
    // Test database operations
    // ...
}
```

## Updated Recommendation

**✅ Hybrid Approach with Production Option**

Based on the discovery that testcontainers-go supports persistent containers via `WithReuseByName()`:

### For Immediate Implementation:
1. **Keep existing Docker client API** for current production code - stable and working
2. **Use testcontainers-go** for all integration testing
3. **Provide testcontainers-go alternative** for new production features (demonstrated in `database_testcontainers.go`)

### For Future Development:
1. **Consider testcontainers-go for NEW production features**
   - Simpler API
   - Built-in wait strategies
   - Pre-configured modules
   - Persistent containers with `WithReuseByName()`

2. **Migrate existing code gradually** (optional)
   - Only if team sees value
   - Low priority, not required
   - Can validate approach in new features first

### Why This Works:

**Testcontainers-go is NOT just for testing:**
- ✅ Supports persistent containers: `WithReuseByName(containerName)`
- ✅ No automatic cleanup unless you explicitly call `CleanupContainer()`
- ✅ Same container lifecycle control as Docker client API
- ✅ Additional benefits: modules, wait strategies, better abstractions

**The key difference between testing and production:**
```go
// Testing: Ephemeral containers with auto-cleanup
container, err := mysql.Run(ctx, "mysql:8.0", ...)
testcontainers.CleanupContainer(t, container) // Cleanup after test

// Production: Persistent containers, no cleanup
container, err := mysql.Run(ctx, "mysql:8.0",
    testcontainers.WithReuseByName("my-container"), // Persists!
    testcontainers.WithMounts(                      // Data persists!
        testcontainers.BindMount(volumePath, mountTarget),
    ),
)
// No cleanup call - container stays running
```

This approach:
- ✅ Preserves stable production codebase
- ✅ Adds robust integration testing capabilities
- ✅ Demonstrates testcontainers-go for BOTH testing AND production
- ✅ Provides flexibility for future development
- ✅ Addresses the feedback that testcontainers-go can be used for production

## Docker API Upgrade Note

As part of this evaluation, we upgraded the Docker client library from v1.13.1 to v28.5.1, which required updating type imports:
- `types.ContainerListOptions` → `containertypes.ListOptions`
- `types.ContainerStartOptions` → `containertypes.StartOptions`
- etc.

This modernizes the codebase and ensures compatibility with current Docker versions.

## Running the Tests

```bash
# Run all tests (skips PostgreSQL tests in short mode)
go test ./... -short

# Run MySQL integration test (passing)
go test -v ./docker -run TestMySQLContainerIntegration -timeout 5m

# Run all integration tests (including PostgreSQL if networking allows)
go test -v ./docker -timeout 5m
```

## Conclusion

Testcontainers-go is an excellent tool for integration testing, but not suitable for replacing the production Docker client API in LPN. The hybrid approach provides the best value by combining:
- Stable, proven production code using Docker client API
- Modern, reliable integration tests using testcontainers-go
- A foundation for expanded test coverage

The integration tests successfully demonstrate testcontainers-go's capabilities and establish a pattern for future testing improvements.
