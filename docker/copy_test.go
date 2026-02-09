package docker

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestCopyFileToContainer tests the CopyFileToContainer function with a real container
func TestCopyFileToContainer(t *testing.T) {
	setupTestConfig()
	ctx := context.Background()

	// Start a simple container for testing (using alpine as it's lightweight)
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "alpine:latest",
			Cmd:   []string{"sleep", "300"},
			Labels: map[string]string{
				"lpn-container-name": "test-deploy-container",
			},
			WaitingFor: wait.ForLog(""),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer func() {
		err := container.Terminate(ctx)
		if err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	// Create a temporary test file to deploy
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-deploy.txt")
	testContent := []byte("This is a test file for deployment")
	err = os.WriteFile(testFile, testContent, 0644)
	require.NoError(t, err)

	// Create a mock Liferay image that points to our test container
	mockImage := &mockLiferayImage{
		containerName: "test-deploy-container",
		deployFolder:  "/tmp/deploy",
		user:          "root",
	}

	// First, create the deploy folder in the container
	_, _, err = container.Exec(ctx, []string{"mkdir", "-p", "/tmp/deploy"})
	require.NoError(t, err)

	// Test: Copy file to container
	err = CopyFileToContainer(mockImage, testFile)
	assert.NoError(t, err, "CopyFileToContainer should succeed")

	// Verify: Check that the file exists in the container using multiplexed Exec
	targetPath := filepath.Join(mockImage.deployFolder, "test-deploy.txt")
	exitCode, reader, err := container.Exec(ctx, []string{"test", "-f", targetPath}, tcexec.Multiplexed())
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode, "File should exist in container")

	// Verify: Read the file content using multiplexed Exec
	exitCode, reader, err = container.Exec(ctx, []string{"cat", targetPath}, tcexec.Multiplexed())
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode, "Should be able to read the file")
	
	// Read the output
	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, string(testContent), string(content), "File content should match exactly")
}

// TestDeployTextFileToContainer tests deploying a text file to a container and verifying it with multiplexed Exec
func TestDeployTextFileToContainer(t *testing.T) {
	setupTestConfig()
	ctx := context.Background()

	// Start a simple container for testing (using alpine as it's lightweight)
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "alpine:latest",
			Cmd:   []string{"sleep", "300"},
			Labels: map[string]string{
				"lpn-container-name": "test-deploy-text-container",
			},
			WaitingFor: wait.ForLog(""),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer func() {
		err := container.Terminate(ctx)
		if err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	// Create a temporary test file with specific content
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "deployment-test.txt")
	testContent := "Hello from LPN deployment test!\nThis is line 2.\nThis is line 3."
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create a mock Liferay image that points to our test container
	mockImage := &mockLiferayImage{
		containerName: "test-deploy-text-container",
		deployFolder:  "/tmp/deploy",
		user:          "root",
	}

	// First, create the deploy folder in the container using multiplexed Exec
	exitCode, reader, err := container.Exec(ctx, []string{"mkdir", "-p", "/tmp/deploy"}, tcexec.Multiplexed())
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "mkdir should succeed")
	io.ReadAll(reader) // drain the reader

	// Deploy the file to the container
	err = CopyFileToContainer(mockImage, testFile)
	require.NoError(t, err, "CopyFileToContainer should succeed")

	// Verify: Check that the file exists using multiplexed Exec
	targetPath := filepath.Join(mockImage.deployFolder, "deployment-test.txt")
	exitCode, reader, err = container.Exec(ctx, []string{"test", "-f", targetPath}, tcexec.Multiplexed())
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode, "Deployed file should exist in container at "+targetPath)
	io.ReadAll(reader) // drain the reader

	// Verify: Read and check the file content using multiplexed Exec
	exitCode, reader, err = container.Exec(ctx, []string{"cat", targetPath}, tcexec.Multiplexed())
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode, "Should be able to read the deployed file")
	
	// Read and verify the content
	deployedContent, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, testContent, string(deployedContent), "Deployed file content should match original")

	// Verify: Check file permissions using multiplexed Exec
	exitCode, reader, err = container.Exec(ctx, []string{"stat", "-c", "%a", targetPath}, tcexec.Multiplexed())
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode, "Should be able to stat the deployed file")
	
	permissions, err := io.ReadAll(reader)
	require.NoError(t, err)
	t.Logf("Deployed file permissions: %s", string(permissions))

	// Verify: Check file ownership using multiplexed Exec
	exitCode, reader, err = container.Exec(ctx, []string{"stat", "-c", "%U:%G", targetPath}, tcexec.Multiplexed())
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode, "Should be able to check file ownership")
	
	ownership, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Contains(t, string(ownership), "root", "File should be owned by root user")
}

// TestCopyFileToContainerNonExistentFile tests error handling for non-existent files
func TestCopyFileToContainerNonExistentFile(t *testing.T) {
	setupTestConfig()
	ctx := context.Background()

	// Start a simple container for testing
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "alpine:latest",
			Cmd:   []string{"sleep", "300"},
			Labels: map[string]string{
				"lpn-container-name": "test-deploy-error-container",
			},
			WaitingFor: wait.ForLog(""),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer func() {
		err := container.Terminate(ctx)
		if err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	mockImage := &mockLiferayImage{
		containerName: "test-deploy-error-container",
		deployFolder:  "/tmp/deploy",
		user:          "root",
	}

	// Test: Try to copy a non-existent file
	err = CopyFileToContainer(mockImage, "/tmp/nonexistent-file-12345.txt")
	assert.Error(t, err, "Should return error for non-existent file")
}

// TestCopyFileToContainerNonExistentContainer tests error handling for non-existent container
func TestCopyFileToContainerNonExistentContainer(t *testing.T) {
	setupTestConfig()

	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	mockImage := &mockLiferayImage{
		containerName: "non-existent-container-12345",
		deployFolder:  "/tmp/deploy",
		user:          "root",
	}

	// Test: Try to copy to non-existent container
	err = CopyFileToContainer(mockImage, testFile)
	assert.Error(t, err, "Should return error for non-existent container")
}

// mockLiferayImage is a mock implementation of liferay.Image for testing
type mockLiferayImage struct {
	containerName string
	deployFolder  string
	user          string
}

func (m *mockLiferayImage) GetContainerName() string {
	return m.containerName
}

func (m *mockLiferayImage) GetDeployFolder() string {
	return m.deployFolder
}

func (m *mockLiferayImage) GetUser() string {
	return m.user
}

// Implement other required methods of liferay.Image interface
func (m *mockLiferayImage) GetFullyQualifiedName() string {
	return "mock/image:test"
}

func (m *mockLiferayImage) GetRepository() string {
	return "mock/image"
}

func (m *mockLiferayImage) GetTag() string {
	return "test"
}

func (m *mockLiferayImage) GetType() string {
	return "ce"
}

func (m *mockLiferayImage) GetDockerHubTagsURL() string {
	return "https://hub.docker.com/mock"
}

func (m *mockLiferayImage) GetLiferayHome() string {
	return "/opt/liferay"
}
