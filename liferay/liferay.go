package liferay

// CommercesRepository Namespace for the Docker releases with Commerce
const CommercesRepository = "liferay/liferay-commerce"

// NightliesRepository Namespace for the Docker nightly builds
const NightliesRepository = "mdelapenya/liferay-portal-nightlies"

// ReleasesRepository Namespace for the Docker releases
const ReleasesRepository = "mdelapenya/liferay-portal"

// Image interface defining the contract for Liferay Portal docker images
type Image interface {
	GetContainerName() string
	GetFullyQualifiedName() string
	GetDeployFolder() string
	GetLiferayHome() string
	GetRepository() string
	GetTag() string
	GetType() string
}
