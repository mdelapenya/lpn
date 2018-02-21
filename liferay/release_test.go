package liferay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLiferayHomeRelease7Ga5(t *testing.T) {
	testGetLiferayHomeRelease7Ga(t, "5")
}

func TestGetLiferayHomeRelease7Ga4(t *testing.T) {
	testGetLiferayHomeRelease7Ga(t, "4")
}

func TestGetLiferayHomeRelease7Ga3(t *testing.T) {
	testGetLiferayHomeRelease7Ga(t, "3")
}

func TestGetLiferayHomeRelease7Ga2(t *testing.T) {
	testGetLiferayHomeRelease7Ga(t, "2")
}

func TestGetLiferayHomeRelease7Ga1(t *testing.T) {
	testGetLiferayHomeRelease7Ga(t, "1")
}

func TestGetLiferayHomeRelease6_2Ga6(t *testing.T) {
	release := Release{}

	release.tag = "6.2-ce-ga6-tomcat-hsql"

	assert := assert.New(t)

	assert.Equal("/usr/local/liferay-portal-6.2-ce-ga1", release.GetLiferayHome())
}

func TestGetLiferayHomeRelease6_1Ga1(t *testing.T) {
	release := Release{}

	release.tag = "6.1-ce-ga1-tomcat-hsql"

	assert := assert.New(t)

	assert.Equal("/usr/local/liferay-portal-6.1.0-ce-ga1", release.GetLiferayHome())
}

func TestGetLiferayHomeReleaseNoTag(t *testing.T) {
	release := Release{}

	assert := assert.New(t)

	assert.Equal("/usr/local/liferay", release.GetLiferayHome())
}

func TestGetReleaseRepository(t *testing.T) {
	release := Release{}

	assert := assert.New(t)
	releases := release.GetRepository()

	assert.Equal("mdelapenya/liferay-portal", releases)
}

func testGetLiferayHomeRelease7Ga(t *testing.T, ga string) {
	release := Release{}

	release.tag = "7-ce-ga" + ga + "-tomcat-hsql"

	assert := assert.New(t)

	assert.Equal("/usr/local/liferay-ce-portal-7.0-ga"+ga, release.GetLiferayHome())
}
