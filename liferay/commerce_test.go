package liferay

import (
	"testing"

	internal "github.com/mdelapenya/lpn/internal"
	"github.com/stretchr/testify/assert"
)

func init() {
	internal.CheckWorkspace()
}

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

	assert.Equal("docker.io/liferay/commerce:foo", commerce.GetFullyQualifiedName())
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
