package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/creds/{id}", downloadCredentials).Methods("GET")
	router.HandleFunc("/api/creds", listCredentials).Methods("GET")

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
