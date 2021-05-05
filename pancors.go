package pancors

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type corsTransport string

func (t corsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Referer", string(t))

	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	res.Header.Set("Access-Control-Allow-Origin", "*")
	res.Header.Set("Access-Control-Allow-Credentials", "true")

	return res, nil
}

// HandleProxy is handler which passes requests through to the host
// and returns their responses with CORS headers
func HandleProxy(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Query().Get("url")

	referer := r.URL.Query().Get("referer")
	if referer == "" {
		referer = r.Header.Get("referer")
	}

	urlParsed, err := url.Parse(urlParam)
	if err != nil || (urlParsed.Scheme != "http" && urlParsed.Scheme != "https") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	proxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = urlParsed
			r.Host = urlParsed.Host
		},
		Transport: corsTransport(referer),
	}

	proxy.ServeHTTP(w, r)
}
