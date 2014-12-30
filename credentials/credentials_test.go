package credentials_test

import (
	"github.com/hanjos/nexus/credentials"
	"testing"
)

func TestNoneImplementsCredentials(t *testing.T) {
	if _, ok := interface{}(credentials.None).(credentials.Credentials); !ok {
		t.Errorf("credentials.None doesn't implement credentials.Credentials!")
	}
}

func TestBasicAuthImplementsCredentials(t *testing.T) {
	if _, ok := interface{}(credentials.BasicAuth{"", ""}).(credentials.Credentials); !ok {
		t.Errorf("credentials.None doesn't implement credentials.Credentials!")
	}
}

func TestOrZeroReturnsTheGivenNonNilArgument(t *testing.T) {
	c := credentials.BasicAuth{"", ""}
	if v := credentials.OrZero(c); v != c {
		t.Errorf("credentials.OrZero(%v) should've returned %v, not %v!", c, c, v)
	}
}

func TestOrZeroReturnsNoneOnNil(t *testing.T) {
	if v := credentials.OrZero(nil); v != credentials.None {
		t.Errorf("credentials.OrZero(nil) should've returned credentials.None, not %v!", v)
	}
}

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
