package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	//
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	//
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	//
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handleChirpValidation)
	//
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	//
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	//
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handleChirpValidation(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON of the chirp
	type chirp struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding chirp body: %s", err)
		w.WriteHeader(500)
		return
	}
	//
	type returnErrorVal struct {
		Error string `json:"error"`
	}
	type returnValidVal struct {
		Valid bool `json:"valid"`
	}
	// check at it's valid and encode a response
	if len(params.Body) > 140 {
		respBody := returnErrorVal{
			Error: "Chirp is too Long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			respBody.Error = "Something went wrong"
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	//
	respBody := returnValidVal{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		respBody.Valid = false
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
