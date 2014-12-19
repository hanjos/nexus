package search_test

import (
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"

	"testing"
)

func TestAllImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.All).(search.Criteria); !ok {
		t.Errorf("search.All does not implement Criteria!")
	}
}

func TestByCoordinatesImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByCoordinates{}).(search.Criteria); !ok {
		t.Errorf("search.ByCoordinates does not implement Criteria!")
	}
}

func TestByClassnameImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByClassname("")).(search.Criteria); !ok {
		t.Errorf("search.ByClassname does not implement Criteria!")
	}
}

func TestByChecksumImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByChecksum("")).(search.Criteria); !ok {
		t.Errorf("search.ByChecksum does not implement Criteria!")
	}
}

func TestByKeywordImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByKeyword("")).(search.Criteria); !ok {
		t.Errorf("search.ByKeyword does not implement Criteria!")
	}
}

func TestByRepositoryImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByRepository("")).(search.Criteria); !ok {
		t.Errorf("search.ByRepository does not implement Criteria!")
	}
}

func TestInRepositoryImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.InRepository{}).(search.Criteria); !ok {
		t.Errorf("search.InRepository does not implement Criteria!")
	}
}

// Examples

func ExampleByKeyword() {
	n := nexus.New("https://maven.java.net", credentials.None)

	// Return all artifacts with javax.enterprise somewhere.
	n.Artifacts(search.ByKeyword("javax.enterprise*"))

	// This search may or may not return an error, depending on the version of
	// the Nexus being accessed. On newer Nexuses (sp?) "*" searches are
	// invalid.
	n.Artifacts(search.ByKeyword("*"))
}

func ExampleByCoordinates() {
	n := nexus.New("https://maven.java.net", credentials.None)

	// Returns all artifacts with a groupId starting with com.sun. Due to Go's
	// struct syntax, we don't need to specify all the coordinates; they
	// default to string's zero value (""), which Nexus ignores.
	n.Artifacts(search.ByCoordinates{GroupId: "com.sun*"})

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
	n.Artifacts(search.ByCoordinates{GroupId: "com*", Packaging: "pom"})
}

func ExampleInRepository() {
	n := nexus.New("https://maven.java.net", credentials.None)

	// Returns all artifacts in the repository releases with groupId starting
	// with com.sun and whose project has packaging "pom".
	n.Artifacts(
		search.InRepository{
			"releases",
			search.ByCoordinates{GroupId: "com.sun*", Packaging: "pom"},
		})

	// Nexus doesn't support * in the repository ID parameter, so this search
	// will return an error.
	n.Artifacts(
		search.InRepository{
			"releases*",
			search.ByCoordinates{GroupId: "com.sun*", Packaging: "pom"},
		})
}
