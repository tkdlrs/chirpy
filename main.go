package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	//
	mux := http.NewServeMux()
	//
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	//
	log.Printf("Servering on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
