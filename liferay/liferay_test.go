package liferay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNightlyRepository(t *testing.T) {
	nightly := Nightly{}

	assert := assert.New(t)
	nightlies := nightly.GetRepository()

	assert.Equal("mdelapenya/liferay-portal-nightlies", nightlies)
}

func TestGetRepositories(t *testing.T) {
	assert := assert.New(t)
	repositories := getRepositories()

	assert.Equal(2, len(repositories), "There must be only two repositories")
}

func TestGetStableRepository(t *testing.T) {
	release := Release{}

	assert := assert.New(t)
	releases := release.GetRepository()

	assert.Equal("mdelapenya/liferay-portal", releases)
}
