package cas

import (
	"crypto/rand"
	"net/http"
	"sync"

	"github.com/golang/glog"
)

const (
	sessionCookieName = "_cas_session"
)

type clientHandler struct {
	c    *Client
	h    http.Handler
	mu   sync.Mutex
	seen map[string]string
}

func newSessionId() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// generate 64 character string
	bytes := make([]byte, 64)
	rand.Read(bytes)

	for k, v := range bytes {
		bytes[k] = alphabet[v%byte(len(alphabet))]
	}

	return string(bytes)
}

func (ch *clientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if glog.V(2) {
		glog.Infof("cas: handling %v request for %v", r.Method, r.URL)
	}

	setClient(r, ch.c)
	defer clear(r)

	if glog.V(2) {
		glog.Infof("Checking request for %v cookie", sessionCookieName)
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		// NOTE: Intentionally not enabling HttpOnly so the cookie can
		//       still be used by Ajax requests.
		cookie = &http.Cookie{
			Name:     sessionCookieName,
			Value:    newSessionId(),
			MaxAge:   86400,
			HttpOnly: false,
		}

		if glog.V(2) {
			glog.Infof("Setting %v cookie with value: %v", cookie.Name, cookie.Value)
		}

		http.SetCookie(w, cookie)
	}

	loggedIn := false

	if glog.V(2) {
		glog.Infof("Checking ticket cache with cookie %v", cookie.Value)
	}

	if ticket, ok := ch.seen[cookie.Value]; ok {
		if glog.V(2) {
			glog.Infof("Found ticket for cookie %v, ticket is %v", cookie.Value, ticket)
			glog.Info("Retrieving ticket response from store")
		}

		// Set ticket on request
		if auth, err := ch.c.store.Read(ticket); err == nil {
			setAuthenticationResponse(r, auth)

			loggedIn = true
		} else {
			if glog.V(2) {
				glog.Infof("Ticket %v not in store", ticket)
			}
		}
	}

	if glog.V(2) {
		if !loggedIn {
			glog.Info("Checking request URL ticket parameter")
		}
	}

	if ticket := r.URL.Query().Get("ticket"); !loggedIn && ticket != "" {
		if err := ch.c.validateTicket(ticket, r.URL); err != nil {
			// log error, invalid ticket, service or something
			// allow them up, but don't set anything to show them
			// as logged in to the higher layers
			ch.h.ServeHTTP(w, r)
			return
		}

		if glog.V(2) {
			glog.Infof("Recording ticket in cache, %v -> %v", cookie.Value, ticket)
		}

		ch.mu.Lock()
		ch.seen[cookie.Value] = ticket
		ch.mu.Unlock()

		if auth, err := ch.c.store.Read(ticket); err == nil {
			setAuthenticationResponse(r, auth)

			loggedIn = true
		} else {
			if glog.V(2) {
				glog.Infof("Ticket %v not in store after ticket validated", ticket)
			}
		}
	}

	ch.h.ServeHTTP(w, r)
	return
}
