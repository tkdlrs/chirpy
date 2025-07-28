package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerIndexChirps(w http.ResponseWriter, r *http.Request) {
	// get chirps from database
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}
	//
	sortDirection := "asc"
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "desc" {
		sortDirection = sortParam
	}
	//
	authorID := uuid.Nil
	authorIDStringFilter := r.URL.Query().Get("author_id")
	if authorIDStringFilter != "" {
		authorID, err = uuid.Parse(authorIDStringFilter)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}
	//
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		// filtering for if filtering by author
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}
		//
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}
	// apply sorting
	if sortDirection == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			first := chirps[i].CreatedAt
			second := chirps[j].CreatedAt
			//
			return first.After(second)
		})
	}
	//
	respondWithJSON(w, http.StatusOK, chirps)
}
