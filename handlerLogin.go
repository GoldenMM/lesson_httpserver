package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
	"github.com/GoldenMM/lesson_httpserver/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	log.Printf(`API endpoint [%s] called`, "POST /api/login")

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type ErrResp struct {
		Err string `json:"error"`
	}

	// Decode the JSON and check if valid
	decoder := json.NewDecoder(r.Body)
	request := Request{}
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf(`Error decoding JSON: %s`, err)
		dat, err := json.Marshal(&ErrResp{Err: "Something went wrong"})
		if err != nil {
			log.Printf(`Error marshaling JSON: %s`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	// Else it is a valid User so check in database
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		log.Printf(`Error getting user by email: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the password is correct
	err = auth.CheckPasswordHash(request.Password, user.HashedPassword)
	if err != nil {
		log.Printf(`Error checking password: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		log.Printf(`Error making JWT: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Make the refresh token
	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf(`Error making refresh token: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Save the refresh token to the database
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refresh_token,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60), // 60 days
			UserID:    user.ID,
		})
	if err != nil {
		log.Printf(`Error creating refresh token: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Load the userResp in preperation to marshal
	userResp := toRespUser(user, token, refresh_token)
	log.Printf(`JWTToken: %s`, userResp.Token)
	log.Printf(`Created RefreshToken: %s`, userResp.RefreshToken)
	dat, err := json.Marshal(userResp)
	if err != nil {
		log.Printf(`Error marshaling JSON: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We successfully logged in
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
