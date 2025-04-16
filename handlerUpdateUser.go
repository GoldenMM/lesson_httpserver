package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
	"github.com/GoldenMM/lesson_httpserver/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	log.Printf(`API endpoint [%s] called`, "PUT /api/users")

	//Get and check the token from the header
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf(`Error getting token from header: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the token is valid
	userID, err := auth.ValidateJWT(jwtToken, cfg.tokenSecret)
	if err != nil {
		log.Printf(`Error checking token: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Decode the JSON and check if valid
	decoder := json.NewDecoder(r.Body)
	request := Request{}
	err = decoder.Decode(&request)
	if err != nil {
		log.Printf(`Error decoding JSON: %s`, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf(`Error hashing password: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the user in the database
	user, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          request.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf(`Error updating user: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResp := UserResp{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		Token:       jwtToken,
		IsChirpyRed: user.IsChirpyRed,
	}

	// Marshal the user to JSON
	dat, err := json.Marshal(&userResp)
	if err != nil {
		log.Printf(`Error marshaling JSON: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}
