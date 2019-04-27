package github

import (
	"context"
	"fmt"

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

// SemverLabelService deals with pull requests
// to manage semver labels
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

// GetFromPullRequest find out semver version label attached to a pull request, if there is none or more than one
// this function returns an error
func (s SemverLabelService) GetFromPullRequest(pullRequestNumber int) (Version, error) {
	labels, _, err := s.client.Issues.ListLabelsByIssue(context.Background(), s.owner, s.repository, pullRequestNumber, nil)
	if err != nil {
		return UNVALID_VERSION, fmt.Errorf("can't fetch github api to get label for pull request #%d : %s", pullRequestNumber, err)
	}

	ls := []github.Label{}

	for _, l := range labels {
		ls = append(ls, *l)
	}

	version, err := extractSemverLabels(ls)
	if err != nil {
		return UNVALID_VERSION, fmt.Errorf("pull request #%d : %s", pullRequestNumber, err)
	}

	return version, nil
}

// CreateList populates a repository with all labels needed
// to version pull requests
func (s SemverLabelService) CreateList() error {
	for _, label := range []struct {
		name        string
		color       string
		description string
	}{
		{
			"norelease",
			"bdbdbd",
			"Produces no new version when pull request is merged on master",
		},
		{
			"alpha",
			"d0bcd5",
			"Produce a new alpha version according to semver when pull request is merged on master",
		},
		{
			"beta",
			"a499b3",
			"Produce a new beta version according to semver when pull request is merged on master",
		},
		{
			"rc",
			"534b62",
			"Produce a new rc version according to semver when pull request is merged on master",
		},
		{
			"patch",
			"0e8a16",
			"Produce a new patch version according to semver when pull request is merged on master",
		},
		{
			"minor",
			"fbca04",
			"Produce a new minor version according to semver when pull request is merged on master",
		},
		{
			"major",
			"d93f0b",
			"Produce a new major version according to semver when pull request is merged on master",
		},
	} {
		if _, _, err := s.client.Issues.CreateLabel(context.Background(), s.owner, s.repository, &github.Label{
			Name:        &label.name,
			Color:       &label.color,
			Description: &label.description,
		}); err != nil {
			return fmt.Errorf("can't create label %s on repository %s : %s", label.name, s.repository, err)
		}
	}

	return nil
}

func extractSemverLabels(labels []github.Label) (Version, error) {
	versions := []Version{}

	for _, label := range labels {
		switch label.GetName() {
		case "norelease":
			versions = append(versions, NORELEASE)
		case "alpha":
			versions = append(versions, ALPHA)
		case "beta":
			versions = append(versions, BETA)
		case "rc":
			versions = append(versions, RC)
		case "patch":
			versions = append(versions, PATCH)
		case "minor":
			versions = append(versions, MINOR)
		case "major":
			versions = append(versions, MAJOR)
		}
	}

	if len(versions) == 0 {
		return UNVALID_VERSION, fmt.Errorf("no semver label found")
	} else if len(versions) > 1 {
		return UNVALID_VERSION, fmt.Errorf("more than one semver label found")
	}

	return versions[0], nil
}
