package auth

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HandleAuthCallback read github auth code and generate auth token
func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authToken, err := generateAuthToken(ctx, getOAuthConfig(r), r.URL.Query().Get("code"))
	if err != nil {
		log.WithError(err).Error("Failed to authenticate user")
		http.Error(w, err.Error(), 400)
		return
	}
	cookie := createCookie(tokenCookieName, *authToken)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, r.URL.Query().Get("state"), http.StatusSeeOther)
}

// HandleLogout remove the auth token
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie := createCookie(tokenCookieName, "")
	http.SetCookie(w, cookie)
	w.Header().Set("content-type", "text/plain")
	fmt.Fprint(w, "Logged out")
}
