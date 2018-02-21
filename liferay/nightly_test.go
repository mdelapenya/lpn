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

func TestGetNightlyRepository(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)
	nightlies := nightly.GetRepository()

	assert.Equal("mdelapenya/liferay-portal-nightlies", nightlies)
}
