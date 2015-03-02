package cas

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

const (
	sessionCookieName = "_cas_session"
)

type clientHandler struct {
	c *Client
	h http.Handler
}

func (ch *clientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if glog.V(2) {
		glog.Infof("cas: handling %v request for %v", r.Method, r.URL)
	}

	setClient(r, ch.c)
	defer clear(r)

	ch.c.getSession(w, r)
	ch.h.ServeHTTP(w, r)
	return
}
