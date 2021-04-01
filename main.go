package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func getListenPort() string {
	port, isSet := os.LookupEnv("PORT")

	if !isSet {
		return ":8080"
	}

	return ":" + port
}

func invalidUrl(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Invalid URL")
}

type CorsTransport http.Header

func (t CorsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(r)

	if err != nil {
		return nil, err
	}

	res.Header.Set("Access-Control-Allow-Origin", "*")
	res.Header.Set("Access-Control-Allow-Credentials", "true")

	return res, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")

		urlParsed, err := url.Parse(urlParam)
		if err != nil || (urlParsed.Scheme != "http" && urlParsed.Scheme != "https") {
			invalidUrl(w)
			return
		}

		proxy := httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL = urlParsed
				r.Host = urlParsed.Host
			},
			Transport: CorsTransport(http.Header{}),
		}
		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Initialized PanCORS")
	log.Fatal(http.ListenAndServe(getListenPort(), nil))
}
