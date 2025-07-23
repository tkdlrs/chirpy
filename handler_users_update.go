package main

import (
	"encoding/json"
	"net/http"

	"github.com/tkdlrs/chirpy/internal/auth"
	"github.com/tkdlrs/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}
	// Get the access token
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}
	// get user id from JWT
	userId, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
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
	hashNewPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash new password", err)
		return
	}
	// Update users password and email in the Database
	updateUserParams := database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashNewPassword,
	}
	user, err := cfg.db.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not update user", err)
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
