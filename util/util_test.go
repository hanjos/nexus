package util

import (
	"reflect"
	"testing"
)

type pair interface {
	expected() string
	errorType() reflect.Type
}

type simplePair struct {
	input string

	_expected  string
	_errorType reflect.Type
}

func (p simplePair) expected() string        { return p._expected }
func (p simplePair) errorType() reflect.Type { return p._errorType }

var cleanSlashesPairs = []simplePair{
	{"http://maven.java.net", "http://maven.java.net", nil},
	{"http://maven.java.net/", "http://maven.java.net/", nil},
	{"http://maven.java.net/////", "http://maven.java.net/", nil},
	{"http://maven.java.net/nexus", "http://maven.java.net/nexus", nil},
	{"http://maven.java.net/////nexus", "http://maven.java.net/nexus", nil},
	{"http:/maven.java.net", "", reflect.TypeOf(&MalformedUrlError{})},
}

func checkResults(t *testing.T, p pair, actual string, err error) {
	if p.expected() != "" { // a value should have been returned
		if err != nil {
			t.Errorf("expected %v, got an error: %v", p.expected, err)
		}
		if actual != p.expected() {
			t.Errorf("expected %v, got %v", p.expected, actual)
		}
	} else if p.errorType() != nil { // an error should have been returned
		if actual != "" {
			t.Errorf("expected \"\" and an error, got %v", actual)
		}

		actualErrorType := reflect.TypeOf(err)
		if p.errorType() != actualErrorType {
			t.Errorf("expected an error %v, got %v", p.errorType(), actualErrorType)
		}
	} else { // if we're here, the test is broken
		t.Errorf("Test malformed; check the pair: %v", p)
	}
}

func TestCleanSlashes(t *testing.T) {
	for _, p := range cleanSlashesPairs {
		actual, err := cleanSlashes(p.input)

		checkResults(t, p, actual, err)
	}
}

type bfuInput struct {
	host  string
	path  string
	query map[string]string
}

type bfuPair struct {
	input bfuInput

	_expected  string
	_errorType reflect.Type
}

func (p bfuPair) expected() string        { return p._expected }
func (p bfuPair) errorType() reflect.Type { return p._errorType }

var bfuPairs = []bfuPair{
	{bfuInput{"http://maven.java.net", "nexus", map[string]string{}}, "http://maven.java.net/nexus", nil},
	{bfuInput{"http://maven.java.net", "///nexus", map[string]string{}}, "http://maven.java.net/nexus", nil},
	{bfuInput{"http://maven.java.net////", "/nexus", map[string]string{}}, "http://maven.java.net/nexus", nil},
	{bfuInput{"http:/maven.java.net", "/nexus", map[string]string{}}, "", reflect.TypeOf(&MalformedUrlError{})},
	{bfuInput{"http://maven.java.net///", "/nexus", map[string]string{"p": "1", "q": "2"}}, "http://maven.java.net/nexus?p=1&q=2", nil},
}

func TestBuildFullUrl(t *testing.T) {
	for _, p := range bfuPairs {
		actual, err := BuildFullUrl(p.input.host, p.input.path, p.input.query)

		checkResults(t, p, actual, err)
	}
}

func TestIfMalformedUrlErrorIsError(t *testing.T) {
	// type assertion only works on interface types, so...
	if _, ok := interface{}(&MalformedUrlError{}).(error); !ok {
		t.Errorf("util.MalformedUrlError does not implement the error interface!")
	}
}
