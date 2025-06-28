package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileserverHits.Add(1)
	return next
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	//
	mux := http.NewServeMux()
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	//
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerServerHits)
	// mux.HandleFunc("reset", apiCfg. )
	//
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	//
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerServerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// io.WriteString(w, "Hits: ")
	fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits.Load())
}

// func (cfg *apiConfig) getServerHits() int32 {
// 	return int32(cfg.fileserverHits.Load())
// }

// func handlerResetServerHits(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	// resetSeverHit()
// 	fmt.Fprintf(w, "Server counts reset")
// }

// func resetSeverHits(cfg apiConfig) {
// }
