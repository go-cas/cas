package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/golang/glog"

	"gopkg.in/cas.v1"
)

type myHandler struct{}

var MyHandler = &myHandler{}
var casURL string

func init() {
	flag.StringVar(&casURL, "url", "", "CAS server URL")
}

func main() {
	flag.Parse()

	if casURL == "" {
		flag.Usage()
		return
	}

	glog.Info("Starting up")

	m := http.NewServeMux()
	m.Handle("/", MyHandler)

	url, _ := url.Parse(casURL)
	client := cas.NewClient(&cas.Options{
		URL: url,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: client.Handle(m),
	}

	if err := server.ListenAndServe(); err != nil {
		glog.Infof("Error from HTTP Server: %v", err)
	}

	glog.Info("Shutting down")
}

type templateBinding struct {
	Username   string
	Attributes cas.UserAttributes
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	}

	if r.URL.Path == "/logout" {
		r.URL.Path = ""
		cas.RedirectToLogout(w, r)
		return
	}

	w.Header().Add("Content-Type", "text/html")

	tmpl, err := template.New("index.html").Parse(index_html)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	binding := &templateBinding{
		Username:   cas.Username(r),
		Attributes: cas.Attributes(r),
	}

	html := new(bytes.Buffer)
	if err := tmpl.Execute(html, binding); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	html.WriteTo(w)
}

const index_html = `<!DOCTYPE html>
<html>
  <head>
    <title>Welcome {{.Username}}</title>
  </head>
  <body>
    <h1>Welcome {{.Username}} <a href="/logout">Logout</a></h1>
    <p>Your attributes are:</p>
    <ul>{{range $key, $values := .Attributes}}
      <li>{{$len := len $values}}{{$key}}:{{if gt $len 1}}
        <ul>{{range $values}}
          <li>{{.}}</li>{{end}}
        </ul>
      {{else}} {{index $values 0}}{{end}}</li>{{end}}
    </ul>
  </body>
</html>
`

const error_500 = `<!DOCTYPE html>
<html>
  <head>
    <title>Error 500</title>
  </head>
  <body>
    <h1>Error 500</h1>
    <p>%v</p>
  </body>
</html>
`
