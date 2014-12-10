package nexus_test

import (
	"fmt"
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"
	"reflect"
)

func Example() {
	n := nexus.New("https://maven.java.net", credentials.None)

	// obtaining all repositories in Nexus
	repositories, err := n.Repositories()
	if err != nil {
		fmt.Printf("%v: %v", reflect.TypeOf(err), err)
		return
	}

	// printing out all artifacts which are in a hosted repository, and have
	// in their groupId 'javax.enterprise' and a 'pom' packaging.
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
	artifacts, err := n.Artifacts(search.ByClassname("javax.servlet.Servlet"))

	// using a composite search
	artifacts, err := n.Artifacts(
		search.InRepository{
			"releases",
			search.ByKeyword("javax.enterprise")})
}
