package cas

import (
	"crypto/rand"
	"net/http"
	"sync"
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
	setClient(r, ch.c)
	defer clear(r)

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

		http.SetCookie(w, cookie)
	}

	loggedIn := false

	if ticket, ok := ch.seen[cookie.Value]; ok {
		// Set ticket on request
		// Store.Read(ticket) if this fails then we need to nuke
		auth, _ := ch.c.store.Read(ticket)
		setAuthenticationResponse(r, auth)

		loggedIn = true
	}

	if ticket := r.URL.Query().Get("ticket"); !loggedIn && ticket != "" {
		if err := ch.c.validateTicket(ticket, r.URL); err != nil {
			// log error, invalid ticket, service or something
			// allow them up, but don't set anything to show them
			// as logged in to the higher layers
			ch.h.ServeHTTP(w, r)
			return
		}

		ch.mu.Lock()
		ch.seen[cookie.Value] = ticket
		ch.mu.Unlock()

		auth, _ := ch.c.store.Read(ticket)
		setAuthenticationResponse(r, auth)

		loggedIn = true
	}

	ch.h.ServeHTTP(w, r)
	return
}
