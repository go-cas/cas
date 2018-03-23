package cas

import (
	"net/url"
	"path"
	"net/http"
	"github.com/golang/glog"
	"io/ioutil"
	"fmt"
)

func ValidateTicket(client *http.Client, casUrl *url.URL, ticket string, serviceUrl *url.URL) (*AuthenticationResponse, error) {
	if glog.V(2) {
		glog.Infof("Validating ticket %v for service %v", ticket, serviceUrl)
	}

	u, err := serviceValidateUrl(casUrl, ticket, serviceUrl)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("User-Agent", "Golang CAS client gopkg.in/cas")

	if glog.V(2) {
		glog.Infof("Attempting ticket validation with %v", r.URL)
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if glog.V(2) {
		glog.Infof("Request %v %v returned %v",
		r.Method, r.URL,
		resp.Status)
	}

	if resp.StatusCode == http.StatusNotFound {
		return validateTicketCas1(client, casUrl, ticket, serviceUrl)
	}

	body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, err
		}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cas: validate ticket: %v", string(body))
	}

	if glog.V(2) {
		glog.Infof("Received authentication response\n%v", string(body))
	}

	success, err := ParseServiceResponse(body)
	if err != nil {
		return nil, err
	}

	if glog.V(2) {
		glog.Infof("Parsed ServiceResponse: %#v", success)
	}

	return success, nil
}

func serviceValidateUrl(casUrl *url.URL, ticket string, serviceUrl *url.URL) (string, error) {
	u, err := casUrl.Parse(path.Join(casUrl.Path, "serviceValidate"))
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("service", sanitisedURLString(serviceUrl))
	q.Add("ticket", ticket)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func validateUrl(casUrl *url.URL, ticket string, serviceUrl *url.URL) (string, error) {
	u, err := casUrl.Parse(path.Join(casUrl.Path, "validate"))
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("service", sanitisedURLString(serviceUrl))
	q.Add("ticket", ticket)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func validateTicketCas1(client *http.Client, casUrl *url.URL, ticket string, serviceUrl *url.URL) (*AuthenticationResponse, error) {
	u, err := validateUrl(casUrl, ticket, serviceUrl)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("User-Agent", "Golang CAS client gopkg.in/cas")

	if glog.V(2) {
		glog.Info("Attempting ticket validation with %v", r.URL)
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if glog.V(2) {
		glog.Info("Request %v %v returned %v",
			r.Method, r.URL,
			resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body := string(data)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cas: validate ticket: %v", body)
	}

	if glog.V(2) {
		glog.Infof("Received authentication response\n%v", body)
	}

	if body == "no\n\n" {
		return nil, nil // not logged in
	}

	success := &AuthenticationResponse{
		User: body[4 : len(body)-1],
	}

	if glog.V(2) {
		glog.Infof("Parsed ServiceResponse: %#v", success)
	}

	return success, nil
}