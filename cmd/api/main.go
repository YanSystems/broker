package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRequests)
	log.Println("Broker is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
