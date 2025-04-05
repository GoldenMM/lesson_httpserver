package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/GoldenMM/lesson_httpserver/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

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

	// Look up refresh token in database
	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf(`Error getting refresh token: %s`, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the refresh token is revoked
	isRevoked := refreshToken.RevokedAt != sql.NullTime{}
	if isRevoked {
		log.Printf(`Error refresh token is already revoked`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf(`Error revoking refresh token: %s`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
