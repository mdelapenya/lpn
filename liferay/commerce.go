package liferay

// Commerce implementation for Liferay nightly images with Commerce
type Commerce struct {
	Tag string
}

// GetDeployFolder returns the deploy folder under Liferay Home
func (c Commerce) GetDeployFolder() string {
	return c.GetLiferayHome() + "/deploy"
}

// GetFullyQualifiedName returns the fully qualified name of the image
func (c Commerce) GetFullyQualifiedName() string {
	return c.GetRepository() + ":" + c.GetTag()
}

// GetLiferayHome returns the Liferay home for nightly builds with Commerce
func (c Commerce) GetLiferayHome() string {
	return "/liferay"
}

// GetRepository returns the repository for nightly builds with Commerce
func (c Commerce) GetRepository() string {
	return CommercesRepository
}

// GetTag returns the tag of the image
func (c Commerce) GetTag() string {
	return c.Tag
}

// GetType returns the type of the image
func (c Commerce) GetType() string {
	return "commerce"
}
