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

import liferay "github.com/mdelapenya/lpn/liferay"

// DBName name of the default database
const DBName = "lportal"

// DBPassword default credentials for the database
const DBPassword = "my-secret-pw"

// DBUser default user for the database
const DBUser = "liferay"

// DatabaseImage interface defining the contract for database docker images
type DatabaseImage interface {
	GetContainerName() string
	GetDataFolder() string
	GetEnvVariables() EnvVariables
	GetJDBCConnection() JDBCConnection
	GetFullyQualifiedName() string
	GetLpnType() string
	GetPort() int
	GetRepository() string
	GetTag() string
	GetType() string
}

// GetAlias returns the alias used to link containers
func GetAlias() string {
	return "db"
}

// GetDatabase returns the proper database model
func GetDatabase(image liferay.Image, datastore string) DatabaseImage {
	if datastore == "mysql" {
		return MySQL{LpnType: image.GetType()}
	} else if datastore == "postgresql" {
		return PostgreSQL{LpnType: image.GetType()}
	}

	return nil
}

// EnvVariables defines how to configure the internal variables for the database
type EnvVariables struct {
	Password string
	Database string
	User     string
}

// JDBCConnection defines the JDBC connection to the database
type JDBCConnection struct {
	DriverClassName string
	Password        string
	User            string
	URL             string
}
