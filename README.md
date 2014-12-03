go-nexus
========

A Go library to pull some data from a [Sonatype Nexus](http://www.sonatype.com/nexus) instance.

How?
====

`go get` should see you through:

```sh
go get github.com/hanjos/nexus
```

And therefore `import`:

```Go
package main

import (
  "fmt"
  "github.com/hanjos/nexus"
  "reflect"
)

func main() {
  n := nexus.New("http://nexus.somewhere.com")
  
  artifacts, err := n.Artifacts(nexus.InRepository{ nexus.ByKeyword("com.sbrubbles*"), "shamalamadingdong" })
  if err != nil {
    fmt.Printf("%v: %v", reflect.TypeOf(err), err)
  }

  for _, a := range artifacts {
    fmt.Println(a)
  }
}
```

Why?
====

Nexus has a large REST API, but some information isn't readily available, requiring several API calls and some mashing 
up to produce. 

And it was a good excuse to try Go out :)

LICENSE
=======

MIT License. See [LICENSE](https://github.com/hanjos/nexus/blob/master/LICENSE) for the gory details.
