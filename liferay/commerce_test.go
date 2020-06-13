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

func TestDeployFolderCommerce(t *testing.T) {
	commerce := Commerce{}

	assert := assert.New(t)

	assert.Equal("/opt/liferay/deploy", commerce.GetDeployFolder())
}

func TestGetContainerNameCommerce(t *testing.T) {
	commerce := Commerce{}

	assert := assert.New(t)

	assert.Equal("lpn-commerce", commerce.GetContainerName())
}

func TestGetDockerHubTagsURLCommerce(t *testing.T) {
	commerce := Commerce{}

	assert := assert.New(t)

	assert.Equal("liferay/commerce", commerce.GetDockerHubTagsURL())
}

func TestGetFullyQualifiedNameCommerce(t *testing.T) {
	commerce := Commerce{Tag: "foo"}

	assert := assert.New(t)

	assert.Equal("liferay/commerce:foo", commerce.GetFullyQualifiedName())
}

func TestGetLiferayHomeCommerce(t *testing.T) {
	commerce := Commerce{}

	assert := assert.New(t)

	assert.Equal("/opt/liferay", commerce.GetLiferayHome())
}

func TestGetCommerceRepository(t *testing.T) {
	commerce := Commerce{}

	assert := assert.New(t)
	commerceRepository := commerce.GetRepository()

	assert.Equal("liferay/commerce", commerceRepository)
}

func TestGetTypeCommerce(t *testing.T) {
	commerce := Commerce{}

	assert := assert.New(t)

	assert.Equal("commerce", commerce.GetType())
}

func TestGetUserCommerce(t *testing.T) {
	commerce := Commerce{}

	assert := assert.New(t)

	assert.Equal("liferay", commerce.GetUser())
}
