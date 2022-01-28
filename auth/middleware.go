package auth

import (
	"context"
	"fmt"
	"net/http"
)

var middlewareKey string

func WithAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(tokenCookieName)
		if err != nil {
			url := getOAuthConfig(r).AuthCodeURL(r.URL.Path)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		auth, err := validateToken(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "invalid cookie")
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), middlewareKey, auth)))
	})
}

func GetAuth(ctx context.Context) *Auth {
	iface := ctx.Value(middlewareKey)
	if iface == nil {
		return nil
	}
	return iface.(*Auth)
}
