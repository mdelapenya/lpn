package docker

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMySQLProductionWithTestcontainers demonstrates using testcontainers-go
// for production-like container management with reusable, persistent containers.
// This shows that testcontainers-go can be used beyond just testing.
func TestMySQLProductionWithTestcontainers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping production testcontainers demo in short mode")
	}

	ctx := context.Background()

	// Create a MySQL database image configuration
	mysqlImage := MySQL{
		LpnType: "test-production",
		Tag:     "8.0",
	}

	// Start the container using testcontainers-go with reuse enabled
	// This container will persist and can be reused across multiple runs
	container, err := RunDatabaseDockerImageWithTestcontainers(mysqlImage)
	require.NoError(t, err)
	require.NotNil(t, container)

	// Get the container's endpoint
	endpoint, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	// Connect to the database
	connStr := "liferay:my-secret-pw@tcp(" + endpoint + ")/lportal"
	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Verify we can use the database
	err = db.Ping()
	require.NoError(t, err)

	// Create a table and insert data
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS production_test (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO production_test (name) VALUES (?)", "production-data")
	require.NoError(t, err)

	// Verify data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM production_test").Scan(&count)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 1)

	// Stop the container (but don't remove it - it can be reused)
	err = StopDatabaseContainerWithTestcontainers(container)
	require.NoError(t, err)

	// In a real production scenario, the container would be restarted later
	// with the same name and the data would persist

	// Cleanup for this test (in production, you wouldn't do this)
	// Start it again briefly to terminate it
	container2, err := RunDatabaseDockerImageWithTestcontainers(mysqlImage)
	require.NoError(t, err)
	err = RemoveDatabaseContainerWithTestcontainers(container2)
	require.NoError(t, err)
}

// TestPostgreSQLProductionWithTestcontainers demonstrates PostgreSQL production usage
func TestPostgreSQLProductionWithTestcontainers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping production testcontainers demo in short mode")
	}

	ctx := context.Background()

	// Create a PostgreSQL database image configuration
	pgImage := PostgreSQL{
		LpnType: "test-production",
		Tag:     "16-alpine",
	}

	// Start the container using testcontainers-go with reuse enabled
	container, err := RunDatabaseDockerImageWithTestcontainers(pgImage)
	require.NoError(t, err)
	require.NotNil(t, container)

	// Verify container is running
	state, err := container.State(ctx)
	require.NoError(t, err)
	assert.True(t, state.Running)

	// Stop without removing
	err = StopDatabaseContainerWithTestcontainers(container)
	require.NoError(t, err)

	// Cleanup
	container2, err := RunDatabaseDockerImageWithTestcontainers(pgImage)
	require.NoError(t, err)
	err = RemoveDatabaseContainerWithTestcontainers(container2)
	require.NoError(t, err)
}
