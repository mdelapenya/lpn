package docker

import (
	"testing"

	"github.com/mdelapenya/lpn/internal"
	"github.com/stretchr/testify/require"
	liferay "github.com/mdelapenya/lpn/liferay"
)

// setupTestConfig initializes the global LpnConfig for tests
func setupTestConfig() {
	if internal.LpnConfig == nil {
		// Set up a temporary workspace for tests
		tmpDir := "/tmp/lpn-test-workspace"
		internal.LpnWorkspace = tmpDir
		
		internal.LpnConfig = &internal.LPNConfig{
			Container: internal.NamesConfig{
				Names: internal.NameConfig{
					Db: map[string]string{
						"ce":         "db-ce",
						"commerce":   "db-commerce",
						"dxp":        "db-dxp",
						"nightly":    "db-nightly",
						"release":    "db-release",
						"test":       "db-test",       // Add test type
						"test-reuse": "db-test-reuse", // Add test-reuse type
					},
					Portal: map[string]string{
						"ce":         "lpn-ce",
						"commerce":   "lpn-commerce",
						"dxp":        "lpn-dxp",
						"nightly":    "lpn-nightly",
						"release":    "lpn-release",
						"test":       "lpn-test",
						"test-reuse": "lpn-test-reuse",
					},
				},
			},
			Images: internal.ImagesConfig{
				Db: map[string]internal.ImageConfig{
					"mysql": {
						Image: "docker.io/mdelapenya/mysql-utf8",
						Tag:   "5.7",
					},
					"postgres": {
						Image: "postgres",
						Tag:   "9.6-alpine",
					},
				},
				Portal: map[string]internal.ImageConfig{
					"ce": {
						Image: "liferay/portal",
						Tag:   "7.0.6-ga7",
					},
				},
			},
		}
	}
}

// TestGetDatabase tests the GetDatabase function returns correct database types
func TestGetDatabase(t *testing.T) {
	mockImage := liferay.CE{Tag: "latest"}

	tests := []struct {
		name      string
		datastore string
		wantType  string
		wantNil   bool
	}{
		{
			name:      "MySQL datastore",
			datastore: "mysql",
			wantType:  "mysql",
			wantNil:   false,
		},
		{
			name:      "PostgreSQL datastore",
			datastore: "postgresql",
			wantType:  "postgresql",
			wantNil:   false,
		},
		{
			name:      "Unsupported datastore",
			datastore: "mongodb",
			wantType:  "",
			wantNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := GetDatabase(mockImage, tt.datastore)
			if tt.wantNil {
				require.Nil(t, db)
			} else {
				require.NotNil(t, db)
				require.Equal(t, tt.wantType, db.GetType())
			}
		})
	}
}

// TestMySQLImage tests MySQL struct methods
func TestMySQLImage(t *testing.T) {
	mysql := MySQL{
		LpnType: "ce",
		Tag:     "5.7",
	}

	t.Run("GetType", func(t *testing.T) {
		require.Equal(t, "mysql", mysql.GetType())
	})

	t.Run("GetPort", func(t *testing.T) {
		require.Equal(t, 3306, mysql.GetPort())
	})

	t.Run("GetDataFolder", func(t *testing.T) {
		require.Equal(t, "/var/lib/mysql", mysql.GetDataFolder())
	})

	t.Run("GetTag", func(t *testing.T) {
		require.Equal(t, "5.7", mysql.GetTag())
	})

	t.Run("GetLpnType", func(t *testing.T) {
		require.Equal(t, "ce", mysql.GetLpnType())
	})

	t.Run("GetJDBCConnection", func(t *testing.T) {
		jdbc := mysql.GetJDBCConnection()
		require.Equal(t, "com.mysql.jdbc.Driver", jdbc.DriverClassName)
		require.Equal(t, DBUser, jdbc.User)
		require.Equal(t, DBPassword, jdbc.Password)
		require.Contains(t, jdbc.URL, "jdbc:mysql://")
		require.Contains(t, jdbc.URL, DBName)
	})

	t.Run("GetEnvVariables", func(t *testing.T) {
		env := mysql.GetEnvVariables()
		require.Equal(t, "MYSQL_DATABASE="+DBName, env.Database)
		require.Equal(t, "MYSQL_ROOT_PASSWORD="+DBPassword, env.Password)
		require.Equal(t, "MYSQL_USER="+DBUser, env.User)
	})
}

// TestPostgreSQLImage tests PostgreSQL struct methods
func TestPostgreSQLImage(t *testing.T) {
	postgres := PostgreSQL{
		LpnType: "ce",
		Tag:     "16-alpine",
	}

	t.Run("GetType", func(t *testing.T) {
		require.Equal(t, "postgresql", postgres.GetType())
	})

	t.Run("GetPort", func(t *testing.T) {
		require.Equal(t, 5432, postgres.GetPort())
	})

	t.Run("GetDataFolder", func(t *testing.T) {
		require.Equal(t, "/var/lib/postgresql/data", postgres.GetDataFolder())
	})

	t.Run("GetTag", func(t *testing.T) {
		require.Equal(t, "16-alpine", postgres.GetTag())
	})

	t.Run("GetLpnType", func(t *testing.T) {
		require.Equal(t, "ce", postgres.GetLpnType())
	})

	t.Run("GetJDBCConnection", func(t *testing.T) {
		jdbc := postgres.GetJDBCConnection()
		require.Equal(t, "org.postgresql.Driver", jdbc.DriverClassName)
		require.Equal(t, DBUser, jdbc.User)
		require.Equal(t, DBPassword, jdbc.Password)
		require.Contains(t, jdbc.URL, "jdbc:postgresql://")
		require.Contains(t, jdbc.URL, DBName)
	})

	t.Run("GetEnvVariables", func(t *testing.T) {
		env := postgres.GetEnvVariables()
		require.Equal(t, "POSTGRES_DB="+DBName, env.Database)
		require.Equal(t, "POSTGRES_PASSWORD="+DBPassword, env.Password)
		require.Equal(t, "POSTGRES_USER="+DBUser, env.User)
	})
}

// TestGetAlias tests the GetAlias function
func TestGetAlias(t *testing.T) {
	require.Equal(t, "db", GetAlias())
}

// TestDatabaseConstants tests the database constants
func TestDatabaseConstants(t *testing.T) {
	require.Equal(t, "lportal", DBName)
	require.Equal(t, "my-secret-pw", DBPassword)
	require.Equal(t, "liferay", DBUser)
}

// TestRunDatabaseDockerImageMySQL is an integration test for MySQL container creation
// This tests our actual RunDatabaseDockerImage function with MySQL
func TestRunDatabaseDockerImageMySQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Initialize test config
	setupTestConfig()

	mysql := MySQL{
		LpnType: "test",
		Tag:     "5.7",
	}

	// Clean up any existing container first
	containerName := mysql.GetContainerName()
	if CheckDockerContainerExists(containerName) {
		t.Logf("Container %s already exists, test may reuse it", containerName)
	}

	// Test our RunDatabaseDockerImage function
	err := RunDatabaseDockerImage(mysql)
	require.NoError(t, err, "RunDatabaseDockerImage should not error")

	// Verify container was created
	require.True(t, CheckDockerContainerExists(containerName), 
		"Container %s should exist after RunDatabaseDockerImage", containerName)

	t.Cleanup(func() {
		t.Logf("Test completed for container %s", containerName)
	})
}

// TestRunDatabaseDockerImagePostgreSQL is an integration test for PostgreSQL container creation
// This tests our actual RunDatabaseDockerImage function with PostgreSQL
func TestRunDatabaseDockerImagePostgreSQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Initialize test config
	setupTestConfig()

	postgres := PostgreSQL{
		LpnType: "test",
		Tag:     "16-alpine",
	}

	// Clean up any existing container first
	containerName := postgres.GetContainerName()
	if CheckDockerContainerExists(containerName) {
		t.Logf("Container %s already exists, test may reuse it", containerName)
	}

	// Test our RunDatabaseDockerImage function
	err := RunDatabaseDockerImage(postgres)
	require.NoError(t, err, "RunDatabaseDockerImage should not error")

	// Verify container was created
	require.True(t, CheckDockerContainerExists(containerName),
		"Container %s should exist after RunDatabaseDockerImage", containerName)

	t.Cleanup(func() {
		t.Logf("Test completed for container %s", containerName)
	})
}

// TestRunDatabaseDockerImageAlreadyExists tests that calling RunDatabaseDockerImage
// on an already running container doesn't fail
func TestRunDatabaseDockerImageAlreadyExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Initialize test config
	setupTestConfig()

	mysql := MySQL{
		LpnType: "test-reuse",
		Tag:     "5.7",
	}

	// First call should create the container
	err := RunDatabaseDockerImage(mysql)
	require.NoError(t, err, "First call to RunDatabaseDockerImage should not error")

	// Second call should not fail (should skip creation and reuse existing container)
	err = RunDatabaseDockerImage(mysql)
	require.NoError(t, err, "Second call to RunDatabaseDockerImage should not error")

	// Verify container still exists
	require.True(t, CheckDockerContainerExists(mysql.GetContainerName()),
		"Container should still exist after second call")

	t.Cleanup(func() {
		t.Logf("Test completed for container %s", mysql.GetContainerName())
	})
}
