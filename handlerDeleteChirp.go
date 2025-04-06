package main

import (
	"log"
	"net/http"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	log.Printf(`API endpoint [%s] called`, "DELETE /api/chirps")

	// Get the chirp ID
	chirp_Id := r.PathValue("chirpID")
	uuid_chirp_Id, err := uuid.Parse(chirp_Id)
	if err != nil {
		log.Printf(`Error parsing UUID: %s`, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the token from the header
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf(`Error getting token from header: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the Chirp from the database
	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), uuid_chirp_Id)
	if err != nil {
		log.Printf(`Error getting chirp by ID: %s`, err)
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Check if the chirp belongs to the user
	userID, err := auth.ValidateJWT(jwtToken, cfg.tokenSecret)
	if err != nil {
		log.Printf(`Error checking token: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if chirp.UserID != userID {
		log.Printf(`Error chirp does not belong to user`)
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	// Delete the chirp from the database
	err = cfg.dbQueries.DeleteChirp(r.Context(), uuid_chirp_Id)
	if err != nil {
		log.Printf(`Error deleting chirp: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204

}
