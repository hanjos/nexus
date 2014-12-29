package credentials_test

import (
	"testing"
	"github.com/hanjos/nexus/credentials"
)

func TestNoneSignDoesntBarfOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%v", r)
		}
	}()

	credentials.None.Sign(nil)
}

func TestBasicAuthSignDoesntBarfOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%v", r)
		}
	}()

	credentials.BasicAuth{"u", "p"}.Sign(nil)
}
