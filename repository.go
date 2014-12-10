package nexus

// A Nexus repository. Nexus actually provides a bit more data, but this should be enough for most uses. Groups aren't
// considered repositories by Nexus' API; there's a separate call for them.
type Repository struct {
	Id        string // e.g. releases
	Name      string // e.g. Releases
	Type      string // e.g. hosted, proxy, virtual...
	Format    string // e.g. maven2, maven1...
	Policy    string // e.g. RELEASE, SNAPSHOT
	RemoteURI string // e.g. http://repo1.maven.org/maven2/
}

// String returns a pleasant but informative string representation of repo.
func (repo Repository) String() string {
	var uri string
	if repo.RemoteURI != "" {
		uri = ", points to " + repo.RemoteURI
	} else {
		uri = ""
	}

	return repo.Id + " ('" + repo.Name + "'){ " +
		repo.Type + ", " + repo.Format + " format, " +
		repo.Policy + " policy" + uri + " }"
}
