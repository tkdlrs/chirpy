package main

import (
	"encoding/json"
	"net/http"

	"github.com/tkdlrs/chirpy/internal/auth"
	"github.com/tkdlrs/chirpy/internal/database"
)

func (cfg *apiConfig) handlerEditUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		NewPassword string `json:"new_password"`
		Email       string `json:"email"`
	}

	type response struct {
		User
	}
	// Get the access token
	jwtAccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}
	// get user id from jwt
	userId, err := auth.ValidateJWT(jwtAccessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}

	// get parameters from request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}
	// Hash the new password
	hashNewPassword, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash new password", err)
		return
	}
	// Update users password and email in the Database
	updateUserParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashNewPassword,
		ID:             userId,
	}
	user, err := cfg.db.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	//
	updatedUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusOK, response{User: updatedUser})
}
