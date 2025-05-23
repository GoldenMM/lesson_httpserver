package main

import (
	"log"
	"time"

	"github.com/GoldenMM/lesson_httpserver/internal/database"
	"github.com/google/uuid"
)

type UserResp struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func toRespUser(u database.User, tokens ...string) UserResp {
	resp := UserResp{
		Id:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}
	if len(tokens) >= 1 {
		resp.Token = tokens[0]
	}
	if len(tokens) == 2 {
		resp.RefreshToken = tokens[1]
	}
	if len(tokens) > 2 {
		log.Fatalln("Too many arguments for token")
	}
	return resp
}

type ChirpResp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func toRespChirp(c database.Chirp) ChirpResp {
	return ChirpResp{
		Id:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}
