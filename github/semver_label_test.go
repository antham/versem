package github

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestSemverLabelServiceGetFromPullRequest(t *testing.T) {
	defer gock.Off()

	scenarios := []struct {
		name  string
		setup func()
		test  func(version Version, err error)
	}{
		{
			"An error occurred when requesting github api",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(400).
					JSON(
						map[string]interface{}{
							"message": "Problems parsing JSON",
						})
			},
			func(version Version, err error) {
				assert.EqualError(t, err, "can't fetch github api to get label for pull request #1 : GET https://api.github.com/repos/antham/versem/issues/1/labels: 400 Problems parsing JSON []")
			},
		},
		{
			"No semver label attached to the pull request",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/bug",
								"name":    "bug",
								"color":   "f29513",
								"default": true,
							},
						},
					)
			},
			func(version Version, err error) {
				assert.EqualError(t, err, "pull request #1 : no semver label found")
			},
		},
		{
			"Several semver label attached to the pull request",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/major",
								"name":    "major",
								"color":   "f29513",
								"default": true,
							},
							{
								"id":      208045947,
								"url":     "https://api.github.com/repos/antham/versem/labels/minor",
								"name":    "minor",
								"color":   "f29513",
								"default": true,
							},
						},
					)
			},
			func(version Version, err error) {
				assert.EqualError(t, err, "pull request #1 : more than one semver label found")
			},
		},
		{
			"Fetch norelease label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/norelease",
								"name":    "norelease",
								"color":   "f29513",
								"default": true,
							},
							{
								"id":      208045947,
								"url":     "https://api.github.com/repos/antham/versem/labels/bug",
								"name":    "bug",
								"color":   "f29513",
								"default": true,
							},
						},
					)
			},
			func(version Version, err error) {
				assert.NoError(t, err)
				assert.Equal(t, NORELEASE, version)
			},
		},
		{
			"Fetch patch label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/patch",
								"name":    "patch",
								"color":   "f29513",
								"default": true,
							},
							{
								"id":      208045947,
								"url":     "https://api.github.com/repos/antham/versem/labels/bug",
								"name":    "bug",
								"color":   "f29513",
								"default": true,
							},
						},
					)
			},
			func(version Version, err error) {
				assert.NoError(t, err)
				assert.Equal(t, PATCH, version)
			},
		},
		{
			"Fetch minor label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/minor",
								"name":    "minor",
								"color":   "f29513",
								"default": true,
							},
							{
								"id":      208045947,
								"url":     "https://api.github.com/repos/antham/versem/labels/bug",
								"name":    "bug",
								"color":   "f29513",
								"default": true,
							},
						},
					)
			},
			func(version Version, err error) {
				assert.NoError(t, err)
				assert.Equal(t, MINOR, version)
			},
		},
		{
			"Fetch major label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/major",
								"name":    "major",
								"color":   "f29513",
								"default": true,
							},
							{
								"id":      208045947,
								"url":     "https://api.github.com/repos/antham/versem/labels/bug",
								"name":    "bug",
								"color":   "f29513",
								"default": true,
							},
						},
					)
			},
			func(version Version, err error) {
				assert.NoError(t, err)
				assert.Equal(t, MAJOR, version)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.setup()
			s := NewSemverLabelService("antham", "versem", "396531004112aa66a7fda31bfdca7d00")
			scenario.test(s.GetFromPullRequest(1))
			assert.True(t, gock.IsDone())
		})
	}
}

func TestSemverLabelServiceGetFromCommitSha(t *testing.T) {
	defer gock.Off()

	scenarios := []struct {
		name  string
		setup func()
		test  func(version Version, err error)
	}{
		{
			"An error occurred when requesting github api",
			func() {
				gock.New("https://api.github.com").
					Get("/search/issues").
					MatchParam("q", "a6e6c8b8c34d2382e591587e960e7e7f825cb221.repo.antham.versem").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(400).
					JSON(
						map[string]interface{}{
							"message": "Problems parsing JSON",
						})
			},
			func(version Version, err error) {
				assert.EqualError(t, err, "can't fetch github api to get label from commit a6e6c8b8c34d2382e591587e960e7e7f825cb221 : GET https://api.github.com/search/issues?q=a6e6c8b8c34d2382e591587e960e7e7f825cb221+repo:antham%2Fversem: 400 Problems parsing JSON []")
			},
		},
		{
			"Commit not found",
			func() {
				gock.New("https://api.github.com").
					Get("/search/issues").
					MatchParam("q", "a6e6c8b8c34d2382e591587e960e7e7f825cb221.repo.antham.versem").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						map[string]interface{}{
							"total_count":        0,
							"incomplete_results": false,
							"items":              []struct{}{},
						},
					)
			},
			func(version Version, err error) {
				assert.EqualError(t, err, "commit a6e6c8b8c34d2382e591587e960e7e7f825cb221 not found in antham/versem")
			},
		},
		{
			"Several commit found",
			func() {
				gock.New("https://api.github.com").
					Get("/search/issues").
					MatchParam("q", "a6e6c8b8c34d2382e591587e960e7e7f825cb221.repo.antham.versem").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						map[string]interface{}{
							"total_count":        0,
							"incomplete_results": false,
							"items": []interface{}{
								map[string]interface{}{
									"url":            "https://api.github.com/repos/antham/versem/issues/132",
									"repository_url": "https://api.github.com/repos/antham/versem",
									"labels_url":     "https://api.github.com/repos/antham/versem/issues/132/labels{/name}",
									"comments_url":   "https://api.github.com/repos/antham/versem/issues/132/comments",
									"events_url":     "https://api.github.com/repos/antham/versem/issues/132/events",
									"html_url":       "https://github.com/antham/versem/issues/132",
									"id":             35802,
									"node_id":        "MDU6SXNzdWUzNTgwMg==",
									"number":         132,
									"title":          "Some pull request",
									"user": map[string]interface{}{
										"login":               "Nick3C",
										"id":                  90254,
										"node_id":             "MDQ6VXNlcjkwMjU0",
										"avatar_url":          "https://secure.gravatar.com/avatar/934442aadfe3b2f4630510de416c5718?d=https://a248.e.akamai.net/assets.github.com%2Fimages%2Fgravatars%2Fgravatar-user-420.png",
										"gravatar_id":         "",
										"url":                 "https://api.github.com/users/Nick3C",
										"html_url":            "https://github.com/Nick3C",
										"followers_url":       "https://api.github.com/users/Nick3C/followers",
										"following_url":       "https://api.github.com/users/Nick3C/following{/other_user}",
										"gists_url":           "https://api.github.com/users/Nick3C/gists{/gist_id}",
										"starred_url":         "https://api.github.com/users/Nick3C/starred{/owner}{/repo}",
										"subscriptions_url":   "https://api.github.com/users/Nick3C/subscriptions",
										"organizations_url":   "https://api.github.com/users/Nick3C/orgs",
										"repos_url":           "https://api.github.com/users/Nick3C/repos",
										"events_url":          "https://api.github.com/users/Nick3C/events{/privacy}",
										"received_events_url": "https://api.github.com/users/Nick3C/received_events",
										"type":                "User",
									},
									"labels": []map[string]interface{}{
										{
											"id":      4,
											"node_id": "MDU6TGFiZWw0",
											"url":     "https://api.github.com/repos/antham/versem/labels/bug",
											"name":    "bug",
											"color":   "ff0000",
										},
									},
									"state":      "open",
									"assignee":   nil,
									"milestone":  nil,
									"comments":   15,
									"created_at": "2009-07-12T20:10:41Z",
									"updated_at": "2009-07-19T09:23:43Z",
									"closed_at":  nil,
									"pull_request": map[string]interface{}{
										"html_url":  nil,
										"diff_url":  nil,
										"patch_url": nil,
									},
									"body":  "...",
									"score": 1.3859273,
								},
								map[string]interface{}{
									"url":            "https://api.github.com/repos/antham/versem/issues/133",
									"repository_url": "https://api.github.com/repos/antham/versem",
									"labels_url":     "https://api.github.com/repos/antham/versem/issues/133/labels{/name}",
									"comments_url":   "https://api.github.com/repos/antham/versem/issues/133/comments",
									"events_url":     "https://api.github.com/repos/antham/versem/issues/133/events",
									"html_url":       "https://github.com/antham/versem/issues/133",
									"id":             35803,
									"node_id":        "MDU6SXNzdWUzNTgwMg==",
									"number":         133,
									"title":          "Some pull request",
									"user": map[string]interface{}{
										"login":               "Nick3C",
										"id":                  90254,
										"node_id":             "MDQ6VXNlcjkwMjU0",
										"avatar_url":          "https://secure.gravatar.com/avatar/934442aadfe3b2f4630510de416c5718?d=https://a248.e.akamai.net/assets.github.com%2Fimages%2Fgravatars%2Fgravatar-user-420.png",
										"gravatar_id":         "",
										"url":                 "https://api.github.com/users/Nick3C",
										"html_url":            "https://github.com/Nick3C",
										"followers_url":       "https://api.github.com/users/Nick3C/followers",
										"following_url":       "https://api.github.com/users/Nick3C/following{/other_user}",
										"gists_url":           "https://api.github.com/users/Nick3C/gists{/gist_id}",
										"starred_url":         "https://api.github.com/users/Nick3C/starred{/owner}{/repo}",
										"subscriptions_url":   "https://api.github.com/users/Nick3C/subscriptions",
										"organizations_url":   "https://api.github.com/users/Nick3C/orgs",
										"repos_url":           "https://api.github.com/users/Nick3C/repos",
										"events_url":          "https://api.github.com/users/Nick3C/events{/privacy}",
										"received_events_url": "https://api.github.com/users/Nick3C/received_events",
										"type":                "User",
									},
									"labels": []map[string]interface{}{
										{
											"id":      4,
											"node_id": "MDU6TGFiZWw0",
											"url":     "https://api.github.com/repos/antham/versem/labels/bug",
											"name":    "bug",
											"color":   "ff0000",
										},
									},
									"state":      "open",
									"assignee":   nil,
									"milestone":  nil,
									"comments":   15,
									"created_at": "2009-07-12T20:10:41Z",
									"updated_at": "2009-07-19T09:23:43Z",
									"closed_at":  nil,
									"pull_request": map[string]interface{}{
										"html_url":  nil,
										"diff_url":  nil,
										"patch_url": nil,
									},
									"body":  "...",
									"score": 1.3859273,
								},
							},
						},
					)
			},
			func(version Version, err error) {
				assert.EqualError(t, err, "several entries found for commit a6e6c8b8c34d2382e591587e960e7e7f825cb221 in antham/versem")
			},
		},
		{
			"Several semver label found",
			func() {
				gock.New("https://api.github.com").
					Get("/search/issues").
					MatchParam("q", "a6e6c8b8c34d2382e591587e960e7e7f825cb221.repo.antham.versem").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						map[string]interface{}{
							"total_count":        0,
							"incomplete_results": false,
							"items": []interface{}{
								map[string]interface{}{
									"url":            "https://api.github.com/repos/antham/versem/issues/132",
									"repository_url": "https://api.github.com/repos/antham/versem",
									"labels_url":     "https://api.github.com/repos/antham/versem/issues/132/labels{/name}",
									"comments_url":   "https://api.github.com/repos/antham/versem/issues/132/comments",
									"events_url":     "https://api.github.com/repos/antham/versem/issues/132/events",
									"html_url":       "https://github.com/antham/versem/issues/132",
									"id":             35802,
									"node_id":        "MDU6SXNzdWUzNTgwMg==",
									"number":         132,
									"title":          "Some pull request",
									"user": map[string]interface{}{
										"login":               "Nick3C",
										"id":                  90254,
										"node_id":             "MDQ6VXNlcjkwMjU0",
										"avatar_url":          "https://secure.gravatar.com/avatar/934442aadfe3b2f4630510de416c5718?d=https://a248.e.akamai.net/assets.github.com%2Fimages%2Fgravatars%2Fgravatar-user-420.png",
										"gravatar_id":         "",
										"url":                 "https://api.github.com/users/Nick3C",
										"html_url":            "https://github.com/Nick3C",
										"followers_url":       "https://api.github.com/users/Nick3C/followers",
										"following_url":       "https://api.github.com/users/Nick3C/following{/other_user}",
										"gists_url":           "https://api.github.com/users/Nick3C/gists{/gist_id}",
										"starred_url":         "https://api.github.com/users/Nick3C/starred{/owner}{/repo}",
										"subscriptions_url":   "https://api.github.com/users/Nick3C/subscriptions",
										"organizations_url":   "https://api.github.com/users/Nick3C/orgs",
										"repos_url":           "https://api.github.com/users/Nick3C/repos",
										"events_url":          "https://api.github.com/users/Nick3C/events{/privacy}",
										"received_events_url": "https://api.github.com/users/Nick3C/received_events",
										"type":                "User",
									},
									"labels": []map[string]interface{}{
										{
											"id":      4,
											"node_id": "MDU6TGFiZWw1",
											"url":     "https://api.github.com/repos/antham/versem/labels/minor",
											"name":    "minor",
											"color":   "ff0000",
										},
										{
											"id":      5,
											"node_id": "MDU6TGFiZWw0",
											"url":     "https://api.github.com/repos/antham/versem/labels/major",
											"name":    "major",
											"color":   "ff0000",
										},
									},
									"state":      "open",
									"assignee":   nil,
									"milestone":  nil,
									"comments":   15,
									"created_at": "2009-07-12T20:10:41Z",
									"updated_at": "2009-07-19T09:23:43Z",
									"closed_at":  nil,
									"pull_request": map[string]interface{}{
										"html_url":  nil,
										"diff_url":  nil,
										"patch_url": nil,
									},
									"body":  "...",
									"score": 1.3859273,
								},
							},
						},
					)
			},
			func(version Version, err error) {
				assert.EqualError(t, err, "commit a6e6c8b8c34d2382e591587e960e7e7f825cb221 in antham/versem : more than one semver label found")
			},
		},
		{
			"Fetch minor label",
			func() {
				gock.New("https://api.github.com").
					Get("/search/issues").
					MatchParam("q", "a6e6c8b8c34d2382e591587e960e7e7f825cb221.repo.antham.versem").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(200).
					JSON(
						map[string]interface{}{
							"total_count":        0,
							"incomplete_results": false,
							"items": []interface{}{
								map[string]interface{}{
									"url":            "https://api.github.com/repos/antham/versem/issues/132",
									"repository_url": "https://api.github.com/repos/antham/versem",
									"labels_url":     "https://api.github.com/repos/antham/versem/issues/132/labels{/name}",
									"comments_url":   "https://api.github.com/repos/antham/versem/issues/132/comments",
									"events_url":     "https://api.github.com/repos/antham/versem/issues/132/events",
									"html_url":       "https://github.com/antham/versem/issues/132",
									"id":             35802,
									"node_id":        "MDU6SXNzdWUzNTgwMg==",
									"number":         132,
									"title":          "Some pull request",
									"user": map[string]interface{}{
										"login":               "Nick3C",
										"id":                  90254,
										"node_id":             "MDQ6VXNlcjkwMjU0",
										"avatar_url":          "https://secure.gravatar.com/avatar/934442aadfe3b2f4630510de416c5718?d=https://a248.e.akamai.net/assets.github.com%2Fimages%2Fgravatars%2Fgravatar-user-420.png",
										"gravatar_id":         "",
										"url":                 "https://api.github.com/users/Nick3C",
										"html_url":            "https://github.com/Nick3C",
										"followers_url":       "https://api.github.com/users/Nick3C/followers",
										"following_url":       "https://api.github.com/users/Nick3C/following{/other_user}",
										"gists_url":           "https://api.github.com/users/Nick3C/gists{/gist_id}",
										"starred_url":         "https://api.github.com/users/Nick3C/starred{/owner}{/repo}",
										"subscriptions_url":   "https://api.github.com/users/Nick3C/subscriptions",
										"organizations_url":   "https://api.github.com/users/Nick3C/orgs",
										"repos_url":           "https://api.github.com/users/Nick3C/repos",
										"events_url":          "https://api.github.com/users/Nick3C/events{/privacy}",
										"received_events_url": "https://api.github.com/users/Nick3C/received_events",
										"type":                "User",
									},
									"labels": []map[string]interface{}{
										{
											"id":      4,
											"node_id": "MDU6TGFiZWw1",
											"url":     "https://api.github.com/repos/antham/versem/labels/minor",
											"name":    "minor",
											"color":   "ff0000",
										},
									},
									"state":      "open",
									"assignee":   nil,
									"milestone":  nil,
									"comments":   15,
									"created_at": "2009-07-12T20:10:41Z",
									"updated_at": "2009-07-19T09:23:43Z",
									"closed_at":  nil,
									"pull_request": map[string]interface{}{
										"html_url":  nil,
										"diff_url":  nil,
										"patch_url": nil,
									},
									"body":  "...",
									"score": 1.3859273,
								},
							},
						},
					)
			},
			func(version Version, err error) {
				assert.Nil(t, err)
				assert.Equal(t, MINOR, version)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.setup()
			s := NewSemverLabelService("antham", "versem", "396531004112aa66a7fda31bfdca7d00")
			scenario.test(s.GetFromCommit("a6e6c8b8c34d2382e591587e960e7e7f825cb221"))
			assert.True(t, gock.IsDone())
		})
	}
}

func TestSemverLabelServiceCreateList(t *testing.T) {
	defer gock.Off()

	scenarios := []struct {
		name  string
		setup func()
		test  func(err error)
	}{
		{
			"An error occurred when requesting github api",
			func() {
				gock.New("https://api.github.com").
					Post("/repos/antham/versem/labels").
					MatchHeader("Authorization", "Bearer 396531004112aa66a7fda31bfdca7d00").
					Reply(400).
					JSON(
						map[string]interface{}{
							"message": "Problems parsing JSON",
						})
			},
			func(err error) {
				assert.EqualError(t, err, "can't create label norelease on repository versem : POST https://api.github.com/repos/antham/versem/labels: 400 Problems parsing JSON []")
			},
		},
		{
			"Create all semver labels on repository",
			func() {
				for _, label := range [][]string{
					{
						"norelease",
						"bdbdbd",
						"Produces no new version when pull request is merged on master",
					},
					{
						"patch",
						"0e8a16",
						"Produce a new semver patch version when pull request is merged on master",
					},
					{
						"minor",
						"fbca04",
						"Produce a new semver minor version when pull request is merged on master",
					},
					{
						"major",
						"d93f0b",
						"Produce a new semver major version when pull request is merged on master",
					},
				} {
					gock.New("https://api.github.com").
						Post("/repos/antham/versem/labels").
						MatchHeader("Authorization", "earer 396531004112aa66a7fda31bfdca7d00").
						MatchType("json").
						JSON(
							map[string]interface{}{
								"name":        label[0],
								"color":       label[1],
								"description": label[2],
							},
						).
						Reply(201)
				}
			},
			func(err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.setup()
			s := NewSemverLabelService("antham", "versem", "396531004112aa66a7fda31bfdca7d00")
			scenario.test(s.CreateList())
			assert.True(t, gock.IsDone())
		})
	}
}

func TestExtractSemverLabels(t *testing.T) {
	scenarios := []struct {
		name        string
		getArgument func() []github.Label
		test        func(Version, error)
	}{
		{
			"Several release label defined",
			func() []github.Label {
				major := "major"
				minor := "minor"

				return []github.Label{
					{
						Name: &major,
					},
					{
						Name: &minor,
					},
				}
			},
			func(versions Version, err error) {
				assert.EqualError(t, err, "more than one semver label found")
			},
		},
		{
			"No release labels found",
			func() []github.Label {
				return []github.Label{}
			},
			func(versions Version, err error) {
				assert.EqualError(t, err, "no semver label found")
			},
		},
		{
			"Patch label found",
			func() []github.Label {
				patch := "patch"
				return []github.Label{
					{
						Name: &patch,
					},
				}
			},
			func(version Version, err error) {
				assert.NoError(t, err)
				assert.Equal(t, PATCH, version)
			},
		},
		{
			"Minor label found",
			func() []github.Label {
				minor := "minor"
				return []github.Label{
					{
						Name: &minor,
					},
				}
			},
			func(version Version, err error) {
				assert.NoError(t, err)
				assert.Equal(t, MINOR, version)
			},
		},
		{
			"Major label found",
			func() []github.Label {
				major := "major"
				return []github.Label{
					{
						Name: &major,
					},
				}
			},
			func(version Version, err error) {
				assert.NoError(t, err)
				assert.Equal(t, MAJOR, version)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(*testing.T) {
			scenario.test(extractSemverLabels(scenario.getArgument()))
		})
	}
}
