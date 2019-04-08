package liferay

// CEDefaultTag Default tag for CE
const CEDefaultTag = "7.0.6-ga7"

// CommerceDefaultTag Default tag for CE
const CommerceDefaultTag = "1.1.1"

// DXPDefaultTag Default tag for DXP
const DXPDefaultTag = "7.0.10.8"

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
