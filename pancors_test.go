package pancors

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProxy(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandleProxy))
	defer testServer.Close()

	refererServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Referer") == "https://example.com" {
			fmt.Fprint(w, "Referer is OK")
		} else {
			http.Error(w, "Referer is not 'https://example.com'", http.StatusBadRequest)
		}
	}))
	defer refererServer.Close()

	type expected struct {
		statusCode int
		headers    map[string]string
	}

	tests := []struct {
		name     string
		url      string
		referer  string
		expected expected
	}{
		{
			"HTTPS URL with params",
			"https://suggest.seznam.cz/slovnik/mix_cz_en?phrase=test&format=json-2",
			"",
			expected{
				http.StatusOK,
				map[string]string{
					"Access-Control-Allow-Origin":      "*",
					"Access-Control-Allow-Credentials": "true",
				},
			},
		},
		{
			"HTTP URL with params",
			"http://suggest.seznam.cz/slovnik/mix_cz_en?phrase=test&format=json-2",
			"",
			expected{
				http.StatusOK,
				map[string]string{
					"Access-Control-Allow-Origin":      "*",
					"Access-Control-Allow-Credentials": "true",
				},
			},
		},
		{
			"Empty URL",
			"",
			"",
			expected{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			"Non-HTTP(S) URL",
			"ftp://example.com",
			"",
			expected{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			"Referer header",
			refererServer.URL,
			"https://example.com",
			expected{
				statusCode: http.StatusOK,
			},
		},
		{
			"Wrong Referer header",
			refererServer.URL,
			"http://google.com",
			expected{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	testServerURL, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatalf("Could not parse test server's URL; got %v", err)
	}
	q := testServerURL.Query()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q.Set("url", tc.url)
			q.Set("referer", tc.referer)
			testServerURL.RawQuery = q.Encode()

			req, err := http.NewRequestWithContext(context.Background(), "GET", testServerURL.String(), nil)
			if err != nil {
				t.Fatalf("Could not prepare a request; got %v", err)
			}

			rsp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Could not fetch testing data")
			}
			defer rsp.Body.Close()

			if rsp.StatusCode != tc.expected.statusCode {
				t.Errorf("Expected HTTP status code %d; got %d", tc.expected.statusCode, rsp.StatusCode)
			}

			for header, expected := range tc.expected.headers {
				actual := rsp.Header.Get(header)
				if actual != expected {
					t.Errorf("Expected header %s = %s; got: %v", header, expected, actual)
				}
			}
		})
	}
}
