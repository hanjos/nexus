package nexus_test

import (
	"fmt"
	"reflect"
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"
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
