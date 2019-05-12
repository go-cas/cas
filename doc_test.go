package cas

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func ExampleRedirectToLogin() {
	u, _ := url.Parse("https://cas.example.com")
	c := NewClient(&Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		RedirectToLogin(w, r)
	})

	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ExampleRedirectToLogout() {
	u, _ := url.Parse("https://cas.example.com")
	c := NewClient(&Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		RedirectToLogout(w, r)
	})

	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ExampleIsAuthenticated() {
	u, _ := url.Parse("https://cas.example.com")
	c := NewClient(&Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogout(w, r)
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
	c := NewClient(&Options{
		URL: u,
	})

	h := c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			RedirectToLogout(w, r)
		}

		fmt.Fprintf(w, "Hello %s\n", Username(r))
	})

	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
