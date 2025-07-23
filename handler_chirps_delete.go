package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tkdlrs/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// Check that user is authorized to do this
	// Get the access token
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}
	// get user id from JWT
	userId, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}
	//
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	//
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp", err)
		return
	}
	// Check that the chirps original owner is the varified user and if so allow the delete
	if dbChirp.UserID == userId {
		cfg.db.DeleteChirp(r.Context(), dbChirp.ID)
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		respondWithError(w, http.StatusForbidden, "Action not permitted", err)
		return
	}

}
