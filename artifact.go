package nexus

import (
	"fmt"
	"github.com/hanjos/nexus/util"
	"strings"
	"time"
)

// Artifact is a Maven coordinate to a single artifact.
type Artifact struct {
	GroupId      string // e.g. org.springframework
	ArtifactId   string // e.g. spring-core
	Version      string // e.g. 4.1.3.RELEASE
	Classifier   string // e.g. sources, javadoc, <the empty string>...
	Extension    string // e.g. jar
	RepositoryId string // e.g. releases
}

// String implements the Stringer interface, for easy printing, as per Maven docs
// (http://maven.apache.org/pom.html#Maven_Coordinates).
func (a Artifact) String() string {
	var parts = []string{a.GroupId, a.ArtifactId, a.Extension}

	if a.Classifier != "" {
		parts = append(parts, a.Classifier)
	}

	return strings.Join(append(parts, a.Version), ":") + "@" + a.RepositoryId
}

// DefaultFileName builds the default name Maven gives an artifact given its metadata.
func (a Artifact) DefaultFileName() string {
	classifier := ""
	if a.Classifier != "" {
		classifier = "-" + a.Classifier
	}

	return a.ArtifactId + "-" + a.Version + classifier + "." + a.Extension
}

// used for the artifact set.
func (a *Artifact) hash() string {
	return a.GroupId + ":" + a.ArtifactId + ":" + a.Version + ":" +
		a.Extension + ":" + a.Classifier + "@" + a.RepositoryId
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

// ArtifactInfo holds extra information about the given artifact.
type ArtifactInfo struct {
	Artifact

	Uploader    string
	Uploaded    time.Time
	LastChanged time.Time
	Sha1        string
	Size        util.FileSize
	MimeType    string
	Url         string
}

// String implements the Stringer interface, for easy printing.
func (info ArtifactInfo) String() string {
	return fmt.Sprintf("%v [SHA1 %v, Mime-Type %v, %v]", info.Artifact, info.Sha1, info.MimeType, info.Size)
}
