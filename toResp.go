package main

import (
	"log"
	"time"

	"github.com/GoldenMM/lesson_httpserver/internal/database"
	"github.com/google/uuid"
)

type UserResp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func toRespUser(u database.User, token ...string) UserResp {
	resp := UserResp{
		Id:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
	if len(token) == 1 {
		resp.Token = token[0]
	} else if len(token) > 1 {
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
