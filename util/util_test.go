package util

import (
	"reflect"
	"testing"
)

var cleanSlashesOK = []struct {
	input    string
	expected string
}{
	{"http://maven.java.net", "http://maven.java.net"},
	{"http://maven.java.net/", "http://maven.java.net/"},
	{"http://maven.java.net/////", "http://maven.java.net/"},
	{"http://maven.java.net/nexus", "http://maven.java.net/nexus"},
	{"http://maven.java.net/////nexus", "http://maven.java.net/nexus"},
}

var cleanSlashesErr = []struct {
	input    string
	expected reflect.Type
}{
	{"http:/maven.java.net", reflect.TypeOf(&MalformedURLError{})},
}

func TestCleanSlashes(t *testing.T) {
	for _, p := range cleanSlashesOK {
		actual, err := cleanSlashes(p.input)

		if err != nil {
			t.Errorf("expected %v, got an error %v", p.expected, err)
		} else if actual != p.expected {
			t.Errorf("expected %v, got %v", p.expected, actual)
		}
	}

	for _, p := range cleanSlashesErr {
		actual, err := cleanSlashes(p.input)

		if actual != "" {
			t.Errorf("expected an error %v, got a value %v", p.expected, actual)
		} else if reflect.TypeOf(err) != p.expected {
			t.Errorf("expected an error %v, got the error %v", p.expected, err)
		}
	}
}

var bfuOk = []struct {
	host  string
	path  string
	query map[string]string

	expected string
}{
	{"http://maven.java.net", "nexus", map[string]string{}, "http://maven.java.net/nexus"},
	{"http://maven.java.net", "///nexus", map[string]string{}, "http://maven.java.net/nexus"},
	{"http://maven.java.net////", "/nexus", map[string]string{}, "http://maven.java.net/nexus"},
	{"http://maven.java.net///", "/nexus", map[string]string{"p": "1"}, "http://maven.java.net/nexus?p=1"},
	{"http://maven.java.net///", "/nexus", map[string]string{"p": "1", "q": "2"}, "http://maven.java.net/nexus?p=1&q=2"},
}

var bfuErr = []struct {
	host  string
	path  string
	query map[string]string

	expected reflect.Type
}{
	{"http:/maven.java.net", "/nexus", map[string]string{}, reflect.TypeOf(&MalformedURLError{})},
}

func TestBuildFullURL(t *testing.T) {
	for _, p := range bfuOk {
		actual, err := BuildFullURL(p.host, p.path, p.query)

		if err != nil {
			t.Errorf("expected %v, got an error %v", p.expected, err)
		} else if actual != p.expected {
			t.Errorf("expected %v, got %v", p.expected, actual)
		}
	}

	for _, p := range bfuErr {
		actual, err := BuildFullURL(p.host, p.path, p.query)

		if actual != "" {
			t.Errorf("expected an error %v, got a value %v", p.expected, actual)
		} else if reflect.TypeOf(err) != p.expected {
			t.Errorf("expected an error %v, got an error %v", p.expected, err)
		}
	}
}

func TestIfMalformedURLErrorIsError(t *testing.T) {
	// type assertion only works on interface types, so...
	if _, ok := interface{}(&MalformedURLError{}).(error); !ok {
		t.Errorf("util.MalformedURLError does not implement the error interface!")
	}
}
