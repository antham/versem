package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestNewTagFromString(t *testing.T) {
	scenarios := []struct {
		name        string
		getArgument func() string
		test        func(Tag, error)
	}{
		{
			"Parse an unvalid semver tag : whatever",
			func() string {
				return "whatever"
			},
			func(tag Tag, err error) {
				assert.EqualError(t, err, "whatever is not a valid semver tag")
			},
		},
		{
			"Parse a valid semver tag : 1.2.3",
			func() string {
				return "1.2.3"
			},
			func(tag Tag, err error) {
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3}, tag)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.test(NewTagFromString(scenario.getArgument()))
		})
	}
}

func TestParseStringTag(t *testing.T) {
	scenarios := []struct {
		name        string
		getArgument func() string
		test        func(Tag, error)
	}{
		{
			"Parse an unvalid semver tag : whatever",
			func() string {
				return "whatever"
			},
			func(tag Tag, err error) {
				assert.EqualError(t, err, "whatever is not a valid semver tag")
			},
		},
		{
			"Parse an unvalid semver tag : 1.0",
			func() string {
				return "1.0"
			},
			func(tag Tag, err error) {
				assert.EqualError(t, err, "1.0 is not a valid semver tag")
			},
		},
		{
			"Parse an alpha 0 tag",
			func() string {
				return "1.2.3-alpha"
			},
			func(tag Tag, err error) {
				var zero int
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero, PreRelease: "alpha"}, tag)
			},
		},
		{
			"Parse an alpha 1 tag",
			func() string {
				return "1.2.3-alpha.1"
			},
			func(tag Tag, err error) {
				one := 1
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &one, PreRelease: "alpha.1"}, tag)
			},
		},
		{
			"Parse a beta 0 tag",
			func() string {
				return "1.2.3-beta"
			},
			func(tag Tag, err error) {
				var zero int
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero, PreRelease: "beta"}, tag)
			},
		},
		{
			"Parse a beta 1 tag",
			func() string {
				return "1.2.3-beta.1"
			},
			func(tag Tag, err error) {
				one := 1
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &one, PreRelease: "beta.1"}, tag)
			},
		},
		{
			"Parse a rc 0 tag",
			func() string {
				return "1.2.3-rc"
			},
			func(tag Tag, err error) {
				var zero int
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero, PreRelease: "rc"}, tag)
			},
		},
		{
			"Parse a rc 1 tag",
			func() string {
				return "1.2.3-rc.1"
			},
			func(tag Tag, err error) {
				one := 1
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &one, PreRelease: "rc.1"}, tag)
			},
		},
		{
			"Parse a tag with custom prerelease",
			func() string {
				return "1.2.3-1.2.3"
			},
			func(tag Tag, err error) {
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, PreRelease: "1.2.3"}, tag)
			},
		},
		{
			"Parse a tag with rc and build metadatas",
			func() string {
				return "1.2.3-rc.1+20150901.sha.5114f85"
			},
			func(tag Tag, err error) {
				one := 1
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &one, PreRelease: "rc.1", BuildMetadata: "20150901.sha.5114f85"}, tag)
			},
		},
		{
			"Parse a tag with build metadatas",
			func() string {
				return "1.2.3+20150901.sha.5114f85"
			},
			func(tag Tag, err error) {
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, BuildMetadata: "20150901.sha.5114f85"}, tag)
			},
		},
		{
			"Parse a tag minor",
			func() string {
				return "1.2.3"
			},
			func(tag Tag, err error) {
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3}, tag)
			},
		},
		{
			"Parse a tag with leading v",
			func() string {
				return "v1.2.3-rc.1"
			},
			func(tag Tag, err error) {
				one := 1
				assert.NoError(t, err)
				assert.Equal(t, Tag{LeadingV: true, Major: 1, Minor: 2, Patch: 3, RC: &one, PreRelease: "rc.1"}, tag)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.test(parseStringTag(scenario.getArgument()))
		})
	}
}

func TestTagString(t *testing.T) {
	scenarios := []struct {
		name        string
		getArgument func() Tag
		test        func(tag string)
	}{
		{
			"Print tag",
			func() Tag {
				return Tag{Major: 1, Minor: 2, Patch: 3}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3", tag)
			},
		},
		{
			"Print tag with prerelease",
			func() Tag {
				one := 1
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &one, PreRelease: "alpha.1"}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3-alpha.1", tag)
			},
		},
		{
			"Print tag with build metadatas",
			func() Tag {
				return Tag{Major: 1, Minor: 2, Patch: 3, BuildMetadata: "20150909"}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3+20150909", tag)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.test(scenario.getArgument().String())
		})
	}
}

func TestGetNextTag(t *testing.T) {
	scenarios := []struct {
		name         string
		getArguments func() (Tag, Version)
		test         func(tag Tag)
	}{
		{
			"Get next tag from alpha version",
			func() (Tag, Version) {
				return Tag{Major: 1, Minor: 2, Patch: 3}, ALPHA
			},
			func(tag Tag) {
				var zero int
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero, PreRelease: "alpha"}, tag)
			},
		},
		{
			"Get next tag from alpha version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero, PreRelease: "alpha"}, ALPHA
			},
			func(tag Tag) {
				var one = 1
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &one, PreRelease: "alpha.1"}, tag)
			},
		},
		{
			"Get next tag from beta version",
			func() (Tag, Version) {
				return Tag{Major: 1, Minor: 2, Patch: 3}, BETA
			},
			func(tag Tag) {
				var zero int
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero, PreRelease: "beta"}, tag)
			},
		},
		{
			"Get next tag from beta version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero, PreRelease: "beta"}, BETA
			},
			func(tag Tag) {
				var one = 1
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &one, PreRelease: "beta.1"}, tag)
			},
		},
		{
			"Get next tag from rc version",
			func() (Tag, Version) {
				return Tag{Major: 1, Minor: 2, Patch: 3}, RC
			},
			func(tag Tag) {
				var zero int
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero, PreRelease: "rc"}, tag)
			},
		},
		{
			"Get next tag from rc version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero, PreRelease: "rc"}, RC
			},
			func(tag Tag) {
				var one = 1
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &one, PreRelease: "rc.1"}, tag)
			},
		},
		{
			"Get next tag from patch version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero, PreRelease: "alpha"}, PATCH
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 4}, tag)
			},
		},
		{
			"Get next tag from patch version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero, PreRelease: "beta"}, PATCH
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 4}, tag)
			},
		},
		{
			"Get next tag from patch version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero, PreRelease: "rc"}, PATCH
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 4}, tag)
			},
		},
		{
			"Get next tag from minor version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero, PreRelease: "alpha"}, MINOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 3, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from minor version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero, PreRelease: "beta"}, MINOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 3, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from minor version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero, PreRelease: "rc"}, MINOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 3, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from major version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero, PreRelease: "alpha"}, MAJOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 2, Minor: 0, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from major version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero, PreRelease: "beta"}, MAJOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 2, Minor: 0, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from major version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero, PreRelease: "rc"}, MAJOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 2, Minor: 0, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from major version and previous rc tag with a leading v",
			func() (Tag, Version) {
				var zero int
				return Tag{LeadingV: true, Major: 1, Minor: 2, Patch: 3, RC: &zero, PreRelease: "rc"}, MAJOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{LeadingV: true, Major: 2, Minor: 0, Patch: 0}, tag)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.test(getNextTag(scenario.getArguments()))
		})
	}
}

func TestReleaseCreateNext(t *testing.T) {
	defer gock.Off()

	scenarios := []struct {
		name        string
		setup       func()
		getArgument func() Version
		test        func(err error)
	}{
		{
			"An error occurred when fetching last tag from github api",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/tags").
					MatchParams(
						map[string]string{
							"page":     "1",
							"per_page": "1",
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(500).
					JSON(
						map[string]interface{}{
							"message": "An error occurred",
						})
			},
			func() Version {
				return MAJOR
			},
			func(err error) {
				assert.EqualError(t, err, "can't fetch latest tag from versem : GET https://api.github.com/repos/antham/versem/tags?page=1&per_page=1: 500 An error occurred []")
			},
		},
		{
			"An error occurred when creating the release",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/tags").
					MatchParams(
						map[string]string{
							"page":     "1",
							"per_page": "1",
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON([]map[string]interface{}{})

				gock.New("https://api.github.com").
					Post("/repos/antham/versem/releases").
					MatchType("json").
					JSON(
						map[string]interface{}{
							"tag_name":   "1.0.0",
							"prerelease": false,
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(500).
					JSON(
						map[string]interface{}{
							"message": "An error occurred",
						})
			},
			func() Version {
				return MAJOR
			},
			func(err error) {
				assert.EqualError(t, err, "can't create release on versem : POST https://api.github.com/repos/antham/versem/releases: 500 An error occurred []")
			},
		},
		{
			"Create a new release",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/tags").
					MatchParams(
						map[string]string{
							"page":     "1",
							"per_page": "1",
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON([]map[string]interface{}{
						{
							"name": "1.0.0",
							"commit": map[string]interface{}{
								"sha": "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc",
								"url": "https://api.github.com/repos/antham/versem/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc",
							},
							"zipball_url": "https://github.com/antham/versem/zipball/1.0.0",
							"tarball_url": "https://github.com/antham/versem/tarball/1.0.0",
						},
					})

				gock.New("https://api.github.com").
					Post("/repos/antham/versem/releases").
					MatchType("json").
					JSON(
						map[string]interface{}{
							"tag_name":   "2.0.0",
							"prerelease": false,
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(201).
					JSON(
						map[string]interface{}{
							"url":              "https://api.github.com/repos/antham/versem/releases/1",
							"html_url":         "https://github.com/antham/versem/releases/2.0.0",
							"assets_url":       "https://api.github.com/repos/antham/versem/releases/1/assets",
							"upload_url":       "https://uploads.github.com/repos/antham/versem/releases/1/assets{?name,label}",
							"tarball_url":      "https://api.github.com/repos/antham/versem/tarball/2.0.0",
							"zipball_url":      "https://api.github.com/repos/antham/versem/zipball/2.0.0",
							"id":               1,
							"node_id":          "MDc6UmVsZWFzZTE=",
							"tag_name":         "2.0.0",
							"target_commitish": "master",
							"name":             "2.0.0",
							"body":             "Description of the release",
							"draft":            false,
							"prerelease":       false,
							"created_at":       "2013-02-27T19:35:32Z",
							"published_at":     "2013-02-27T19:35:32Z",
							"author": map[string]interface{}{
								"login":               "antham",
								"id":                  1,
								"node_id":             "MDQ6VXNlcjE=",
								"avatar_url":          "https://github.com/images/error/antham_happy.gif",
								"gravatar_id":         "",
								"url":                 "https://api.github.com/users/antham",
								"html_url":            "https://github.com/antham",
								"followers_url":       "https://api.github.com/users/antham/followers",
								"following_url":       "https://api.github.com/users/antham/following{/other_user}",
								"gists_url":           "https://api.github.com/users/antham/gists{/gist_id}",
								"starred_url":         "https://api.github.com/users/antham/starred{/owner}{/repo}",
								"subscriptions_url":   "https://api.github.com/users/antham/subscriptions",
								"organizations_url":   "https://api.github.com/users/antham/orgs",
								"repos_url":           "https://api.github.com/users/antham/repos",
								"events_url":          "https://api.github.com/users/antham/events{/privacy}",
								"received_events_url": "https://api.github.com/users/antham/received_events",
								"type":                "User",
								"site_admin":          false,
							},
							"assets": []string{},
						})
			},
			func() Version {
				return MAJOR
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Create a new prerelease",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/tags").
					MatchParams(
						map[string]string{
							"page":     "1",
							"per_page": "1",
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"name": "1.0.0",
								"commit": map[string]interface{}{
									"sha": "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc",
									"url": "https://api.github.com/repos/antham/versem/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc",
								},
								"zipball_url": "https://github.com/antham/versem/zipball/1.0.0",
								"tarball_url": "https://github.com/antham/versem/tarball/1.0.0",
							},
						})

				gock.New("https://api.github.com").
					Post("/repos/antham/versem/releases").
					MatchType("json").
					JSON(
						map[string]interface{}{
							"tag_name":   "1.0.0-rc",
							"prerelease": true,
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(201).
					JSON(
						map[string]interface{}{
							"url":              "https://api.github.com/repos/antham/versem/releases/1",
							"html_url":         "https://github.com/antham/versem/releases/1.0.0-rc",
							"assets_url":       "https://api.github.com/repos/antham/versem/releases/1/assets",
							"upload_url":       "https://uploads.github.com/repos/antham/versem/releases/1/assets{?name,label}",
							"tarball_url":      "https://api.github.com/repos/antham/versem/tarball/1.0.0-rc",
							"zipball_url":      "https://api.github.com/repos/antham/versem/zipball/1.0.0-rc",
							"id":               1,
							"node_id":          "MDc6UmVsZWFzZTE=",
							"tag_name":         "1.0.0-rc",
							"target_commitish": "master",
							"name":             "1.0.0-rc",
							"body":             "Description of the release",
							"draft":            false,
							"prerelease":       false,
							"created_at":       "2013-02-27T19:35:32Z",
							"published_at":     "2013-02-27T19:35:32Z",
							"author": map[string]interface{}{
								"login":               "antham",
								"id":                  1,
								"node_id":             "MDQ6VXNlcjE=",
								"avatar_url":          "https://github.com/images/error/antham_happy.gif",
								"gravatar_id":         "",
								"url":                 "https://api.github.com/users/antham",
								"html_url":            "https://github.com/antham",
								"followers_url":       "https://api.github.com/users/antham/followers",
								"following_url":       "https://api.github.com/users/antham/following{/other_user}",
								"gists_url":           "https://api.github.com/users/antham/gists{/gist_id}",
								"starred_url":         "https://api.github.com/users/antham/starred{/owner}{/repo}",
								"subscriptions_url":   "https://api.github.com/users/antham/subscriptions",
								"organizations_url":   "https://api.github.com/users/antham/orgs",
								"repos_url":           "https://api.github.com/users/antham/repos",
								"events_url":          "https://api.github.com/users/antham/events{/privacy}",
								"received_events_url": "https://api.github.com/users/antham/received_events",
								"type":                "User",
								"site_admin":          false,
							},
							"assets": []string{},
						})
			},
			func() Version {
				return RC
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Create a release when no release exist yet",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/tags").
					MatchParams(
						map[string]string{
							"page":     "1",
							"per_page": "1",
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{})

				gock.New("https://api.github.com").
					Post("/repos/antham/versem/releases").
					MatchType("json").
					JSON(
						map[string]interface{}{
							"tag_name":   "1.0.0",
							"prerelease": false,
						},
					).
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(201).
					JSON(
						map[string]interface{}{
							"url":              "https://api.github.com/repos/antham/versem/releases/1",
							"html_url":         "https://github.com/antham/versem/releases/1.0.0",
							"assets_url":       "https://api.github.com/repos/antham/versem/releases/1/assets",
							"upload_url":       "https://uploads.github.com/repos/antham/versem/releases/1/assets{?name,label}",
							"tarball_url":      "https://api.github.com/repos/antham/versem/tarball/1.0.0",
							"zipball_url":      "https://api.github.com/repos/antham/versem/zipball/1.0.0",
							"id":               1,
							"node_id":          "MDc6UmVsZWFzZTE=",
							"tag_name":         "1.0.0",
							"target_commitish": "master",
							"name":             "1.0.0",
							"body":             "Description of the release",
							"draft":            false,
							"prerelease":       false,
							"created_at":       "2013-02-27T19:35:32Z",
							"published_at":     "2013-02-27T19:35:32Z",
							"author": map[string]interface{}{
								"login":               "antham",
								"id":                  1,
								"node_id":             "MDQ6VXNlcjE=",
								"avatar_url":          "https://github.com/images/error/antham_happy.gif",
								"gravatar_id":         "",
								"url":                 "https://api.github.com/users/antham",
								"html_url":            "https://github.com/antham",
								"followers_url":       "https://api.github.com/users/antham/followers",
								"following_url":       "https://api.github.com/users/antham/following{/other_user}",
								"gists_url":           "https://api.github.com/users/antham/gists{/gist_id}",
								"starred_url":         "https://api.github.com/users/antham/starred{/owner}{/repo}",
								"subscriptions_url":   "https://api.github.com/users/antham/subscriptions",
								"organizations_url":   "https://api.github.com/users/antham/orgs",
								"repos_url":           "https://api.github.com/users/antham/repos",
								"events_url":          "https://api.github.com/users/antham/events{/privacy}",
								"received_events_url": "https://api.github.com/users/antham/received_events",
								"type":                "User",
								"site_admin":          false,
							},
							"assets": []string{},
						})
			},
			func() Version {
				return MAJOR
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.setup()
			s := NewReleaseService("antham", "versem", "396531004112aa66a7fda31bfdca7d00")
			scenario.test(s.CreateNext(scenario.getArgument()))
			assert.True(t, gock.IsDone())
		})
	}
}
