# CAS Client library

## Example

    mux := http.NewServeMux()
    mux.Handle("/", MyHandler)

    url, _ := url.Parse("https://sso.example.com")
    client := cas.NewClient(&cas.Options{
      URL: url,
    })

    handler := client.Handle(mux)
    server := &http.Server{
      Handler: handler,
    }

    server.ListenAndServe()
