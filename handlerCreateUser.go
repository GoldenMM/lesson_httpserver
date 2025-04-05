package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
	"github.com/GoldenMM/lesson_httpserver/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	log.Printf(`API endpoint [%s] called`, "POST /api/users")

	type NewUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type ErrResp struct {
		Err string `json:"error"`
	}
	// Decode the JSON and check if valid
	decoder := json.NewDecoder(r.Body)
	newUser := NewUser{}

	err := decoder.Decode(&newUser)
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

	// Check if email or password is empty
	if newUser.Email == "" || newUser.Password == "" {
		log.Printf(`Error with request: Email or Password is empty`)
		_, err := json.Marshal(&ErrResp{Err: "Email or Password is empty"})
		if err != nil {
			log.Printf(`Error marshaling JSON: %s`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Else it is a valid User
	// Create the hashed password for the user
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		log.Printf(`Error hashing password: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	createUserParams := database.CreateUserParams{
		Email:          newUser.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), createUserParams)
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
