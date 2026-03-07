package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./")))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Starting server on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
