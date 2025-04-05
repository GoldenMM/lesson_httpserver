package auth

import (
	"net/http/httptest"
	"testing"

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
func TestCheckPasswordHash_Valid(t *testing.T) {
	hash, _ := HashPassword("123456")
	err := CheckPasswordHash("123456", hash)
	if err != nil {
		t.Errorf("CheckPasswordHash() error = %v", err)
	}
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	hash, _ := HashPassword("123456")
	err := CheckPasswordHash("invalidpassword", hash)
	if err == nil {
		t.Errorf("CheckPasswordHash() expected error, got nil")
	}
}

// TestMakeJWT tests the MakeJWT function is working
func TestMakeJWT(t *testing.T) {
	jwt, err := MakeJWT(uuid.New(), "1234")
	if err != nil {
		t.Errorf("MakeJWT() error = %v", err)
	}
	if jwt == "" {
		t.Errorf("MakeJWT() jwt is empty")
	}
}

// TestValidateJWT tests the ValidateJWT function is working with valid token
func TestValidateJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	token, _ := MakeJWT(userID, "1234")
	returnedID, err := ValidateJWT(token, "1234")
	if err != nil {
		t.Errorf("ValidateJWT() error = %v", err)
	}
	if returnedID != userID {
		t.Errorf("ValidateJWT() returnedID = %v, want %v", returnedID, userID)
	}
}

// TestValidateJWTInvalid tests the ValidateJWT function is working with invalid token
func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := ValidateJWT("invalid token", "1234")
	if err == nil {
		t.Errorf("ValidateJWT() error = %v", err)
	}
}

// TestGetBearerToken tests the GetBearerToken function is working
func TestGetBearerToken_ValidAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer validtoken123")
	token, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("GetBearerToken() error = %v", err)
	}
	if token != "validtoken123" {
		t.Errorf("GetBearerToken() token = %v, want %v", token, "validtoken123")
	}
}

func TestGetBearerToken_EmptyAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "")
	_, err := GetBearerToken(req.Header)
	if err == nil {
		t.Errorf("GetBearerToken() expected error, got nil")
	}
}

func TestGetBearerToken_NoAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	_, err := GetBearerToken(req.Header)
	if err == nil {
		t.Errorf("GetBearerToken() expected error, got nil")
	}
}

func TestGetBearerToken_InvalidAuthorizationHeaderFormat(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidHeaderFormat")
	_, err := GetBearerToken(req.Header)
	if err == nil {
		t.Errorf("GetBearerToken() expected error, got nil")
	}
}

func TestMakeRefreshToken(t *testing.T) {
	refreshToken, err := MakeRefreshToken()
	if err != nil {
		t.Errorf("MakeRefreshToken() error = %v", err)
	}
	if refreshToken == "" {
		t.Errorf("MakeRefreshToken() refreshToken is empty")
	}
}
