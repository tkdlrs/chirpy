package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerIndexChirps(w http.ResponseWriter, r *http.Request) {
	// get chirps from database
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}
	//
	authorFilter := r.URL.Query().Get("author_id")
	//
	chirps := []Chirp{}
	//
	if authorFilter != "" {
		for _, dbChirp := range dbChirps {
			// filtering for if filtering by author
			if authorFilter == dbChirp.UserID.String() {
				chirps = append(chirps, Chirp{
					ID:        dbChirp.ID,
					CreatedAt: dbChirp.CreatedAt,
					UpdatedAt: dbChirp.UpdatedAt,
					UserID:    dbChirp.UserID,
					Body:      dbChirp.Body,
				})
			}
		}
	} else {
		for _, dbChirp := range dbChirps {
			// filtering for if filtering by author
			if authorFilter == dbChirp.UserID.String() {
				chirps = append(chirps, Chirp{
					ID:        dbChirp.ID,
					CreatedAt: dbChirp.CreatedAt,
					UpdatedAt: dbChirp.UpdatedAt,
					UserID:    dbChirp.UserID,
					Body:      dbChirp.Body,
				})
			}
		}
	}
	//
	respondWithJSON(w, http.StatusOK, chirps)
}
