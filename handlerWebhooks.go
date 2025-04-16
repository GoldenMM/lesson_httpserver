package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
	"github.com/GoldenMM/lesson_httpserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// Handle the Polka webhook

	// Check the api key in the header
	headerKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf(`Error getting API key from header: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if headerKey != cfg.polkaKey {
		log.Printf(`Invalid API key: %s`, headerKey)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Read the request body
	type WebhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var webhookRequest WebhookRequest

	// Decode the JSON
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&webhookRequest)
	if err != nil {
		log.Printf(`Error decoding JSON: %s`, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check the event type we only care about "user.upgraded"
	if webhookRequest.Event != "user.upgraded" {
		log.Printf(`Unsupported event type: %s`, webhookRequest.Event)
		w.WriteHeader(http.StatusNoContent) //204
		return
	} else { // Handle the event and upgrade the user in the database
		log.Printf(`Webhook event: %s`, webhookRequest.Event)
		// Parse the user ID into a UUID
		userID, err := uuid.Parse(webhookRequest.Data.UserID)
		if err != nil {
			log.Printf(`Error parsing user ID: %s`, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Update the user in the database
		err = cfg.dbQueries.UpgradeUserChirpyRed(r.Context(),
			database.UpgradeUserChirpyRedParams{
				ID:          userID,
				IsChirpyRed: true,
			})
		if err != nil {
			log.Printf(`Error upgrading user: %s`, err)
			w.WriteHeader(http.StatusNotFound) //404
			return
		}
		// Else log the sucess and return 204
		log.Printf(`User %s upgraded to Chirpy Red`, userID)
		w.WriteHeader(http.StatusNoContent) //204
	}

}
