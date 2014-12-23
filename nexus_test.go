package nexus_test

import (
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"

	"fmt"
	"reflect"
	"testing"
)

func TestNexus2xImplementsClient(t *testing.T) {
	if _, ok := interface{}(nexus.Nexus2x{}).(nexus.Client); !ok {
		t.Errorf("nexus.Nexus2x does not implement nexus.Client!")
	}
}

func Example() {
	n := nexus.New("https://maven.java.net", credentials.None)

	// obtaining all repositories in Nexus
	repositories, err := n.Repositories()
	if err != nil {
		fmt.Printf("%v: %v", reflect.TypeOf(err), err)
		return
	}

	// printing out all artifacts which are in a hosted repository, and have
	// both 'javax.enterprise' in their groupId and a 'sources' classifier.
	for _, repo := range repositories {
		if repo.Type != "hosted" {
			continue
		}

		artifacts, err := n.Artifacts(
			search.InRepository{
				repo.Id,
				search.ByCoordinates{
					GroupId:    "javax.enterprise*",
					Classifier: "sources"}})

		if err != nil {
			fmt.Printf("%v: %v", reflect.TypeOf(err), err)
			return
		}

		for _, a := range artifacts {
			fmt.Println(a)
		}
	}
}

func ExampleNexus2x_Artifacts() {
	n := nexus.New("http://maven.java.net", credentials.None)

	// using a simple search
	n.Artifacts(search.ByClassname("javax.servlet.Servlet"))

	// using a composite search
	n.Artifacts(
		search.InRepository{
			"releases",
			search.ByKeyword("javax.enterprise")})

	// searching for every artifact in Nexus (WARNING: this can take a LOOONG time - and memory!)
	n.Artifacts(search.All)
}
