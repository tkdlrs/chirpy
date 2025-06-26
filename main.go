package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "/app"
	const port = "8080"
	//
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	//
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	//
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
