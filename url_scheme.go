package cas

import (
	"net/url"
	"path"
)

// URLScheme creates the url which are required to handle the cas protocol.
type URLScheme interface {
	Login() 		  (*url.URL, error)
	Logout() 		  (*url.URL, error)
	Validate() 		  (*url.URL, error)
	ServiceValidate() (*url.URL, error)
}

// NewDefaultURLScheme creates a URLScheme which uses the cas default urls
func NewDefaultURLScheme(base *url.URL) *DefaultURLScheme {
	return &DefaultURLScheme{
		base: 				 base,
		LoginPath: 			 "login",
		LogoutPath: 	     "logout",
		ValidatePath: 		 "validate",
		ServiceValidatePath: "serviceValidate",
	}
}

// DefaultURLScheme is a configurable URLScheme. Use NewDefaultURLScheme to create DefaultURLScheme with the default cas
// urls.
type DefaultURLScheme struct {
	base 				*url.URL
	LoginPath 			string
	LogoutPath 			string
	ValidatePath 		string
	ServiceValidatePath string
}

// Login returns the url for the cas login page
func (scheme *DefaultURLScheme) Login() (*url.URL, error) {
	return scheme.createURL(scheme.LoginPath)
}

// Logout returns the url for the cas logut page
func (scheme *DefaultURLScheme) Logout() (*url.URL, error) {
	return scheme.createURL(scheme.LogoutPath)
}

// Validate returns the url for the request validation endpoint
func (scheme *DefaultURLScheme) Validate() (*url.URL, error) {
	return scheme.createURL(scheme.ValidatePath)
}

// ServiceValidate returns the url for the service validation endpoint
func (scheme *DefaultURLScheme) ServiceValidate() (*url.URL, error) {
	return scheme.createURL(scheme.ServiceValidatePath)
}

func (scheme *DefaultURLScheme) createURL(urlPath string) (*url.URL, error) {
	return scheme.base.Parse(path.Join(scheme.base.Path, urlPath))
}