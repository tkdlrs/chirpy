package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Check for correct enivornment
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Incorrect platform", nil)
	}
	// Delete all users from the database.
	err := cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to delete all users", err)
		return
	}
	//
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and all users removed from Postgres Database.\n"))
}
