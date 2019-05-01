package cmd

import (
	"github.com/antham/versem/github"
	"github.com/spf13/viper"
)

type semverService interface {
	GetFromPullRequest(int) (github.Version, error)
	GetFromCommit(string) (github.Version, error)
	CreateList() error
}

type releaseService interface {
	CreateNext(github.Version, string) (github.Tag, error)
}

func getCredentials() (owner string, repository string, token string) {
	for _, s := range []struct {
		name string
		ptr  *string
	}{
		{
			githubOwner,
			&owner,
		},
		{
			githubRepository,
			&repository,
		},
		{
			githubToken,
			&token,
		},
	} {
		*s.ptr = viper.GetString(s.name)
	}
	return
}

func getSemverLabelService() github.SemverLabelService {
	return github.NewSemverLabelService(getCredentials())
}

func getReleaseService() github.ReleaseService {
	return github.NewReleaseService(getCredentials())
}
