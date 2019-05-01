package cmd

import (
	"os"
	"testing"

	"github.com/antham/versem/github"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type semverServiceMock struct {
	version         github.Version
	err             error
	methodCallCount map[string]int
}

func (s semverServiceMock) GetFromPullRequest(int) (github.Version, error) {
	s.methodCallCount["GetFromPullRequest"]++
	return s.version, s.err
}

func (s semverServiceMock) GetFromCommit(string) (github.Version, error) {
	s.methodCallCount["GetFromCommit"]++
	return s.version, s.err
}

func (s semverServiceMock) CreateList() error {
	s.methodCallCount["CreateList"]++
	return s.err
}

func TestGetCredentials(t *testing.T) {
	os.Setenv("GITHUB_OWNER", "antham")
	os.Setenv("GITHUB_REPOSITORY", "versem")
	os.Setenv("GITHUB_TOKEN", "token")
	viper.AutomaticEnv()
	owner, repository, token := getCredentials()
	assert.Equal(t, "antham", owner)
	assert.Equal(t, "versem", repository)
	assert.Equal(t, "token", token)
	os.Unsetenv("GITHUB_OWNER")
	os.Unsetenv("GITHUB_REPOSITORY")
	os.Unsetenv("GITHUB_TOKEN")
}
