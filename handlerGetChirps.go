package main

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	// Query the database for all chirps
	db_chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		log.Printf(`Error getting all chirps: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the chirps to response structure
	response_chirps := make([]ChirpResp, 0)
	for _, db_chirp := range db_chirps {
		response_chirps = append(response_chirps, toRespChirp(db_chirp))
	}

	// Convert the response structure to JSON
	response_json, err := json.Marshal(response_chirps)
	if err != nil {
		log.Printf("Failed to marshal chirps to JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	// Get the chirp ID
	chirp_Id := r.PathValue("chirpID")
	uuid_chirp_Id, err := uuid.Parse(chirp_Id)
	if err != nil {
		log.Printf(`Error parsing UUID: %s`, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Query the database for the chirp
	db_chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), uuid_chirp_Id)
	if err != nil {
		log.Printf(`Error getting chirp: %s`, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal the chirp to JSON
	dat, err := json.Marshal(toRespChirp(db_chirp))
	if err != nil {
		log.Printf(`Error marshaling JSON: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
