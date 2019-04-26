package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Version represents a semver version
// NORELEASE => no version
// PATCH     => 0.0.1
// MINOR     => 0.1.0
// MAJOR     => 1.0.0
// ALPHA     => 1.0.0-alpha
// BETA      => 1.0.0-beta
// RC        => 1.0.0-rc
type Version int

//go:generate stringer -type=Version
const (
	// UNVALID_VERSION is default when no semver version exists
	UNVALID_VERSION Version = iota
	// NORELEASE represents for instance 1.0.0-rc
	NORELEASE
	// ALPHA represents for instance 1.0.0-alpha
	ALPHA
	// BETA represents for instance 1.0.0-beta
	BETA
	// RC represents for instance 1.0.0-rc
	RC
	// PATCH represents for instance 0.0.1
	PATCH
	// MINOR represents for instance 0.1.0
	MINOR
	// MAJOR represents for instance 1.0.0
	MAJOR
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
