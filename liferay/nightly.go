package liferay

// Nightly implementation for Liferay nightly images
type Nightly struct {
}

// GetRepository returns the repository for nightly builds
func (n Nightly) GetRepository() string {
	return Nightlies
}
