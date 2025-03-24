package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		ExpiresInSec int    `json:"expires_in_seconds"`
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

	// Truncate the ExpiresInSec to 1 hour if it is empty or over 1 hour
	if request.ExpiresInSec == 0 { // If it is empty
		request.ExpiresInSec = 3600
	}
	if request.ExpiresInSec > 3600 { // If it is over 1 hour
		request.ExpiresInSec = 3600
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

	// Password is correct
	// Convert the user to JSON and return
	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Second*time.Duration(request.ExpiresInSec))
	if err != nil {
		log.Printf(`Error making JWT: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userResp := toRespUser(user, token)
	dat, err := json.Marshal(userResp)
	if err != nil {
		log.Printf(`Error marshaling JSON: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
