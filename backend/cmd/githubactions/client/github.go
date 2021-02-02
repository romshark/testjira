package client

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

// GetGithubClient returns an authenticated github client
func GetGithubClientFromEnv(ctx context.Context) *github.Client {
	// get github token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN must not be empty")
	}

	return NewGithubClient(ctx, token)
}

func NewGithubClient(ctx context.Context, token string) *github.Client {
	// authenticate github client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
