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

func TestDeployFolderNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("/opt/liferay/deploy", nightly.GetDeployFolder())
}

func TestGetContainerNameNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("lpn-nightly", nightly.GetContainerName())
}

func TestGetDockerHubTagsURLNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("liferay/portal-snapshot", nightly.GetDockerHubTagsURL())
}

func TestGetFullyQualifiedNameNightly(t *testing.T) {
	nightly := Nightly{Tag: "foo"}

	assert := assert.New(t)

	assert.Equal("liferay/portal-snapshot:foo", nightly.GetFullyQualifiedName())
}

func TestGetLiferayHomeNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("/opt/liferay", nightly.GetLiferayHome())
}

func TestGetNightliesRepository(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)
	nightlies := nightly.GetRepository()

	assert.Equal("liferay/portal-snapshot", nightlies)
}

func TestGetTypeNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("nightly", nightly.GetType())
}

func TestGetUserNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("liferay", nightly.GetUser())
}
