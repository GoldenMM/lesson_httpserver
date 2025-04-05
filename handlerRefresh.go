package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	log.Printf(`API endpoint [%s] called`, "POST /api/refresh")

	// No body but the token is in the header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf(`Error getting token from header: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if token == "" {
		log.Printf(`Error getting token from header`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// DEBUG:
	log.Printf(`Target Refresh Token: %s`, token)

	// Get all the refresh tokens
	refreshTokens, err := cfg.dbQueries.GetAllRefreshTokens(r.Context())
	if err != nil {
		log.Printf(`Error getting all refresh tokens: %s`, err)
		return
	}
	for i, rt := range refreshTokens {
		log.Printf(`Stored Refresh Token [%v]: %s`, i, rt.Token)
	}

	// Look up refresh token in database
	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf(`Error getting refresh token: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the refresh token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		log.Printf(`Error refresh token is expired`)
		// Set token to revoked
		err = cfg.dbQueries.RevokeRefreshToken(r.Context(), token)
		if err != nil {
			log.Printf(`Error revoking refresh token: %s`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Check if the refresh token is revoked
	isRevoked := refreshToken.RevokedAt != sql.NullTime{}
	if isRevoked {
		log.Printf(`Error refresh token is revoked`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret)
	if err != nil {
		log.Printf(`Error generating access token: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Marshal the access token to JSON
	type Resp struct {
		Token string `json:"token"`
	}
	dat, err := json.Marshal(&Resp{Token: newToken})
	if err != nil {
		log.Printf(`Error marshaling JSON: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}
