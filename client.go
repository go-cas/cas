package cas

import (
	"net/http"
	"net/url"
)

type Client struct {
	URL *url.URL
}

func NewClient(options *Options) *Client {
	return &Client{URL: options.URL}
}

func (c *Client) Handle(h http.Handler) http.Handler {
	return &clientHandler{c: c, h: h}
}

func (c *Client) HandleFunc(h func(http.ResponseWriter, *http.Request)) http.Handler {
	return &clientHandler{c: c, h: http.HandlerFunc(h)}
}

type clientHandler struct {
	c *Client
	h http.Handler
}

func (ch *clientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.h.ServeHTTP(w, r)
}
