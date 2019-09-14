package liferay

import (
	"testing"

	internal "github.com/mdelapenya/lpn/internal"
	"github.com/stretchr/testify/assert"
)

func init() {
	internal.CheckWorkspace()
}

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

	assert.Equal("docker.io/liferay/portal-snapshot:foo", nightly.GetFullyQualifiedName())
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
