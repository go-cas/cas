package cas_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/cas.v1"
)

func ExampleRedirectToLogin() {
	u, _ := url.Parse("https://cas.example.com")
	c := cas.NewClient(&cas.Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		cas.RedirectToLogin(w, r)
	})

	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ExampleRedirectToLogout() {
	u, _ := url.Parse("https://cas.example.com")
	c := cas.NewClient(&cas.Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		cas.RedirectToLogout(w, r)
	})

	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ExampleIsAuthenticated() {
	u, _ := url.Parse("https://cas.example.com")
	c := cas.NewClient(&cas.Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cas.IsAuthenticated(r) {
			cas.RedirectToLogout(w, r)
		}

		fmt.Fprintf(w, "Hello World\n")
	})

	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ExampleUsername() {
	u, _ := url.Parse("https://cas.example.com")
	c := cas.NewClient(&cas.Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cas.IsAuthenticated(r) {
			cas.RedirectToLogout(w, r)
		}

		fmt.Fprintf(w, "Hello %s\n", cas.Username(r))
	})

	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
