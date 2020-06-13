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

	"github.com/stretchr/testify/assert"
)

func TestDeployFolderDXP(t *testing.T) {
	dxp := DXP{}

	assert := assert.New(t)

	assert.Equal(dxp.GetLiferayHome()+"/deploy", dxp.GetDeployFolder())
}

func TestGetContainerNameDXP(t *testing.T) {
	dxp := DXP{}

	assert := assert.New(t)

	assert.Equal("lpn-dxp", dxp.GetContainerName())
}

func TestGetDockerHubTagsURLDXP(t *testing.T) {
	dxp := DXP{}

	assert := assert.New(t)

	assert.Equal("liferay/dxp", dxp.GetDockerHubTagsURL())
}

func TestGetFullyQualifiedNameDXP(t *testing.T) {
	dxp := DXP{Tag: "foo"}

	assert := assert.New(t)

	assert.Equal("liferay/dxp:foo", dxp.GetFullyQualifiedName())
}

func TestGetLiferayHomeDXP(t *testing.T) {
	dxp := DXP{}

	assert := assert.New(t)

	assert.Equal("/opt/liferay", dxp.GetLiferayHome())
}

func TestGetDXPsRepository(t *testing.T) {
	dxp := DXP{}

	assert := assert.New(t)
	ces := dxp.GetRepository()

	assert.Equal("liferay/dxp", ces)
}

func TestGetTypeDXP(t *testing.T) {
	dxp := DXP{}

	assert := assert.New(t)

	assert.Equal("dxp", dxp.GetType())
}

func TestGetUserDXP(t *testing.T) {
	dxp := DXP{}

	assert := assert.New(t)

	assert.Equal("liferay", dxp.GetUser())
}
