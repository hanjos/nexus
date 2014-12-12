// Package util store helper code that's useful but not... business logic, so to speak.
package util

import (
	"fmt"
	"github.com/hanjos/nexus/errors"
	"regexp"
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
