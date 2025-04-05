package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	c := jwt.RegisteredClaims{
		Issuer:    "chripy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) { // Function to return the key for validation
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(tokenSecret), nil
		})
	if err != nil {
		log.Println("Error validating JWT token:", err)
		return uuid.Nil, err
	}

	// Check if token is valid
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	stringUserID, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("Error extracting ID from token:", err)
		return uuid.Nil, err
	}

	return uuid.MustParse(stringUserID), nil
}

func GetBearerToken(headers http.Header) (string, error) {
	uncleanToken := headers.Get("Authorization")

	if uncleanToken == "" {
		return "", errors.New("no token provided")
	}
	if uncleanToken == "Bearer " {
		return "", errors.New("empty token provided")
	}
	if uncleanToken[:7] != "Bearer " {
		return "", errors.New("invalid token format")
	}
	return uncleanToken[7:], nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b), nil
}
