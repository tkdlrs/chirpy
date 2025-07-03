package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleChirpValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	//
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	//
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	BodyToScrub := params.Body
	splitBody := strings.Split(BodyToScrub, " ")
	scrubbedSlice := make([]string, len(splitBody))
	for i, v := range splitBody {
		lowerCaseV := strings.ToLower(v)
		if lowerCaseV == "kerfuffle" || lowerCaseV == "sharbert" || lowerCaseV == "fornax" {
			scrubbedSlice[i] = "****"
		} else {
			scrubbedSlice[i] = v
		}
	}
	scrubbedBody := strings.Join(scrubbedSlice, " ")

	//
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: scrubbedBody,
	})
}
