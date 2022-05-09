package cas

import (
	"net/url"
)

var (
	urlCleanParameters = []string{"gateway", "renew", "service", "ticket"}
)

// sanitisedURL cleans a URL of CAS specific parameters
func sanitisedURL(unclean *url.URL) (*url.URL, error) {
	// Parse maybe occur errors, cause unclean is dealt with requestURL method
	u, err := url.Parse(unclean.String())
	if err != nil {
		return nil, err
	}
	q := u.Query()

	for _, param := range urlCleanParameters {
		q.Del(param)
	}

	u.RawQuery = q.Encode()
	return u, nil
}

// sanitisedURLString cleans a URL and returns its string value
func sanitisedURLString(unclean *url.URL) (string, error) {
	u, err := sanitisedURL(unclean)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
