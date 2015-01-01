/*
Package credentials provides credentials to an http.Request. Some Nexus API calls can only be done by users
with the proper authorization.
*/
package credentials

import (
	"fmt"
	"net/http"
)

// Credentials represents credentials which can be provided to an http.Request.
type Credentials interface {
	// Adds a set of credentials to an http.Request for authorization.
	// Does nothing on nil requests.
	Sign(request *http.Request)
}

// None is the zero value for Credentials. Its Sign() removes Authorization data from the header.
const None = noCredentials(true)

// bool trick for Go to allow a const
type noCredentials bool

func (auth noCredentials) Sign(request *http.Request) {
	if request == nil {
		return
	}

	request.Header.Del("Authorization")
}

func (auth noCredentials) String() string {
	return "No credentials"
}

// OrZero returns the given credentials untouched if it's not nil, and credentials.None otherwise. Useful for when one
// must ensure that a given set of credentials is non-nil.
func OrZero(c Credentials) Credentials {
	if c == nil {
		return None
	}

	return c
}

type basicAuth struct {
	Username string
	Password string
}

// BasicAuth returns a Credentials which signs the header using HTTP Basic Authentication.
func BasicAuth(username, password string) Credentials {
	return basicAuth{Username: username, Password: password}
}

func (auth basicAuth) Sign(request *http.Request) {
	if request == nil {
		return
	}

	request.SetBasicAuth(auth.Username, auth.Password)
}

func (auth basicAuth) String() string {
	return "BasicAuth{" + auth.Username + ", ***}"
}

// Error is returned when the given credentials aren't authorized to reach the given URL.
type Error struct {
	URL         string      // e.g. http://nexus.somewhere.com
	Credentials Credentials // e.g. credentials.BasicAuth("username", "password")
}

// Error implements the error interface.
func (err Error) Error() string {
	// err.Credentials may not implement fmt.Stringer, so fmt.Sprintf is safer to use
	return fmt.Sprintf("%v doesn't have access to %v", err.Credentials, err.URL)
}
