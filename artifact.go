package nexus

import "strings"

// A full Maven coordinate to a single artifact.
type Artifact struct {
	GroupId    string
	ArtifactId string
	Version    string
	Classifier string
	Extension  string
}

func (a *Artifact) String() string {
	var parts = []string{a.GroupId, a.ArtifactId, a.Version, a.Extension}

	if a.Classifier != "" {
		parts = append(parts, a.Classifier)
	}

	return strings.Join(parts, ":")
}

func (a *Artifact) IsPom() bool {
	return a.Classifier == "" && a.Extension == "pom"
}

// this is for the artifact set
func (a *Artifact) hash() string {
	return a.GroupId + ":" + a.ArtifactId + ":" + a.Version + ":" + a.Extension + ":" + a.Classifier
}

// Since Go doesn't have a built-in set implementation, a make-shift one follows, using a map for the heavy duty.
// Artifact's hash method is used to distinguish between artifacts; there's no Java-like Equals contract to follow.
type artifactSet struct {
	data    []*Artifact
	hashMap map[string]bool
}

func newArtifactSet() *artifactSet {
	return &artifactSet{
		data:    []*Artifact{},
		hashMap: make(map[string]bool),
	}
}

func (set *artifactSet) add(artifacts []*Artifact) {
	for _, artifact := range artifacts {
		hash := artifact.hash()
		_, contains := set.hashMap[hash]

		set.hashMap[hash] = true
		if !contains {
			set.data = append(set.data, artifact)
		}
	}
}
