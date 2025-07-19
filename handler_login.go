package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tkdlrs/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	// get parameters from request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}
	// get alleged user from database
	allegedUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	//
	err = auth.CheckPasswordHash(params.Password, allegedUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// expiration time should be what they specifed up to an hour. Otherwise it should expire after one hour.
	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	// Generate a JWT for the user
	accessToken, err := auth.MakeJWT(allegedUser.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate access JWT", err)
		return
	}
	// happy path
	userResponse := response{
		User: User{
			ID:        allegedUser.ID,
			CreatedAt: allegedUser.CreatedAt,
			UpdatedAt: allegedUser.UpdatedAt,
			Email:     allegedUser.Email,
		},
		Token: accessToken,
	}
	respondWithJSON(w, http.StatusOK, userResponse)
}
