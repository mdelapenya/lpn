package liferay

// Nightly implementation for Liferay nightly images
type Nightly struct {
	Tag string
}

// GetDeployFolder returns the deploy folder under Liferay Home
func (n Nightly) GetDeployFolder() string {
	return n.GetLiferayHome() + "/deploy"
}

// GetFullyQualifiedName returns the fully qualified name of the image
func (n Nightly) GetFullyQualifiedName() string {
	return n.GetRepository() + ":" + n.GetTag()
}

// GetLiferayHome returns the Liferay home for nightly builds
func (n Nightly) GetLiferayHome() string {
	return "/liferay"
}

// GetRepository returns the repository for nightly builds
func (n Nightly) GetRepository() string {
	return NightliesRepository
}

// GetTag returns the tag of the image
func (n Nightly) GetTag() string {
	return n.Tag
}

// GetType returns the type of the image
func (n Nightly) GetType() string {
	return "nightly"
}
