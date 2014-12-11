/*
 Package credentials provides an http.Request with a set of credentials. Some Nexus API calls can only be done by users
 with the proper authorization.
*/
package credentials

import "net/http"

// Credentials is satisfied by whoever can configure an http.Request properly.
type Credentials interface {
	// Provides an http.Request with a set of credentials for authorization.
	Sign(request *http.Request)
}

// None is the zero value for Credentials. Its Sign() removes Authorization data from the header.
const None = noCredentials(true)

// bool trick for Go to allow a const
type noCredentials bool

func (auth noCredentials) Sign(request *http.Request) {
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

// BasicAuth signs the header using HTTP Basic Authentication.
type BasicAuth struct {
	Username string
	Password string
}

// Sign implements the Credentials interface, signing the header using HTTP Basic Authentication.
func (auth BasicAuth) Sign(request *http.Request) {
	request.SetBasicAuth(auth.Username, auth.Password)
}

// String implements the Stringer interface, for easy printing.
func (auth BasicAuth) String() string {
	return "BasicAuth{" + auth.Username + ", ***}"
}
