package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/michaljanocko/pancors"
)

func getListenPort() string {
	port, isSet := os.LookupEnv("PORT")

	if !isSet {
		return ":8080"
	}

	return ":" + port
}

func main() {
	http.HandleFunc("/", pancors.HandleProxy)

	fmt.Println("Initialized PanCORS")
	log.Fatal(http.ListenAndServe(getListenPort(), nil))
}
