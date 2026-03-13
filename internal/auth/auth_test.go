package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestValidateJWT2(t *testing.T) {
	const tokenSecret = "super-secret"
	userID := uuid.New()

	validToken, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT() expired token error = %v", err)
	}

	tests := []struct {
		name       string
		token      string
		secret     string
		wantUserID uuid.UUID
		wantErr    bool
		errIs      error
	}{
		{
			name:       "Valid token",
			token:      validToken,
			secret:     tokenSecret,
			wantUserID: userID,
			wantErr:    false,
		},
		{
			name:    "Expired token",
			token:   expiredToken,
			secret:  tokenSecret,
			wantErr: true,
			errIs:   jwt.ErrTokenExpired,
		},
		{
			name:    "Wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
			errIs:   jwt.ErrTokenSignatureInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if !errors.Is(err, tt.errIs) {
					t.Fatalf("ValidateJWT() error = %v, want %v", err, tt.errIs)
				}
				return
			}

			if gotUserID != tt.wantUserID {
				t.Fatalf("ValidateJWT() userID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}
