package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerShowChirp(w http.ResponseWriter, r *http.Request) {
	//
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Chirp not found"))
		return
	}
	getChirpParams := uuid.MustParse(chirpID)
	//
	dbChirp, err := cfg.db.GetChirp(r.Context(), getChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirp", err)
		return
	}
	// Happy path
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
