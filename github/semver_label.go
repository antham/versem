package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const noreleaseStr = "norelease"
const patchStr = "patch"
const minorStr = "minor"
const majorStr = "major"

// Version represents a semver version
// NORELEASE  => no version
// PATCH      => 0.0.1
// MINOR      => 0.1.0
// MAJOR      => 1.0.0
type Version int

//go:generate stringer -type=Version
const (
	// UNVALIDVERSION is default when no semver version exists
	UNVALIDVERSION Version = iota
	// NORELEASE don't create any release
	NORELEASE
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
	httpClient *http.Client
}

// NewSemverLabelService creates a new instance of SemverLabelService
func NewSemverLabelService(owner string, repository string, token string) SemverLabelService {
	client := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	)

	return SemverLabelService{
		owner,
		repository,
		github.NewClient(client),
		client,
	}
}

// GetFromPullRequest find out semver version label attached to a pull request, if there is none or more than one
// this function returns an error
func (s SemverLabelService) GetFromPullRequest(pullRequestNumber int) (Version, error) {
	labels, _, err := s.client.Issues.ListLabelsByIssue(context.Background(), s.owner, s.repository, pullRequestNumber, nil)
	if err != nil {
		return UNVALIDVERSION, fmt.Errorf("can't fetch github api to get label for pull request #%d : %s", pullRequestNumber, err)
	}

	ls := []github.Label{}

	for _, l := range labels {
		ls = append(ls, *l)
	}

	version, err := extractSemverLabels(ls)
	if err != nil {
		return UNVALIDVERSION, fmt.Errorf("an error occurred when parsing version from pull request #%d : %s", pullRequestNumber, err)
	}

	return version, nil
}

// GetFromCommit find out version label attached to a commit,
// if the commit doesn't exist or multiple pull requests exist for this commit, it returns an error,
// if there is more than one this function returns an error,
// if the commit doesn't belong to a pull request it returns a NORELEASE version
func (s SemverLabelService) GetFromCommit(commitSha string) (Version, error) {
	type PullRequest struct {
		Labels []github.Label
	}

	results := []PullRequest{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s/pulls", s.owner, s.repository, commitSha), nil)
	if err != nil {
		return UNVALIDVERSION, fmt.Errorf("can't create github request : %s", err)
	}

	req.Header.Add("Accept", "application/vnd.github.groot-preview+json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return UNVALIDVERSION, fmt.Errorf("can't fetch github api to get label from commit %s : %s", commitSha, err)
	}

	if res.StatusCode != http.StatusOK {
		return UNVALIDVERSION, fmt.Errorf("can't fetch github api to get label from commit %s : status code %d", commitSha, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return UNVALIDVERSION, fmt.Errorf("can't parse github response : %s", err)
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		return UNVALIDVERSION, fmt.Errorf("can't decode github json response : %s", err)
	}

	if len(results) == 0 {
		return NORELEASE, nil
	} else if len(results) > 1 {
		return UNVALIDVERSION, fmt.Errorf("several entries found associated with commit %s", commitSha)
	}

	version, err := extractSemverLabels(results[0].Labels)
	if err != nil {
		return UNVALIDVERSION, fmt.Errorf("can't parse version from commit %s : %s", commitSha, err)
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
			noreleaseStr,
			"bdbdbd",
			"Produces no new version when pull request is merged on master",
		},
		{
			patchStr,
			"0e8a16",
			"Produce a new semver patch version when pull request is merged on master",
		},
		{
			minorStr,
			"fbca04",
			"Produce a new semver minor version when pull request is merged on master",
		},
		{
			majorStr,
			"d93f0b",
			"Produce a new semver major version when pull request is merged on master",
		},
	} {
		label := label

		if _, _, err := s.client.Issues.CreateLabel(context.Background(), s.owner, s.repository, &github.Label{
			Name:        &label.name,
			Color:       &label.color,
			Description: &label.description,
		}); err != nil {
			return fmt.Errorf("can't create label %s : %s", label.name, err)
		}
	}

	return nil
}

func extractSemverLabels(labels []github.Label) (Version, error) {
	versions := []Version{}

	for _, label := range labels {
		switch label.GetName() {
		case noreleaseStr:
			versions = append(versions, NORELEASE)
		case patchStr:
			versions = append(versions, PATCH)
		case minorStr:
			versions = append(versions, MINOR)
		case majorStr:
			versions = append(versions, MAJOR)
		}
	}

	if len(versions) == 0 {
		return UNVALIDVERSION, fmt.Errorf("no semver label found")
	} else if len(versions) > 1 {
		return UNVALIDVERSION, fmt.Errorf("more than one semver label found")
	}

	return versions[0], nil
}
