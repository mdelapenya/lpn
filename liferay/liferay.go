package liferay

// nightlies Namespace for the Docker nightly builds
const nightlies = "mdelapenya/liferay-portal-nightlies"

// releases Namespace for the Docker releases
const releases = "mdelapenya/liferay-portal"

var repositories []string

// Image interface defining the contract for Liferay Portal docker images
type Image interface {
	GetRepository() string
}

func init() {
	repositories = append(repositories, releases, nightlies)
}

// GetNightlyBuildsRepository Return the nightly builds repository
func GetNightlyBuildsRepository() string {
	return repositories[1]
}

// GetReleasesRepository Return the releases repository
func GetReleasesRepository() string {
	return repositories[0]
}

// getRepositories Return an array with the available repositories
func getRepositories() []string {
	return repositories
}
