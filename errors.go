package nexus

import "fmt"

// BadResponseError is returned when there's an error on an attempt to access Nexus.
type BadResponseError struct {
	Url        string // e.g. http://nexus.somewhere.com
	StatusCode int    // e.g. 400
	Status     string // e.g. 400 Bad response
}

func (err *BadResponseError) Error() string {
	return fmt.Sprintf("Bad response (%v) from %v", err.Status, err.Url)
}
