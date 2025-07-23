package main

import (
	"net/http"

	"github.com/tkdlrs/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	type response struct {
	}
	// Get the refresh token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get refresh token", err)
		return
	}
	// Look up the Refresh Token
	refreshTokenInfo, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to find refresh token", err)
		return
	}
	// revoke the refresh token
	cfg.db.RevokeRefreshToken(r.Context(), refreshTokenInfo.Token)
	//
	respondWithJSON(w, http.StatusNoContent, response{})
}
