package search_test

import (
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"

	"testing"
)

func checkMap(t *testing.T, expected map[string]string, actual map[string]string) {
	if len(expected) != len(actual) {
		t.Errorf("Wrong number of fields: expected %v, got %v", len(expected), len(actual))
	}

	for k, v := range actual {
		vExp, ok := expected[k]

		if !ok {
			t.Errorf("Unexpected field %q", k)
		} else if vExp != v {
			t.Errorf("Expected value %q for field %q, got %q", vExp, k, v)
		}
	}
}

func TestAllImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.All).(search.Criteria); !ok {
		t.Errorf("search.All does not implement Criteria!")
	}
}

func TestAllProvidesNoCriteria(t *testing.T) {
	criteria := search.All.Parameters()
	if len(criteria) != 0 {
		t.Errorf("expected an empty map, got %v", criteria)
	}
}

func TestByCoordinatesImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByCoordinates{}).(search.Criteria); !ok {
		t.Errorf("search.ByCoordinates does not implement Criteria!")
	}
}

type pair struct {
	expected string
	actual   string
}

func TestByCoordinatesSetsTheProperFields(t *testing.T) {
	criteria := search.ByCoordinates{GroupId: "g", ArtifactId: "a", Version: "v", Packaging: "p", Classifier: "c"}.Parameters()

	expected := []string{"g", "a", "v", "p", "c"}
	missing := []string{}
	wrong := []pair{}

	for _, exp := range expected {
		v, ok := criteria[exp]

		if !ok {
			missing = append(missing, exp)
		}
		if v != exp {
			wrong = append(wrong, pair{exp, v})
		}
	}

	if len(missing) != 0 {
		t.Errorf("Missing fields %v", missing)
	}

	if len(wrong) != 0 {
		t.Errorf("Fields with wrong values:\n")
		for _, p := range wrong {
			t.Errorf("Field %q expected value %q, got %q", p.expected, p.expected, p.actual)
		}
	}
}

func TestByCoordinatesSetsOnlyTheGivenFields(t *testing.T) {
	criteria := search.ByCoordinates{GroupId: "g", ArtifactId: "a"}.Parameters()

	checkMap(t, map[string]string{"g": "g", "a": "a"}, criteria)
}

func TestByClassnameImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByClassname("")).(search.Criteria); !ok {
		t.Errorf("search.ByClassname does not implement Criteria!")
	}
}

func TestByClassnameSetsTheProperFields(t *testing.T) {
	criteria := search.ByClassname("cn").Parameters()

	checkMap(t, map[string]string{"cn": "cn"}, criteria)
}

func TestByChecksumImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByChecksum("")).(search.Criteria); !ok {
		t.Errorf("search.ByChecksum does not implement Criteria!")
	}
}

func TestByChecksumSetsTheProperFields(t *testing.T) {
	criteria := search.ByChecksum("sha1").Parameters()

	checkMap(t, map[string]string{"sha1": "sha1"}, criteria)
}

func TestByKeywordImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByKeyword("")).(search.Criteria); !ok {
		t.Errorf("search.ByKeyword does not implement Criteria!")
	}
}

func TestByKeywordSetsTheProperFields(t *testing.T) {
	criteria := search.ByKeyword("q").Parameters()

	checkMap(t, map[string]string{"q": "q"}, criteria)
}

func TestByRepositoryImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.ByRepository("")).(search.Criteria); !ok {
		t.Errorf("search.ByRepository does not implement Criteria!")
	}
}

func TestByRepositorySetsTheProperFields(t *testing.T) {
	criteria := search.ByRepository("repositoryId").Parameters()

	checkMap(t, map[string]string{"repositoryId": "repositoryId"}, criteria)
}

func TestInRepositoryImplementsCriteria(t *testing.T) {
	if _, ok := interface{}(search.InRepository{}).(search.Criteria); !ok {
		t.Errorf("search.InRepository does not implement Criteria!")
	}
}

func TestInRepositorySetsTheProperFields(t *testing.T) {
	criteria := search.InRepository{"repositoryId", search.ByChecksum("sha1")}.Parameters()

	checkMap(t, map[string]string{"repositoryId": "repositoryId", "sha1": "sha1"}, criteria)
}

func TestInRepositoryWithSearchAllIsTheSameAsByRepository(t *testing.T) {
	inRepo := search.InRepository{"repositoryId", search.All}.Parameters()
	byRepo := search.ByRepository("repositoryId").Parameters()

	checkMap(t, byRepo, inRepo)
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
