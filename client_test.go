package cas

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLearnHttpTestServer(t *testing.T) {
	expected := "Hello, client"

	url, _ := url.Parse("https://cas.host")
	client := NewClient(&Options{
		URL: url,
	})

	ts := httptest.NewServer(client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expected)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Error(err)
	}

	if expected != strings.Trim(string(body), "\n") {
		t.Errorf("expected body to equal <%s>, got <%s>", expected, body)
	}
}

func TestLearnHttpTestResponseRecorder(t *testing.T) {
	expected := "Hello, client"

	url, _ := url.Parse("https://cas.host")
	client := NewClient(&Options{
		URL: url,
	})

	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expected)
	})

	req, err := http.NewRequest("GET", "http://test.host/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	body := w.Body.String()

	if expected != strings.Trim(body, "\n") {
		t.Errorf("expected body to equal <%s>, got <%s>", expected, body)
	}
}
