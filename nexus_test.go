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

	artifacts, err := n.Artifacts(
		search.ByKeyword("javax.enterprise"))

	if err != nil {
		fmt.Printf("%v: %v", reflect.TypeOf(err), err)
		return
	}

	for _, a := range artifacts {
		fmt.Println(a)
	}
}

func ExampleNexus2x_Artifacts() {
	n := nexus.New("http://maven.java.net", credentials.None)

	// using a simple search
	artifacts, err := n.Artifacts(search.ByClassname("javax.servlet.Servlet"))

	// using a composite search
	artifacts, err := n.Artifacts(search.InRepository{"releases", search.ByKeyword("javax.enterprise")})
}

func ExampleNexus2x_Artifacts_full_search() {
	n := nexus.New("http://maven.java.net", credentials.None)

	// returns an error
	artifacts, err := n.Artifacts(search.None)

	// if you want all artifacts in this Nexus, you can search in each repository one by one. Generally you don't want
	// to do that, especially if you have proxy repositories; central has, at the time of this comment, over 800,000
	// artifacts (!), which in this implementation will be all loaded into memory (!!). But, if you insist, the easiest
	// way to do it is:

	// 1) get all repositories
	repositories, err := n.Repositories()
	if err != nil {
		// handle the error
	}

	// and 2) accumulate the results from every repository
	artifacts := []*nexus.Artifact{}
	for _, repo := range repositories {
		a, err := n.Artifacts(search.ByRepository(repo.Id))
		if err != nil {
			// handle the error
		}

		artifacts = append(artifacts, a)
	}

	// Concurrent searches are left as an exercise for the reader :)
}
