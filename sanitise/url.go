package sanitise

import (
	"net/url"
)

var (
	urlCleanParameters = []string{"gateway", "renew", "service", "ticket"}
)

func URL(unclean *url.URL) *url.URL {
	// Shouldn't be any errors parsing an existing *url.URL
	u, _ := url.Parse(unclean.String())
	q := u.Query()

	for _, param := range urlCleanParameters {
		q.Del(param)
	}

	u.RawQuery = q.Encode()
	return u
}

func URLString(unclean *url.URL) string {
	return URL(unclean).String()
}
