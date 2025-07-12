package main

import (
	"encoding/json"
	"net/http"

	"github.com/tkdlrs/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
	// happy path
	user := User{
		ID:        allegedUser.ID,
		Email:     allegedUser.Email,
		CreatedAt: allegedUser.CreatedAt,
		UpdatedAt: allegedUser.UpdatedAt,
	}
	respondWithJSON(w, http.StatusOK, response{User: user})
}
