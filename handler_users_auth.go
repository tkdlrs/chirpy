package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tkdlrs/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// compare alleged user hashed password with hashed password on file
	fmt.Println("------- provided password ")
	fmt.Println(params.Password)
	fmt.Println("------- recorded password ")
	fmt.Println(allegedUser.HashedPassword)
	//
	err = auth.CheckPasswordHash(params.Password, allegedUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// happy path
	user := User{
		ID:        allegedUser.ID,
		CreatedAt: allegedUser.CreatedAt,
		UpdatedAt: allegedUser.UpdatedAt,
		Email:     allegedUser.Email,
	}
	respondWithJSON(w, http.StatusOK, user)
}
