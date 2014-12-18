package util

import (
	"reflect"
	"testing"
)

type pair struct {
	input     string
	expected  string
	errorType reflect.Type
}

var cleanSlashesPairs = []pair{
	{"http://maven.java.net", "http://maven.java.net", nil},
	{"http://maven.java.net/", "http://maven.java.net/", nil},
	{"http://maven.java.net/////", "http://maven.java.net/", nil},
	{"http://maven.java.net/nexus", "http://maven.java.net/nexus", nil},
	{"http://maven.java.net/////nexus", "http://maven.java.net/nexus", nil},
	{"http:/maven.java.net", "", reflect.TypeOf(&MalformedUrlError{})},
}

func TestCleanSlashes(t *testing.T) {
	for _, p := range cleanSlashesPairs {
		actual, err := cleanSlashes(p.input)

		if p.expected != "" { // a value should have been returned
			if err != nil {
				t.Errorf("expected %v, got an error: %v", p.expected, err)
			}
			if actual != p.expected {
				t.Errorf("expected %v, got %v", p.expected, actual)
			}
		} else if p.errorType != nil { // an error should have been returned
			if actual != "" {
				t.Errorf("expected \"\" and an error, got %v", actual)
			}

			actualErrorType := reflect.TypeOf(err)
			if p.errorType != actualErrorType {
				t.Errorf("expected an error %v, got %v", p.errorType, actualErrorType)
			}
		} else { // if we're here, the test is broken
			t.Error("Test malformed; check the pairs")
		}
	}
}
