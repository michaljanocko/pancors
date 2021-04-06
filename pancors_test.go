package pancors

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	go func() {
		http.HandleFunc("/", HandleProxy)

		fmt.Println("Initialized PanCORS")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	url := "http://localhost:8080/?url=https%3A%2F%2Fsuggest.seznam.cz%2Fslovnik%2Fmix_cz_en%3Fphrase%3Dtest%26format%3Djson-2"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Couldn't create testing request")
	}

	req.Header.Add("origin", "test")

	res, _ := http.DefaultClient.Do(req)
	if err != nil {
		panic("Couldn't fetch testing data")
	}

	defer res.Body.Close()

	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		panic("Access-Control-Allow-Origin is not '*'")
	}
	if res.Header.Get("Access-Control-Allow-Credentials") != "true" {
		panic("Access-Control-Allow-Credentials is not 'true'")
	}
}
