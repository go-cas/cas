package cas

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"net/url"
)

func TestRequestGrantingTicket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cas/v1/tickets" || r.Method != "POST" {
			w.WriteHeader(404)
			return
		}

		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			w.WriteHeader(415)
			return
		}

		if r.FormValue("username") != "tricia" && r.FormValue("password") != "hitchhiker" {
			w.WriteHeader(400)
			return
		}

		w.Header().Set("Location", "/cas/v1/tickets/TGT-abc")
		w.WriteHeader(201)
	}))
	defer server.Close()

	casUrl, err := url.Parse(server.URL + "/cas/")
	if err != nil {
		t.Error("failed to create cas url from test server")
	}

	restClient := NewRestClient(&RestOptions{
		CasURL: casUrl,
		Client: server.Client(),
	})

	tgt, err := restClient.RequestGrantingTicket("tricia", "hitchhiker")
	if err != nil {
		t.Errorf("requesting granting ticket failed: %v", err)
	}

	if tgt != "TGT-abc" {
		t.Errorf("expected %s but received %v", "TGT-abc", tgt)
	}

	_, err = restClient.RequestGrantingTicket("arthur", "dent")
	if err == nil {
		t.Errorf("authentication should fail for arthur")
	}
}

func TestRequestServiceTicket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cas/v1/tickets/TGT-abc" || r.Method != "POST" {
			w.WriteHeader(404)
			return
		}

		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			w.WriteHeader(415)
			return
		}

		if r.FormValue("service") != "https://hitchhiker.com/heartOfGold" {
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("ST-123"))
	}))
	defer server.Close()

	casUrl, err := url.Parse(server.URL + "/cas/")
	if err != nil {
		t.Error("failed to create cas url from test server")
	}

	restClient := NewRestClient(&RestOptions{
		CasURL: casUrl,
		Client: server.Client(),
	})

	serviceUrl, err := url.Parse("https://hitchhiker.com/heartOfGold")
	if err != nil {
		t.Error("failed to create service url")
	}

	st, err := restClient.RequestServiceTicket(TicketGrantingTicket("TGT-abc"), serviceUrl)
	if err != nil {
		t.Errorf("requesting service ticket failed: %v", err)
	}

	if st != "ST-123" {
		t.Errorf("expected %s but received %v", "ST-123", st)
	}

	_, err = restClient.RequestServiceTicket(TicketGrantingTicket("TGT-xyz"), serviceUrl)
	if err == nil {
		t.Errorf("service ticket request should fail for TGT-xyz")
	}

	serviceUrl, err = url.Parse("https://hitchhiker.com/restaurantAtTheEndOfTheUniverse")
	if err != nil {
		t.Error("failed to create service url")
	}

	_, err = restClient.RequestServiceTicket(TicketGrantingTicket("TGT-abc"), serviceUrl)
	if err == nil {
		t.Errorf("service ticket request should fail for this service")
	}
}

func TestLogout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cas/v1/tickets/TGT-abc" || r.Method != "DELETE" {
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)
	}))
	defer server.Close()

	casUrl, err := url.Parse(server.URL + "/cas/")
	if err != nil {
		t.Error("failed to create cas url from test server")
	}

	restClient := NewRestClient(&RestOptions{
		CasURL: casUrl,
		Client: server.Client(),
	})

	err = restClient.Logout(TicketGrantingTicket("TGT-abc"))
	if err != nil {
		t.Errorf("logout failed %v", err)
	}

	err = restClient.Logout(TicketGrantingTicket("TGT-xyz"))
	if err == nil {
		t.Errorf("logout should failed for this TGT")
	}
}