package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// Set up stuff to get info out of request
	type Data struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}
	//
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}
	//
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	//
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user", err)
		return
	}
	//
	_, err = cfg.db.UpgradeUserToRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not upgrade user to Red", err)
		return
	}
	//
	w.WriteHeader(http.StatusNoContent)
}
