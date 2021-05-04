package cas

import (
	"net/url"
	"testing"
)

func TestDefaultURLScheme(t *testing.T) {
	url, _ := url.Parse("https://cas.org/cas")
	scheme := NewDefaultURLScheme(url)

	u, err := scheme.Login()
	assertURL(t, "/cas/login", u, err)
	u, err = scheme.Logout()
	assertURL(t, "/cas/logout", u, err)
	u, err = scheme.Validate()
	assertURL(t, "/cas/validate", u, err)
	u, err = scheme.ServiceValidate()
	assertURL(t, "/cas/serviceValidate", u, err)
	u, err = scheme.RestGrantingTicket()
	assertURL(t, "/cas/v1/tickets", u, err)
	u, err = scheme.RestServiceTicket("TGT-123")
	assertURL(t, "/cas/v1/tickets/TGT-123", u, err)
	u, err = scheme.RestLogout("TGT-123")
	assertURL(t, "/cas/v1/tickets/TGT-123", u, err)
}

func assertURL(t *testing.T, expected string, u *url.URL, err error) {
	if err != nil {
		t.Fatalf("returned error")
	}

	if expected != u.Path {
		t.Errorf("%s should be equal to %s", u.Path, expected)
	}
}
