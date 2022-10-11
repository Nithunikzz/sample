package main

import (
	"fmt"
	"log"
	"net/http"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home page Endpoint")
}

func handleRequests() {
	http.HandleFunc("/", homepage)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func main() {
	handleRequests()
}
