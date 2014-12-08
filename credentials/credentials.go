/*
 Package credentials provides code to provide an http.Request with a set of credentials. Some API calls can only be
 done by users with the proper authorization, so a nexus.Client must have some way of providing it.
*/
package credentials

import "net/http"

// Credentials is satisfied by whoever can configure an http.Request properly.
type Credentials interface {
	// Provides an http.Request with a set of credentials for authorization.
	Sign(request *http.Request)
}

// None is the zero value for Credentials. It removes Authorization data from the header.
const None = noCredentials(false)

type noCredentials bool

func (auth noCredentials) Sign(request *http.Request) {
	request.Header.Del("Authorization")
}

// OrZero returns the given credentials untouched if it's not nil, and credentials.None otherwise. Useful for when one
// must ensure that a given set of credentials is non-nil.
func OrZero(c Credentials) Credentials {
	if c == nil {
		return None
	}

	return c
}

// BasicAuth signs the header using HTTP Basic Authentication.
type BasicAuth struct {
	Username string
	Password string
}

func (auth BasicAuth) Sign(request *http.Request) {
	request.SetBasicAuth(auth.Username, auth.Password)
}
