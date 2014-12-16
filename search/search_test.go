package search_test

import (
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"
)

func ExampleByCoordinates_issues() {
	// Nexus' coordinate search has several issues and peculiarities, which
	// affect the results this library gives. Some examples follow below.

	n := nexus.New("https://maven.java.net", credentials.None)

	// A coordinate search requires specifying at least either a groupId, an
	// artifactId or a version. This search will (after some time), return
	// nothing. This doesn't mean there are no projects with packaging "pom";
	// this is a limitation of Nexus' search.
	n.Artifacts(search.ByCoordinates{Packaging: "pom"})

	// This search may or may not return an error, depending on the version of
	// the Nexus being accessed. On newer Nexuses (sp?) "*" searches are
	// invalid.
	n.Artifacts(search.ByCoordinates{GroupId: "*", Packaging: "pom"})

	// ByCoordinates searches in Maven *projects*, not artifacts. So this
	// search will return all com.sun* artifacts in projects with packaging
	// "pom", not all POM artifacts with groupId com.sun*! Packaging is not
	// the same as extension.
	n.Artifacts(search.ByCoordinates{GroupId: "com.sun*", Packaging: "pom"})
}
