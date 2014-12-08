package nexus

import "strings"

// Artifact is a Maven coordinate to a single artifact.
type Artifact struct {
	GroupId    string
	ArtifactId string
	Version    string
	Classifier string
	Extension  string
}

// String returns the string representation of an artifact, as per Maven docs
// (http://maven.apache.org/pom.html#Maven_Coordinates).
func (a *Artifact) String() string {
	var parts = []string{a.GroupId, a.ArtifactId, a.Extension}

	if a.Classifier != "" {
		parts = append(parts, a.Classifier)
	}

	return strings.Join(append(parts, a.Version), ":")
}

// used for the artifact set.
func (a *Artifact) hash() string {
	return a.GroupId + ":" + a.ArtifactId + ":" + a.Version + ":" + a.Extension + ":" + a.Classifier
}

// since Go doesn't have a built-in set implementation, a make-shift one follows, using a map for the heavy duty.
// Artifact's hash method is used to distinguish between artifacts; there's no Java-like Equals contract to follow.
type artifactSet struct {
	// piles up the artifacts
	data []*Artifact

	// the set behavior
	hashMap map[string]bool
}

// creates and initializes a new set of artifacts.
func newArtifactSet() *artifactSet {
	return &artifactSet{
		data:    []*Artifact{},
		hashMap: make(map[string]bool),
	}
}

// adds a bunch of artifacts to this set.
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
