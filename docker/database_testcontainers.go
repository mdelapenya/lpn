// Copyright (c) 2000-present Liferay, Inc. All rights reserved.
//
// This library is free software; you can redistribute it and/or modify it under
// the terms of the GNU Lesser General Public License as published by the Free
// Software Foundation; either version 2.1 of the License, or (at your option)
// any later version.
//
// This library is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more
// details.

package docker

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	internal "github.com/mdelapenya/lpn/internal"
	log "github.com/sirupsen/logrus"
)

// RunDatabaseDockerImageWithTestcontainers demonstrates using testcontainers-go
// for production container management with persistent, reusable containers.
// This is an alternative implementation to RunDatabaseDockerImage.
//
// Key features:
// - WithReuseByName: Container persists and is reused across invocations
// - No automatic cleanup: Container lifecycle managed explicitly
// - Production-ready: Same behavior as Docker client API implementation
//
// For additional persistence guarantees, you can disable Ryuk globally:
// - Set TESTCONTAINERS_RYUK_DISABLED=true environment variable, OR
// - Create .testcontainers.properties with ryuk.disabled=true
// This prevents the reaper sidecar container from auto-removing containers.
func RunDatabaseDockerImageWithTestcontainers(image DatabaseImage) (testcontainers.Container, error) {
	ctx := context.Background()
	containerName := image.GetContainerName()

	// Check if container already exists and is running
	if CheckDockerContainerExists(containerName) {
		log.WithFields(log.Fields{
			"container": containerName,
		}).Debug("Container already exists, attempting to reuse")
	}

	// Create mount path for data persistence
	volumePath := internal.LpnWorkspace + "/" + containerName
	
	var container testcontainers.Container
	var err error

	// Use appropriate module based on database type
	switch image.GetType() {
	case "mysql":
		container, err = mysql.Run(
			ctx,
			image.GetFullyQualifiedName(),
			mysql.WithDatabase(DBName),
			mysql.WithUsername(DBUser),
			mysql.WithPassword(DBPassword),
			// Reuse existing container if it exists
			testcontainers.WithReuseByName(containerName),
			// Mount volume for data persistence
			testcontainers.WithMounts(
				testcontainers.BindMount(volumePath, testcontainers.ContainerMountTarget(image.GetDataFolder())),
			),
			// Add labels for identification
			testcontainers.WithLabels(map[string]string{
				"db-type":  image.GetType(),
				"lpn-type": image.GetLpnType(),
			}),
			// Wait for database to be ready
			testcontainers.WithWaitStrategy(
				wait.ForLog("port: 3306  MySQL Community Server"),
			),
		)

	case "postgresql":
		container, err = postgres.Run(
			ctx,
			image.GetFullyQualifiedName(),
			postgres.WithDatabase(DBName),
			postgres.WithUsername(DBUser),
			postgres.WithPassword(DBPassword),
			// Reuse existing container if it exists
			testcontainers.WithReuseByName(containerName),
			// Mount volume for data persistence
			testcontainers.WithMounts(
				testcontainers.BindMount(volumePath, testcontainers.ContainerMountTarget(image.GetDataFolder())),
			),
			// Add labels for identification
			testcontainers.WithLabels(map[string]string{
				"db-type":  image.GetType(),
				"lpn-type": image.GetLpnType(),
			}),
		)

	default:
		return nil, fmt.Errorf("unsupported database type: %s", image.GetType())
	}

	if err != nil {
		log.WithFields(log.Fields{
			"container": containerName,
			"image":     image.GetFullyQualifiedName(),
			"error":     err,
		}).Error("Failed to run database container with testcontainers")
		return nil, err
	}

	log.WithFields(log.Fields{
		"container": containerName,
		"image":     image.GetFullyQualifiedName(),
	}).Info("Database container started with testcontainers (reusable)")

	return container, nil
}

// StopDatabaseContainerWithTestcontainers stops a testcontainers-managed container
// without removing it, allowing it to be reused later.
func StopDatabaseContainerWithTestcontainers(container testcontainers.Container) error {
	ctx := context.Background()
	
	// Stop the container but don't terminate (remove) it
	if err := container.Stop(ctx, nil); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to stop container")
		return err
	}

	log.Info("Container stopped (will be reused on next start)")
	return nil
}

// RemoveDatabaseContainerWithTestcontainers permanently removes a testcontainers-managed container.
func RemoveDatabaseContainerWithTestcontainers(container testcontainers.Container) error {
	ctx := context.Background()
	
	if err := container.Terminate(ctx); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to terminate container")
		return err
	}

	log.Info("Container terminated and removed")
	return nil
}
