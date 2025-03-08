package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type Email struct {
		Email string `json:"email"`
	}

	type ErrResp struct {
		Err string `json:"error"`
	}
	// Decode the JSON and check if valid
	decoder := json.NewDecoder(r.Body)
	email := Email{}

	err := decoder.Decode(&email)
	// If something went wrong
	if err != nil {
		log.Printf(`Error decoding JSON: %s`, err)
		dat, err := json.Marshal(&ErrResp{Err: "Something went wrong"})
		if err != nil {
			log.Printf(`Error marshaling JSON: %s`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
		return
	}

	// Else it is a valid User
	user, err := cfg.dbQueries.CreateUser(r.Context(), email.Email)
	if err != nil {
		log.Printf(`Error creating user: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		dat, err := json.Marshal(toRespUser(user))
		if err != nil {
			log.Printf(`Error marshaling JSON: %s`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(dat)
	}

}
