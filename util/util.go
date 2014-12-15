// Package util stores useful helper code.
package util

import (
	"fmt"
	"github.com/hanjos/nexus/errors"
	"regexp"
	"strings"
)

// FileSize represents an amount of bytes.
type FileSize int64

const (
	Byte     = FileSize(1)
	Kilobyte = FileSize(1 << 10)
	Megabyte = FileSize(1 << 20)
	Gigabyte = FileSize(1 << 30)
)

// String implements the Stringer interface, for easy printing.
func (size FileSize) String() string {
	switch true {
	case size <= Kilobyte:
		return fmt.Sprintf("%d B", int(size))
	case size <= Megabyte:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(Kilobyte))
	case size <= Gigabyte:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(Megabyte))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(Gigabyte))
	}
}

var urlRe = regexp.MustCompile("^(?P<scheme>[^:]+)://(?P<rest>.+)")
var slashesRe = regexp.MustCompile("//+")

// CleanSlashes removes extraneous slashes (like nexus.com///something), which Nexus' API doesn't recognize as valid.
// Returns an errors.MalformedUrlError if the given URL can't be parsed.
func CleanSlashes(url string) (string, error) {
	if !urlRe.MatchString(url) {
		return "", &errors.MalformedUrlError{url}
	}

	scheme := urlRe.ReplaceAllString(url, "${scheme}")
	rest := urlRe.ReplaceAllString(url, "${rest}")

	return scheme + "://" + slashesRe.ReplaceAllString(rest, "/"), nil
}

// BuildFullUrl builds a complete URL string in the format host/path?query, where query's keys and values will be
// formatted as k=v. This function is a (very) simplified version of url.URL.String().
func BuildFullUrl(host string, path string, query map[string]string) string {
	params := []string{}

	for k, v := range query {
		params = append(params, k+"="+v)
	}

	if len(params) == 0 {
		return host + "/" + path
	} else {
		return host + "/" + path + "?" + strings.Join(params, "&")
	}
}
