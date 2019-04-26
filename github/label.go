package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type SemverLabelService struct {
	owner      string
	repository string
	client     *github.Client
}

// NewSemverLabelService creates a new instance of SemverLabelService
func NewSemverLabelService(owner string, repository string, token string) SemverLabelService {
	return SemverLabelService{
		owner,
		repository,
		github.NewClient(
			oauth2.NewClient(
				context.Background(),
				oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
			),
		),
	}
}
