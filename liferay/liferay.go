package liferay

// CommerceRepository Namespace for the Docker releases with Commerce
const CommerceRepository = "liferay/liferay-commerce"

// Nightlies Namespace for the Docker nightly builds
const Nightlies = "mdelapenya/liferay-portal-nightlies"

// Releases Namespace for the Docker releases
const Releases = "mdelapenya/liferay-portal"

// Image interface defining the contract for Liferay Portal docker images
type Image interface {
	GetFullyQualifiedName() string
	GetDeployFolder() string
	GetLiferayHome() string
	GetRepository() string
	GetTag() string
}
