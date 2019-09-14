package liferay

import (
	"testing"

	internal "github.com/mdelapenya/lpn/internal"
	"github.com/stretchr/testify/assert"
)

func init() {
	internal.CheckWorkspace()
}

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

	assert.Equal("docker.io/liferay/dxp:foo", dxp.GetFullyQualifiedName())
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
