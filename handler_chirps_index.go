package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerIndexChirps(w http.ResponseWriter, r *http.Request) {
	//
	// type response struct {
	// 	Chirps []Chirp
	// }
	// get chirps from database
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}
	//
	formatChirps := make([]Chirp, len(chirps))

	for index, value := range chirps {
		formatChirps[index] = Chirp{
			ID:        value.ID,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
			UserID:    value.UserID,
			Body:      value.Body,
		}
	}
	// Happy path
	respondWithJSON(w, http.StatusOK, formatChirps)

}
