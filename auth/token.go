package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtAudience = os.Getenv("OAUTH_AUDIENCE")

func validateToken(tokenString string) (*Auth, error) {
	secretKey := os.Getenv("OAUTH_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("can not parse auth token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claim := token.Claims.(*jwtClaim)
	err = claim.Valid()
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !claim.VerifyAudience(jwtAudience, true) {
		return nil, fmt.Errorf("invalid audience")
	}

	return &Auth{*claim}, nil
}

func createToken(claim *jwtClaim) (string, error) {
	secretKey := os.Getenv("OAUTH_SECRET")
	claim.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
	claim.Audience = jwtAudience
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secretKey))
}
