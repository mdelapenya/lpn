package liferay

// CEDefaultTag Default tag for CE
const CEDefaultTag = "7.0.6-ga7"

// CERepository Namespace for the official Docker releases for CE
const CERepository = "liferay/portal"

// CommerceDefaultTag Default tag for CE
const CommerceDefaultTag = "1.1.1"

// CommerceRepository Namespace for the Docker releases with Commerce
const CommerceRepository = "liferay/commerce"

// DXPDefaultTag Default tag for DXP
const DXPDefaultTag = "7.0.10.8"

// DXPRepository Namespace for the official Docker releases for DXP
const DXPRepository = "liferay/dxp"

// NightliesRepository Namespace for the Docker nightly builds
const NightliesRepository = "liferay/portal-snapshot"

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
	GetUser() string
}
