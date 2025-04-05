package auth

import (
	"errors"
	"fmt"

	"github.com/CP-Payne/taskflow/user/config/authkeys"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	keys authkeys.AuthKeys
	aud  string
	iss  string
}

func NewJWTAuthenticator(keys authkeys.AuthKeys, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		keys: keys,
		aud:  aud,
		iss:  iss,
	}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	prvKey := a.keys.PrivateKey()
	if prvKey == nil {
		return "", errors.New("failed to load private key")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return "", fmt.Errorf("GenerateToken: parse key: %w", err)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, err
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	pubKey := a.keys.PublicKey()
	if pubKey == nil {
		return nil, errors.New("failed to load public key")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, fmt.Errorf("ValidateToken: parse key: %w", err)
	}

	return jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return publicKey, nil
	})

	// Additional checks can be performed here to check the Issuer and Audience
}
