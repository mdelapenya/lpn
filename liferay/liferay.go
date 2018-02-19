package liferay

// Nightlies Namespace for the Docker nightly builds
const Nightlies = "mdelapenya/liferay-portal-nightlies"

// Releases Namespace for the Docker releases
const Releases = "mdelapenya/liferay-portal"

var repositories []string

func init() {
	repositories = append(repositories, Releases, Nightlies)
}

// GetNightlyBuildsRepository Return the nightly builds repository
func GetNightlyBuildsRepository() string {
	return repositories[1]
}

// GetReleasesRepository Return the releases repository
func GetReleasesRepository() string {
	return repositories[0]
}

// GetRepositories Return an array with the available repositories
func GetRepositories() []string {
	return repositories
}
