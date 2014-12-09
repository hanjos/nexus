package search_test

import (
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"
)

func Example() {
	n := nexus.New("http://maven.java.net", credentials.None)

	// using a simple search
	artifacts, err := n.Artifacts(search.ByClassname("javax.servlet.Servlet"))

	// using a composite search
	artifacts, err := n.Artifacts(search.InRepository{"releases", search.ByKeyword("javax.enterprise")})
}
