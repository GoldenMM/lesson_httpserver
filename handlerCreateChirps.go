package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
	"github.com/GoldenMM/lesson_httpserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	log.Printf(`API endpoint [%s] called`, "POST /api/chirps")

	type Chirp struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type ErrResp struct {
		Err string `json:"error"`
	}
	type ValidRes struct {
		CleanedBody string `json:"cleaned_body"`
	}

	chirp := Chirp{}

	// Decode the JSON and check if valid
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
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

	// Check if the user is authenticated
	// Get the token from the header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf(`Error getting token: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Validate the token
	chirp.UserID, err = auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		log.Printf(`Error validating token: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if chirp is less than 140 characters
	if len(chirp.Body) >= 140 {
		log.Printf(`Error with request: Chirp "body" is greater than 140 characters`)
		dat, err := json.Marshal(&ErrResp{Err: "Chirp is too long"})
		if err != nil {
			log.Printf(`Error marshaling JSON: %s`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
		return
	}

	// Clean the chirp
	cleaned_chirp := cleanChirp(chirp.Body)

	// Else it is a valid Chirp
	_, err = json.Marshal(&ValidRes{CleanedBody: cleaned_chirp})
	if err != nil {
		log.Printf(`Error sending valid response`)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the chirp in the database
	chirp_db, err := cfg.dbQueries.CreateChirp(r.Context(),
		database.CreateChirpParams{
			Body:   cleaned_chirp,
			UserID: chirp.UserID,
		})
	if err != nil {
		log.Printf(`Error creating chirp: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chirp_resp := toRespChirp(chirp_db)
	dat, err := json.Marshal(chirp_resp)
	if err != nil {
		log.Printf(`Error marshaling JSON: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
	w.Write(dat)
}

var bannedWords = [...]string{"kerfuffle", "sharbert", "fornax"}

func cleanChirp(chirp string) string {
	// split the chirp into words
	words := strings.Split(chirp, " ")
	var cleaned_words = make([]string, len(words))
	for i, word := range words {
		// normalize the word
		l_word := strings.ToLower(word)
		// check if the word is in the banned words list
		for _, bannedWord := range bannedWords {
			if l_word == bannedWord {
				log.Printf(`Found banned word: %s`, word)
				// replace the word with "****"
				cleaned_words[i] = "****"
				log.Printf(`Message so far: %s`, cleaned_words)
				break
			} else {
				cleaned_words[i] = word
			}
		}
	}
	// join the words back together
	return strings.Join(cleaned_words, " ")
}
