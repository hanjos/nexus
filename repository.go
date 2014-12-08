package nexus

// A Nexus repository. Nexus actually provides a bit more data, but this should be enough for most uses.
type Repository struct {
	Id        string
	Name      string
	Type      string
	Format    string
	Policy    string
	RemoteURI string
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
