package liferay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReleaseRepository(t *testing.T) {
	release := Release{}

	assert := assert.New(t)
	releases := release.GetRepository()

	assert.Equal("mdelapenya/liferay-portal", releases)
}
