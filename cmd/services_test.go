package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

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
