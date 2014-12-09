package errors

import (
	"fmt"
	"github.com/hanjos/nexus/credentials"
)

// BadResponseError is returned when there's an error on an attempt to access Nexus.
type BadResponseError struct {
	Url        string // e.g. http://nexus.somewhere.com
	StatusCode int    // e.g. 400
	Status     string // e.g. 400 Bad response
}

func (err BadResponseError) Error() string {
	return fmt.Sprintf("Bad response (%v) from %v", err.Status, err.Url)
}

// UnauthorizedError is returned when the given credentials aren't authorized to reach the given URL.
type UnauthorizedError struct {
	Url        string // e.g. http://nexus.somewhere.com
	Credentials credentials.Credentials // e.g. credentials.BasicAuth{"username", "password"}
}

func (err UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %v doesn't have access to %v", err.Credentials, err.Url)
}
