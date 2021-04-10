package main

import (
	"log"
	"net/http"
	"os"

	"github.com/michaljanocko/pancors"
)

func getListenPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return ":8080"
}

func main() {
	http.HandleFunc("/", pancors.HandleProxy)

	port := getListenPort()
	log.Printf("PanCORS started listening at %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
