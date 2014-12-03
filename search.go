package nexus

/*
 Nexus' API supports 4 different types of searches, but in the end, all we need is a map holding the parameters to pass
 along. Criteria enables a mini-DSL for nexus.Client.Artifacts(). The types and functions here cover the basics, but
 anything satisfying Criteria can be used.
*/
type Criteria interface {
	Parameters() map[string]string
}

// The zero value for criteria. There's no reason for more than one value to exist, so it's unexported and
// made bool for Go to allow a const.
type zeroCriteria bool

func (empty zeroCriteria) Parameters() map[string]string {
	return map[string]string{}
}

// The zero value for criteria. Its Parameters() method returns an empty map.
const CriteriaZero = zeroCriteria(false)

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

// Searches all artifacts in the given repository ID.
type ByRepository string

func (byRepo ByRepository) Parameters() map[string]string {
	return map[string]string{
		"repositoryId": string(byRepo),
	}
}

// Searches for all artifacts following the given criteria in the given repository ID.
type InRepository struct {
	Criteria

	RepositoryId string
}

func (inRepo InRepository) Parameters() map[string]string {
	params := inRepo.Criteria.Parameters()
	params["repositoryId"] = inRepo.RepositoryId

	return params
}
