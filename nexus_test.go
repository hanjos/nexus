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
	artifacts, err := n.Artifacts(
		search.InRepository{
			"releases",
			search.ByKeyword("javax.enterprise")})
}
