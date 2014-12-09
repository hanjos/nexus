// Package search provides a mini-DSL for nexus.Client.Artifacts().
package search

import (
	"fmt"
	"strings"
)

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

func (empty noCriteria) String() string {
	return "search.None"
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

func (gav ByCoordinates) String() string {
	str := []string{}

	if gav.GroupId != "" {
		str = append(str, "g: " + gav.GroupId)
	}
	if gav.ArtifactId != "" {
		str = append(str, "a: " + gav.ArtifactId)
	}
	if gav.Version != "" {
		str = append(str, "v: " + gav.Version)
	}
	if gav.Packaging != "" {
		str = append(str, "p: " + gav.Packaging)
	}
	if gav.Classifier != "" {
		str = append(str, "c: " + gav.Classifier)
	}

	return "search.ByCoordinates(" + strings.Join(str, ", ") + ")"
}

// Searches by keywords.
type ByKeyword string

func (q ByKeyword) Parameters() map[string]string {
	return map[string]string{
		"q": string(q),
	}
}

func (q ByKeyword) String() string {
	return "search.ByKeyword(" + string(q) + ")"
}

// Searches by class name.
type ByClassname string

func (cn ByClassname) Parameters() map[string]string {
	return map[string]string{
		"cn": string(cn),
	}
}

func (cn ByClassname) String() string {
	return "search.ByClassname(" + string(cn) + ")"
}

// Searches by SHA1 checksum.
type ByChecksum string

func (sha1 ByChecksum) Parameters() map[string]string {
	return map[string]string{
		"sha1": string(sha1),
	}
}

func (sha1 ByChecksum) String() string {
	return "search.ByChecksum(" + string(sha1) + ")"
}

// Searches for all artifacts in the given repository ID.
type ByRepository string

func (byRepo ByRepository) Parameters() map[string]string {
	return map[string]string{
		"repositoryId": string(byRepo),
	}
}

func (byRepo ByRepository) String() string {
	return "search.ByRepository(" + string(byRepo) + ")"
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

func (inRepo InRepository) String() string {
	return "search.InRepository(" + inRepo.RepositoryId + ", " + fmt.Sprintf("%v", inRepo.Criteria) + ")"
}
