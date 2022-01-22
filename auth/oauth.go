package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
)

func getOAuthConfig(r *http.Request) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),

		Endpoint: oauth2.Endpoint{
			AuthURL:  githubAuthorizeURL,
			TokenURL: githubTokenURL,
		},
		Scopes: []string{},
	}
}

func generateAuthToken(ctx context.Context, cfg *oauth2.Config, code string) (*string, error) {
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("can not exchange auth code: %w", err)
	}

	ghClient := githubClient(ctx, token.AccessToken)
	user, _, err := ghClient.Users.Get(ctx, "")
	if err != nil {
		log.WithError(err).Error("Can not get user info from github")
		return nil, fmt.Errorf("can not get user info from github: %w", err)
	}

	orgMember, _, err := ghClient.Organizations.GetOrgMembership(ctx, user.GetLogin(), "traveloka")
	if err != nil {
		log.WithField("user", user.GetLogin()).WithError(err).Error("Can not get organization member status")
		return nil, fmt.Errorf("error when retrieving organization member status: %w", err)
	}
	if orgMember == nil {
		return nil, fmt.Errorf("user %s is not member of org %s", user.GetLogin(), "traveloka")
	}

	log.WithField("githubLogin", user.GetLogin()).WithField("name", user.GetName()).Info("User autheticated")
	claim := &jwtClaim{
		GithubToken: token.AccessToken,
		Name:        user.GetName(),
	}
	claim.Subject = "github:" + user.GetLogin()
	tokenString, err := createToken(claim)
	return &tokenString, err
}
