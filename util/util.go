// Package util stores useful helper code.
package util

import (
	"fmt"
	"regexp"
	"strings"
)

// ByteSize represents an amount of bytes. float64 is needed since division
// is required.
type ByteSize float64

// Some pre-constructed file size units.
const (
	Byte ByteSize = 1 << (10 * iota)
	Kilobyte
	Megabyte
	Gigabyte
)

// String implements the fmt.Stringer interface.
func (size ByteSize) String() string {
	switch {
	case size <= Kilobyte:
		return fmt.Sprintf("%d B", size)
	case size <= Megabyte:
		return fmt.Sprintf("%.2f KB", size/Kilobyte)
	case size <= Gigabyte:
		return fmt.Sprintf("%.2f MB", size/Megabyte)
	default:
		return fmt.Sprintf("%.2f GB", size/Gigabyte)
	}
}

var urlRe = regexp.MustCompile(`^(?P<scheme>[^:]+)://(?P<rest>.+)`)
var slashesRe = regexp.MustCompile(`//+`)

// Removes extraneous slashes (like nexus.com///something), which Nexus' API doesn't recognize as valid.
// Returns an util.MalformedURLError if the given URL can't be parsed.
func cleanSlashes(url string) (string, error) {
	matches := urlRe.FindStringSubmatch(url)
	if matches == nil {
		return "", &MalformedURLError{url}
	}

	// if we got here, scheme = matches[1] and rest = matches[2]. Clean the extraneous slashes
	return matches[1] + "://" + slashesRe.ReplaceAllString(matches[2], "/"), nil
}

// BuildFullURL builds a complete URL string in the format host/path?query, where query's keys and values will be
// formatted as k=v. Returns an util.MalformedURLError if the given URL can't be parsed. This function is a (very)
// simplified version of url.URL.String().
func BuildFullURL(host string, path string, query map[string]string) (string, error) {
	params := []string{}

	for k, v := range query {
		params = append(params, k+"="+v)
	}

	if len(params) == 0 {
		return cleanSlashes(host + "/" + path)
	}

	return cleanSlashes(host + "/" + path + "?" + strings.Join(params, "&"))
}

// MalformedURLError is returned when the given URL could not be parsed.
type MalformedURLError struct {
	URL string // e.g. http:/:malformed.url.com
}

// Error implements the error interface.
func (err MalformedURLError) Error() string {
	return fmt.Sprintf("Malformed URL: %v", err.URL)
}
