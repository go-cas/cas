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

func setClient(r *http.Request, c *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	clients[r] = c
}

func getClient(r *http.Request) *Client {
	mutex.RLock()
	defer mutex.RUnlock()

	return clients[r]
}

func RedirectToCas(w http.ResponseWriter, r *http.Request) {
	c := getClient(r)
	if c == nil {
		err := "cas: redirect to cas failed as no client associated with request"
		http.Error(w, err, http.StatusInternalServerError)
		return
	}

	c.RedirectToCas(w, r)
}

func setAuthenticationResponse(r *http.Request, a *AuthenticationResponse) {
	mutex.Lock()
	defer mutex.Unlock()

	data[r] = a
}

func getAuthenticationResponse(r *http.Request) *AuthenticationResponse {
	mutex.RLock()
	defer mutex.RUnlock()

	return data[r]
}

func clear(r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(clients, r)
	delete(data, r)
}

func IsAuthenticated(r *http.Request) bool {
	if a := getAuthenticationResponse(r); a != nil {
		return true
	}

	return false
}

func Username(r *http.Request) string {
	if a := getAuthenticationResponse(r); a != nil {
		return a.User
	}

	return ""
}

func Attributes(r *http.Request) UserAttributes {
	if a := getAuthenticationResponse(r); a != nil {
		return a.Attributes
	}

	return nil
}

func AuthenticationDate(r *http.Request) time.Time {
	var t time.Time
	if a := getAuthenticationResponse(r); a != nil {
		t = a.AuthenticationDate
	}

	return t
}

func IsNewLogin(r *http.Request) bool {
	if a := getAuthenticationResponse(r); a != nil {
		return a.IsNewLogin
	}

	return false
}

func IsRememberedLogin(r *http.Request) bool {
	if a := getAuthenticationResponse(r); a != nil {
		return a.IsRememberedLogin
	}

	return false
}

func MemberOf(r *http.Request) []string {
	if a := getAuthenticationResponse(r); a != nil {
		return a.MemberOf
	}

	return nil
}
