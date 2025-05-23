package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	GenerateToken(jwt jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
