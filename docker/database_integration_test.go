package docker

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// TestPostgreSQLContainerIntegration demonstrates using testcontainers-go
// to validate PostgreSQL container functionality. This shows how testcontainers-go
// provides a cleaner, more reliable approach for integration testing compared to
// managing container lifecycle manually.
func TestPostgreSQLContainerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PostgreSQL integration test in short mode due to IPv6 networking issues")
	}
	
	ctx := context.Background()

	// Start PostgreSQL container with testcontainers-go
	// This handles all container lifecycle automatically
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(DBName),
		postgres.WithUsername(DBUser),
		postgres.WithPassword(DBPassword),
	)
	// CRITICAL: Register cleanup BEFORE error check to prevent resource leaks
	// The cleanup function handles nil containers gracefully
	testcontainers.CleanupContainer(t, pgContainer)
	require.NoError(t, err)

	// Get connection string - testcontainers handles port mapping automatically
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Verify database connectivity
	err = db.Ping()
	require.NoError(t, err)

	// Verify database exists and credentials work
	var dbName string
	err = db.QueryRow("SELECT current_database()").Scan(&dbName)
	require.NoError(t, err)
	assert.Equal(t, DBName, dbName)

	// Verify we can create tables and insert data
	_, err = db.Exec(`
		CREATE TABLE test_table (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO test_table (name) VALUES ($1)", "testdata")
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

// TestMySQLContainerIntegration demonstrates using testcontainers-go
// to validate MySQL container functionality. This provides automatic
// container lifecycle management and port allocation.
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
	// CRITICAL: Register cleanup BEFORE error check to prevent resource leaks
	// The cleanup function handles nil containers gracefully
	testcontainers.CleanupContainer(t, mysqlContainer)
	require.NoError(t, err)

	// Get connection string
	connStr, err := mysqlContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Connect to the database
	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Verify database connectivity
	err = db.Ping()
	require.NoError(t, err)

	// Verify database exists
	var dbName string
	err = db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	require.NoError(t, err)
	assert.Equal(t, DBName, dbName)

	// Verify we can create tables and insert data
	_, err = db.Exec(`
		CREATE TABLE test_table (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO test_table (name) VALUES (?)", "testdata")
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

// TestPostgreSQLSnapshot demonstrates testcontainers-go's advanced features
// like database snapshots for test isolation. This is a unique feature that
// helps with maintaining clean test state.
func TestPostgreSQLSnapshot(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PostgreSQL snapshot test in short mode due to IPv6 networking issues")
	}
	
	ctx := context.Background()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(DBName),
		postgres.WithUsername(DBUser),
		postgres.WithPassword(DBPassword),
	)
	// CRITICAL: Register cleanup BEFORE error check to prevent resource leaks
	testcontainers.CleanupContainer(t, pgContainer)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Create initial data
	_, err = db.Exec(`
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO users (name) VALUES ($1)", "Alice")
	require.NoError(t, err)

	// Take a snapshot
	err = pgContainer.Snapshot(ctx, postgres.WithSnapshotName("initial"))
	require.NoError(t, err)

	// Make changes
	_, err = db.Exec("INSERT INTO users (name) VALUES ($1)", "Bob")
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Restore to snapshot
	err = pgContainer.Restore(ctx, postgres.WithSnapshotName("initial"))
	require.NoError(t, err)

	// Reconnect after restore
	db.Close()
	db, err = sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Verify Bob is gone, only Alice remains
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	var name string
	err = db.QueryRow("SELECT name FROM users LIMIT 1").Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "Alice", name)
}
