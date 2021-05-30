package main

import (
	"log"
	"net/http"
	"os"

	"github.com/michaljanocko/pancors"
)

func getAllowOrigin() string {
	if origin, ok := os.LookupEnv("ALLOW_ORIGIN"); ok {
		return origin
	}
	return "*"
}

func getAllowCredentials() string {
	if credentials, ok := os.LookupEnv("ALLOW_CREDENTIALS"); ok {
		return credentials
	}
	return "true"
}

func getListenPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return ":8080"
}

func main() {
	http.HandleFunc("/", pancors.HandleProxyWith(getAllowOrigin(), getAllowCredentials()))

	port := getListenPort()
	log.Printf("PanCORS started listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
