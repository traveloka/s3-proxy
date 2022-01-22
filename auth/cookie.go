package auth

import (
	"net/http"
	"os"
)

const tokenCookieName = "s3-proxy-auth"

func createCookie(name, value string) *http.Cookie {
	maxAge := 3600
	if value == "" {
		maxAge = -1
	}
	secure := os.Getenv("GO_ENV") == "production"
	var sameSite http.SameSite
	if secure {
		sameSite = http.SameSiteNoneMode
	}
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   secure,
		SameSite: sameSite,
	}
}
