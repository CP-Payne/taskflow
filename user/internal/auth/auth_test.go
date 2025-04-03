package auth_test

import (
	"testing"
	"time"

	"github.com/CP-Payne/taskflow/user/config"
	"github.com/CP-Payne/taskflow/user/config/authkeys"
	"github.com/CP-Payne/taskflow/user/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_Success(t *testing.T) {
	cfg := config.New("../../config/.env")

	if cfg.KeyPath == "" {
		t.Errorf("failed to load configuration")
	}

	authKeys := authkeys.NewAuthKeys()
	authKeys.Load(cfg.KeyPath)

	authenticator := auth.NewJWTAuthenticator(authKeys, "user", "user-service")

	claims := make(jwt.MapClaims)
	claims["iss"] = "test-iss"
	claims["aud"] = "test-aud"
	claims["sub"] = int64(1)
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		t.Errorf("GenerateToken expected nil error, received error: %v ", err)
	}

	if token == "" {
		t.Error("Expected signed JWT token, received empty string")
	}
}

func TestValidateToken(t *testing.T) {
	cfg := config.New("../../config/.env")

	if cfg.KeyPath == "" {
		t.Errorf("failed to load configuration")
	}

	authKeys := authkeys.NewAuthKeys()
	authKeys.Load(cfg.KeyPath)

	authenticator := auth.NewJWTAuthenticator(authKeys, "user", "user-service")

	claims := make(jwt.MapClaims)
	claims["iss"] = "test-iss"
	claims["aud"] = "test-aud"
	claims["sub"] = int64(1)
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		t.Errorf("GenerateToken expected nil error, received error: %v ", err)
	}

	if token == "" {
		t.Error("Expected signed JWT token, received empty string")
	}

	_, err = authenticator.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken: Expected nil, received: %v", err)
	}

	invalidToken := token + "invalid"

	_, err = authenticator.ValidateToken(invalidToken)
	if err == nil {
		t.Errorf("ValidateToken: Expected err, received %v", err)
	}
}
