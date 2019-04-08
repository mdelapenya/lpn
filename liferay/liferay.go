package liferay

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
