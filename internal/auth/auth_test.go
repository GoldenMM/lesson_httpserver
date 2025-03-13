package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestHashPassword tests the HashPassword function is working
func TestHashPassword(t *testing.T) {
	_, err := HashPassword("password")
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
	}
}

// TestCheckPasswordHash tests the CheckPasswordHash function is and matched created hash
func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword("123456")
	err := CheckPasswordHash("123456", hash)
	if err != nil {
		t.Errorf("CheckPasswordHash() error = %v", err)
	}
}

// TestMakeJWT tests the MakeJWT function is working
func TestMakeJWT(t *testing.T) {
	jwt, err := MakeJWT(uuid.New(), "1234", time.Hour)
	if err != nil {
		t.Errorf("MakeJWT() error = %v", err)
	}
	if jwt == "" {
		t.Errorf("MakeJWT() jwt is empty")
	}
}

// TestValidateJWT tests the ValidateJWT function is working with valid token
func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	token, _ := MakeJWT(userID, "1234", time.Hour)
	returnedID, err := ValidateJWT(token, "1234")
	if err != nil {
		t.Errorf("ValidateJWT() error = %v", err)
	}
	if returnedID != userID {
		t.Errorf("ValidateJWT() returnedID = %v, want %v", returnedID, userID)
	}
}

// TestValidateJWTInvalid tests the ValidateJWT function is working with invalid token
func TestValidateJWTInvalid(t *testing.T) {
	_, err := ValidateJWT("invalid token", "1234")
	if err == nil {
		t.Errorf("ValidateJWT() error = %v", err)
	}
}

// TestValidateJWTExpired tests the ValidateJWT function is working with expired token
func TestValidateJWTExpired(t *testing.T) {
	userID := uuid.New()
	token, _ := MakeJWT(userID, "1234", -time.Hour)
	_, err := ValidateJWT(token, "1234")
	if err == nil {
		t.Errorf("ValidateJWT() error = %v", err)
	}
}
