package auth

import (
	"context"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

func githubClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return (github.NewClient(tc))
}
