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
