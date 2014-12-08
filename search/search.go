// Package search provides a mini-DSL for nexus.Client.Artifacts().
package search

// Criteria compiles to a single map with the parameters Nexus expects. Nexus' API supports 4 different types of
// searches, but in the end, all we need is a map holding the parameters to pass along.
type Criteria interface {
	Parameters() map[string]string
}

// None is the zero value for Criteria. It returns an empty map.
const None = noCriteria(false)

// there's no reason for more than one value to exist, so it's unexported and
// made bool for Go to allow a const.
type noCriteria bool

func (empty noCriteria) Parameters() map[string]string {
	return map[string]string{}
}

// OrZero returns the given criteria untouched if it's not nil, and search.None otherwise. Useful for when one must
// ensure that the given criteria is non-nil.
func OrZero(c Criteria) Criteria {
	if c == nil {
		return None
	}

	return c
}

// Searches by Maven coordinates.
type ByCoordinates struct {
	GroupId    string
	ArtifactId string
	Version    string
	Packaging  string
	Classifier string
}

func (gav ByCoordinates) Parameters() map[string]string {
	return map[string]string{
		"g": gav.GroupId,
		"a": gav.ArtifactId,
		"v": gav.Version,
		"p": gav.Packaging,
		"c": gav.Classifier,
	}
}

// Searches by keywords.
type ByKeyword string

func (q ByKeyword) Parameters() map[string]string {
	return map[string]string{
		"q": string(q),
	}
}

// Searches by class name.
type ByClassname string

func (cn ByClassname) Parameters() map[string]string {
	return map[string]string{
		"cn": string(cn),
	}
}

// Searches by SHA1 checksum.
type ByChecksum string

func (sha1 ByChecksum) Parameters() map[string]string {
	return map[string]string{
		"sha1": string(sha1),
	}
}

// Searches for all artifacts in the given repository ID.
type ByRepository string

func (byRepo ByRepository) Parameters() map[string]string {
	return map[string]string{
		"repositoryId": string(byRepo),
	}
}

// Searches for all artifacts in the given repository ID following the given criteria.
type InRepository struct {
	RepositoryId string

	Criteria
}

func (inRepo InRepository) Parameters() map[string]string {
	params := inRepo.Criteria.Parameters()
	params["repositoryId"] = inRepo.RepositoryId

	return params
}
