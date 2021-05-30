package pancors

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type expected struct {
	statusCode int
	headers    map[string]string
}

type params struct {
	userAgent string
	referer   string
}

type test struct {
	name     string
	url      string
	params   params
	expected expected
}

func TestProxy(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HandleProxy))
	defer testServer.Close()
	testCustomServer := httptest.NewServer(http.HandlerFunc(HandleProxyWith("example.com", "false")))
	defer testCustomServer.Close()

	echoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
		t.Log("one")
	}))
	defer echoServer.Close()
	refererServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Referer") == "https://example.com" {
			fmt.Fprint(w, "Referer is OK")
		} else {
			http.Error(w, "Referer is not 'https://example.com'", http.StatusBadRequest)
		}
	}))
	defer refererServer.Close()

	tests := []test{
		{
			"HTTPS URL with params",
			echoServer.URL,
			params{userAgent: "Test"},
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
			echoServer.URL,
			params{userAgent: "Test"},
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
			params{userAgent: "Test"},
			expected{statusCode: http.StatusBadRequest},
		},
		{
			"Non-HTTP(S) URL",
			"ftp://example.com",
			params{userAgent: "Test"},
			expected{statusCode: http.StatusBadRequest},
		},
		{
			"Referer header",
			refererServer.URL,
			params{userAgent: "Test", referer: "https://example.com"},
			expected{statusCode: http.StatusOK},
		},
		{
			"Wrong Referer header",
			refererServer.URL,
			params{userAgent: "Test", referer: "https://elpmaxe.com"},
			expected{statusCode: http.StatusBadRequest},
		},
		{
			"Missing User-Agent header",
			echoServer.URL,
			params{},
			expected{statusCode: http.StatusBadRequest},
		},
	}

	customTests := []test{
		{
			"Success reponse with specified custom headers",
			echoServer.URL,
			params{userAgent: "Test"},
			expected{
				http.StatusOK,
				map[string]string{
					"Access-Control-Allow-Origin":      "example.com",
					"Access-Control-Allow-Credentials": "false",
				},
			},
		},
	}

	testServerURL, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatalf("Could not parse test server's URL, got %v", err)
	}
	q := testServerURL.Query()
	for _, testCase := range tests {
		t.Run(testCase.name, runTest(testCase, testServerURL, q))
	}

	testServerURL, err = url.Parse(testCustomServer.URL)
	if err != nil {
		t.Fatalf("Could not parse test server's URL, got %v", err)
	}
	q = testServerURL.Query()
	for _, testCase := range customTests {
		t.Run(testCase.name, runTest(testCase, testServerURL, q))
	}
}

func runTest(testCase test, testServerURL *url.URL, q url.Values) func(*testing.T) {
	return func(t *testing.T) {
		q.Set("url", testCase.url)
		q.Set("referer", testCase.params.referer)
		testServerURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(context.Background(), "GET", testServerURL.String(), nil)
		if err != nil {
			t.Fatalf("Could not prepare a request, got %v", err)
		}
		req.Header.Add("User-Agent", testCase.params.userAgent)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Could not fetch testing data")
		}
		defer res.Body.Close()

		if res.StatusCode != testCase.expected.statusCode {
			t.Errorf("Expected HTTP status code %d, got %d", testCase.expected.statusCode, res.StatusCode)
		}

		for header, expected := range testCase.expected.headers {
			actual := res.Header.Get(header)
			if actual != expected {
				t.Errorf("Expected header %s to be %s, got: %v", header, expected, actual)
			}
		}
	}
}
