package cas

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUnauthenticatedRequestShouldRedirectToCasURL(t *testing.T) {
	url, _ := url.Parse("https://cas.example.com/")
	client := NewClient(&Options{
		URL: url,
	})

	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogin(w, r)
			return
		}

		fmt.Fprintln(w, "You are logged in, but you shouldn't be, oh noes!!")
	})

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected HTTP response code to be <%v>, got <%v>", http.StatusFound, w.Code)
	}

	loc := w.Header().Get("Location")
	exp := "https://cas.example.com/login?service=http%3A%2F%2Fexample.com%2F"
	if loc != exp {
		t.Errorf("Expected HTTP redirect to <%s>, got <%s>", exp, loc)
	}

	setCookie := w.Header().Get("Set-Cookie")
	if !strings.HasPrefix(setCookie, sessionCookieName) {
		t.Errorf("Expected response to have Set-Cookie header with <%v>, got <%v>",
			sessionCookieName, setCookie)
	}
}

func TestInvalidServiceTicket(t *testing.T) {
	server := &TestServer{}
	defer server.Close()
	ts := httptest.NewServer(server)
	defer ts.Close()

	url, _ := url.Parse(ts.URL)
	client := NewClient(&Options{
		URL: url,
	})

	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogin(w, r)
			return
		}

		fmt.Fprintln(w, "You are logged in, but you shouldn't be, oh noes!!")
	})

	req, err := http.NewRequest("GET", "http://example.com/?ticket=ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected HTTP response code to be <%v>, got <%v>", http.StatusFound, w.Code)
	}

	loc := w.Header().Get("Location")
	exp, _ := url.Parse("/login?service=http%3A%2F%2Fexample.com%2F")
	if loc != exp.String() {
		t.Errorf("Expected HTTP redirect to <%s>, got <%s>", exp, loc)
	}

	setCookie := w.Header().Get("Set-Cookie")
	if !strings.HasPrefix(setCookie, sessionCookieName) {
		t.Errorf("Expected response to have Set-Cookie header with <%v>, got <%v>",
			sessionCookieName, setCookie)
	}
}

func TestValidServiceTicket(t *testing.T) {
	server := &TestServer{}
	ticket := server.NewTicket("ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c")
	ticket.Service = "http://example.com/"
	ticket.Username = "TestValidServiceTicket"
	server.AddTicket(ticket)

	defer server.Close()

	ts := httptest.NewServer(server)
	defer ts.Close()

	url, _ := url.Parse(ts.URL)
	client := NewClient(&Options{
		URL: url,
	})

	message := "You are logged in, welcome client"
	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogin(w, r)
			return
		}

		fmt.Fprintln(w, message)
	})

	req, err := http.NewRequest("GET", "http://example.com/?ticket=ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP response code to be <%v>, got <%v>", http.StatusOK, w.Code)
	}

	if message != strings.Trim(w.Body.String(), "\n") {
		t.Errorf("Expected body to be <%s>, got <%s>", message, strings.Trim(w.Body.String(), "\n"))
	}

	setCookie := w.Header().Get("Set-Cookie")
	if !strings.HasPrefix(setCookie, sessionCookieName) {
		t.Errorf("Expected response to have Set-Cookie header with <%v>, got <%v>",
			sessionCookieName, setCookie)
	}
}

func TestGetUsernameFromServiceTicket(t *testing.T) {
	server := &TestServer{}
	ticket := server.NewTicket("ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c")
	ticket.Service = "http://example.com/"
	ticket.Username = "enoch.root"
	server.AddTicket(ticket)
	defer server.Close()

	ts := httptest.NewServer(server)
	defer ts.Close()

	url, _ := url.Parse(ts.URL)
	client := NewClient(&Options{
		URL: url,
	})

	message := "You are logged in, welcome"
	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogin(w, r)
			return
		}

		user := Username(r)
		fmt.Fprintln(w, message, user)
	})

	req, err := http.NewRequest("GET", "http://example.com/?ticket=ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP response code to be <%v>, got <%v>", http.StatusOK, w.Code)
	}

	expected := fmt.Sprintf("%s %s", message, ticket.Username)
	if expected != strings.Trim(w.Body.String(), "\n") {
		t.Errorf("Expected body to be <%s>, got <%s>", expected, strings.Trim(w.Body.String(), "\n"))
	}

	setCookie := w.Header().Get("Set-Cookie")
	if !strings.HasPrefix(setCookie, sessionCookieName) {
		t.Errorf("Expected response to have Set-Cookie header with <%v>, got <%v>",
			sessionCookieName, setCookie)
	}
}

func TestGetAttributesFromServiceTicket(t *testing.T) {
	server := &TestServer{}
	ticket := server.NewTicket("ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c")
	ticket.Service = "http://example.com/"
	ticket.Username = "enoch.root"
	ticket.Attributes.Add("admin", "true")
	ticket.Attributes.Add("account", "testing")
	server.AddTicket(ticket)
	defer server.Close()

	ts := httptest.NewServer(server)
	defer ts.Close()

	url, _ := url.Parse(ts.URL)
	client := NewClient(&Options{
		URL: url,
	})

	message := "You are logged in, welcome %s%s, your account is %s"
	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogin(w, r)
			return
		}

		user := Username(r)
		attr := Attributes(r)

		admin := ""
		if attr.Get("admin") == "true" {
			admin = "Sir "
		}

		account := attr.Get("account")
		fmt.Fprintf(w, message, admin, user, account)
		fmt.Fprintf(w, "\n")
	})

	req, err := http.NewRequest("GET", "http://example.com/?ticket=ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP response code to be <%v>, got <%v>", http.StatusOK, w.Code)
	}

	expected := fmt.Sprintf(message, "Sir ", ticket.Username, "testing")
	if expected != strings.Trim(w.Body.String(), "\n") {
		t.Errorf("Expected body to be <%s>, got <%s>", expected, strings.Trim(w.Body.String(), "\n"))
	}

	setCookie := w.Header().Get("Set-Cookie")
	if !strings.HasPrefix(setCookie, sessionCookieName) {
		t.Errorf("Expected response to have Set-Cookie header with <%v>, got <%v>",
			sessionCookieName, setCookie)
	}
}

func TestSecondRequestShouldBeCookied(t *testing.T) {
	server := &TestServer{}
	ticket := server.NewTicket("ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c")
	ticket.Service = "http://example.com/"
	ticket.Username = "enoch.root"
	ticket.Attributes.Add("admin", "true")
	ticket.Attributes.Add("account", "testing")
	server.AddTicket(ticket)
	defer server.Close()

	ts := httptest.NewServer(server)
	defer ts.Close()

	url, _ := url.Parse(ts.URL)
	client := NewClient(&Options{
		URL: url,
	})

	message := "You are logged in, welcome %s%s, your account is %s"
	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogin(w, r)
			return
		}

		user := Username(r)
		attr := Attributes(r)

		admin := ""
		if attr.Get("admin") == "true" {
			admin = "Sir "
		}

		account := attr.Get("account")
		fmt.Fprintf(w, message, admin, user, account)
		fmt.Fprintf(w, "\n")
	})

	req, err := http.NewRequest("GET", "http://example.com/?ticket=ST-l8d6b51d8e9c4569345a30e2f904626a1066384db8694784a60b515d62f6c", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected First HTTP response code to be <%v>, got <%v>", http.StatusOK, w.Code)
	}

	setCookie := w.Header().Get("Set-Cookie")
	if !strings.HasPrefix(setCookie, sessionCookieName) {
		t.Errorf("Expected response to have Set-Cookie header with <%v>, got <%v>",
			sessionCookieName, setCookie)
	}

	req, err = http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Error(err)
	}

	// Parse response headers and add them to the new request
	resp := http.Response{Header: w.Header()}
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected Second HTTP response code to be <%v>, got <%v>", http.StatusOK, w.Code)
	}
}
