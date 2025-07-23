package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tkdlrs/chirpy/internal/auth"
	"github.com/tkdlrs/chirpy/internal/database"
)

func (cfg *apiConfig) handlerAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
	// Generate a JWT for the user
	accessToken, err := auth.MakeJWT(
		allegedUser.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate access JWT", err)
		return
	}
	// Generate a Refresh Token for the User
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate refresh token", err)
		return
	}
	// Save Refresh Token in Database
	saveRefreshTokenParams := database.CreateRefreshTokenParams{
		UserID:    allegedUser.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), saveRefreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save refresh token", err)
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
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, userResponse)
}
