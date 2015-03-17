package cas

import (
	"net/http"
	"sync"
	"time"
)

var (
	mutex   sync.RWMutex
	clients = make(map[*http.Request]*Client)
	data    = make(map[*http.Request]*AuthenticationResponse)
)

// setClient associates a Client with a http.Request.
func setClient(r *http.Request, c *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	clients[r] = c
}

// getClient retrieves the Client associated with the http.Request.
func getClient(r *http.Request) *Client {
	mutex.RLock()
	defer mutex.RUnlock()

	return clients[r]
}

// RedirectToLogin allows CAS protected handlers to redirect a request
// to the CAS login page.
func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	c := getClient(r)
	if c == nil {
		err := "cas: redirect to cas failed as no client associated with request"
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	c.RedirectToLogin(w, r)
}

// RedirectToLogout allows CAS protected handlers to redirect a request
// to the CAS logout page.
func RedirectToLogout(w http.ResponseWriter, r *http.Request) {
	c := getClient(r)
	if c == nil {
		err := "cas: redirect to cas failed as no client associated with request"
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	c.RedirectToLogout(w, r)
}

// setAuthenticationResponse associates an AuthenticationResponse with
// a http.Request.
func setAuthenticationResponse(r *http.Request, a *AuthenticationResponse) {
	mutex.Lock()
	defer mutex.Unlock()

	data[r] = a
}

// getAuthenticationResponse retrieves the AuthenticationResponse associated
// with a http.Request.
func getAuthenticationResponse(r *http.Request) *AuthenticationResponse {
	mutex.RLock()
	defer mutex.RUnlock()

	return data[r]
}

// clear removes the Client and AuthenticationResponse associations of a http.Request.
func clear(r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(clients, r)
	delete(data, r)
}

// IsAuthenticated indicates whether the request has been authenticated with CAS.
func IsAuthenticated(r *http.Request) bool {
	if a := getAuthenticationResponse(r); a != nil {
		return true
	}

	return false
}

// Username returns the authenticated users username
func Username(r *http.Request) string {
	if a := getAuthenticationResponse(r); a != nil {
		return a.User
	}

	return ""
}

// Attributes returns the authenticated users attributes.
func Attributes(r *http.Request) UserAttributes {
	if a := getAuthenticationResponse(r); a != nil {
		return a.Attributes
	}

	return nil
}

// AuthenticationDate returns the date and time that authentication was performed.
//
// This may return time.IsZero if Authentication Date information is not included
// in the CAS service validation response. This will be the case for CAS 2.0
// protocol servers.
func AuthenticationDate(r *http.Request) time.Time {
	var t time.Time
	if a := getAuthenticationResponse(r); a != nil {
		t = a.AuthenticationDate
	}

	return t
}

// IsNewLogin indicates whether the CAS service ticket was granted following a
// new authentication.
//
// This may incorrectly return false if Is New Login information is not included
// in the CAS service validation response. This will be the case for CAS 2.0
// protocol servers.
func IsNewLogin(r *http.Request) bool {
	if a := getAuthenticationResponse(r); a != nil {
		return a.IsNewLogin
	}

	return false
}

// IsRememberedLogin indicates whether the CAS service ticket was granted by the
// presence of a long term authentication token.
//
// This may incorrectly return false if Remembered Login information is not included
// in the CAS service validation response. This will be the case for CAS 2.0
// protocol servers.
func IsRememberedLogin(r *http.Request) bool {
	if a := getAuthenticationResponse(r); a != nil {
		return a.IsRememberedLogin
	}

	return false
}

// MemberOf returns the list of groups which the user belongs to.
func MemberOf(r *http.Request) []string {
	if a := getAuthenticationResponse(r); a != nil {
		return a.MemberOf
	}

	return nil
}
