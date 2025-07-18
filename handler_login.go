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
	expirationTime := time.Duration(params.ExpiresInSeconds)
	if expirationTime == 0 || time.Duration(expirationTime) > time.Hour {
		expirationTime = time.Hour
	}

	// Generate a JWT for the user
	token, err := auth.MakeJWT(allegedUser.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue generating Token", err)
		return
	}
	// happy path
	user := User{
		ID:        allegedUser.ID,
		UpdatedAt: allegedUser.UpdatedAt,
		CreatedAt: allegedUser.CreatedAt,
		Email:     allegedUser.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, response{User: user})
}
