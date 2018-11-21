package liferay

// CERepository Namespace for the official Docker releases for CE
const CERepository = "liferay/portal"

// CommercesRepository Namespace for the Docker releases with Commerce
const CommercesRepository = "liferay/liferay-commerce"

// DXPRepository Namespace for the official Docker releases for DXP
const DXPRepository = "liferay/dxp"

// NightliesRepository Namespace for the Docker nightly builds
const NightliesRepository = "mdelapenya/liferay-portal-nightlies"

// ReleasesRepository Namespace for the Docker releases
const ReleasesRepository = "mdelapenya/liferay-portal"

// Image interface defining the contract for Liferay Portal docker images
type Image interface {
	GetContainerName() string
	GetFullyQualifiedName() string
	GetDeployFolder() string
	GetDockerHubTagsURL() string
	GetLiferayHome() string
	GetRepository() string
	GetTag() string
	GetType() string
}
