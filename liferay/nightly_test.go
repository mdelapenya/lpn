package liferay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployFolderNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("/liferay/deploy", nightly.GetDeployFolder())
}

func TestGetContainerNameNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("lpn-nightly", nightly.GetContainerName())
}

func TestGetDockerHubTagsURLNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("liferay/liferay-portal-nightlies", nightly.GetDockerHubTagsURL())
}

func TestGetFullyQualifiedNameNightly(t *testing.T) {
	nightly := Nightly{Tag: "foo"}

	assert := assert.New(t)

	assert.Equal("mdelapenya/liferay-portal-nightlies:foo", nightly.GetFullyQualifiedName())
}

func TestGetLiferayHomeNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("/liferay", nightly.GetLiferayHome())
}

func TestGetNightliesRepository(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)
	nightlies := nightly.GetRepository()

	assert.Equal("mdelapenya/liferay-portal-nightlies", nightlies)
}

func TestGetTypeNightly(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)

	assert.Equal("nightly", nightly.GetType())
}
