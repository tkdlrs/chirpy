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
		respondWithError(w, http.StatusUnauthorized, "Unable to get refresh token", err)
		return
	}
	// Look up the Refresh Token
	refreshTokenInfo, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to find refresh token", err)
		return
	}
	// Check that it has not been revoked.
	if refreshTokenInfo.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Attempting to use a revoked refresh token", err)
		return
	}
	// Check that is not expired
	expiredCheck := time.Now().Compare(refreshTokenInfo.ExpiresAt)
	if expiredCheck > 0 {
		// is expired
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired", err)
		return
	}
	// Make a new JWT
	newAccessToken, err := auth.MakeJWT(refreshTokenInfo.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate access JWT", err)
		return
	}
	tokenResponse := response{
		Token: newAccessToken,
	}
	//
	respondWithJSON(w, http.StatusOK, tokenResponse)
}
