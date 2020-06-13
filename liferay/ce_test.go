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

package liferay

import (
	"testing"

	internal "github.com/mdelapenya/lpn/internal"
	"github.com/stretchr/testify/assert"
)

func init() {
	internal.CheckWorkspace()
}

func TestDeployFolderCE(t *testing.T) {
	ce := CE{}

	assert := assert.New(t)

	assert.Equal(ce.GetLiferayHome()+"/deploy", ce.GetDeployFolder())
}

func TestGetContainerNameCE(t *testing.T) {
	ce := CE{}

	assert := assert.New(t)

	assert.Equal("lpn-ce", ce.GetContainerName())
}

func TestGetDockerHubTagsURLCE(t *testing.T) {
	ce := CE{}

	assert := assert.New(t)

	assert.Equal("liferay/portal", ce.GetDockerHubTagsURL())
}

func TestGetFullyQualifiedNameCE(t *testing.T) {
	ce := CE{Tag: "foo"}

	assert := assert.New(t)

	assert.Equal("docker.io/liferay/portal:foo", ce.GetFullyQualifiedName())
}

func TestGetLiferayHomeCE(t *testing.T) {
	ce := CE{}

	assert := assert.New(t)

	assert.Equal("/opt/liferay", ce.GetLiferayHome())
}

func TestGetCEsRepository(t *testing.T) {
	ce := CE{}

	assert := assert.New(t)
	ces := ce.GetRepository()

	assert.Equal("liferay/portal", ces)
}

func TestGetTypeCE(t *testing.T) {
	ce := CE{}

	assert := assert.New(t)

	assert.Equal("ce", ce.GetType())
}

func TestGetUserCE(t *testing.T) {
	ce := CE{}

	assert := assert.New(t)

	assert.Equal("liferay", ce.GetUser())
}
