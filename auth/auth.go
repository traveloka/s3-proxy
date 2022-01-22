package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// Auth contains authentication info.
type Auth struct {
	claim jwtClaim
}

type jwtClaim struct {
	jwt.StandardClaims
	GithubToken string `json:"github:token"`
	Name        string `json:"profile:name"`
}
