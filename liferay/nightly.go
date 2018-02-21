package liferay

// Nightly implementation for Liferay nightly images
type Nightly struct {
}

// GetLiferayHome returns the Liferay home for nightly builds
func (n Nightly) GetLiferayHome() string {
	return "/liferay"
}

// GetRepository returns the repository for nightly builds
func (n Nightly) GetRepository() string {
	return Nightlies
}
