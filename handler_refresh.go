package main

import (
	"net/http"
	"time"

	"github.com/tkdlrs/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	// Get the refresh token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find token", err)
		return
	}
	// Look up the user based on their refresh token
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get user for refresh token", err)
		return
	}
	// Make a new JWT
	newAccessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}
	//
	tokenResponse := response{
		Token: newAccessToken,
	}
	respondWithJSON(w, http.StatusOK, tokenResponse)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}
	// revoke the refresh token
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not revoke session", err)
		return
	}
	//
	w.WriteHeader(http.StatusNoContent)
}
