package liferay

// Release implementation for Liferay released images
type Release struct {
}

// GetRepository returns the repository for releases
func (r Release) GetRepository() string {
	return GetReleasesRepository()
}
