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
