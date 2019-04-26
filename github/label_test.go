package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestSemverLabelServiceGet(t *testing.T) {
	scenarios := []struct {
		name  string
		setup func()
		test  func(version Version, err error)
	}{
		{
			"An error occured when requesting github api",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
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
				assert.EqualError(t, err, "no semver label attached to the pull request #1")
			},
		},
		{
			"Several semver label attached to the pull request",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
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
				assert.EqualError(t, err, "more than one semver label attached to the pull request #1")
			},
		},
		{
			"Fetch norelease label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
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
			"Fetch alpha label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/alpha",
								"name":    "alpha",
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
				assert.Equal(t, ALPHA, version)
			},
		},
		{
			"Fetch beta label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/beta",
								"name":    "beta",
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
				assert.Equal(t, BETA, version)
			},
		},
		{
			"Fetch rc label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
					Reply(200).
					JSON(
						[]map[string]interface{}{
							{
								"id":      208045946,
								"url":     "https://api.github.com/repos/antham/versem/labels/rc",
								"name":    "rc",
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
				assert.Equal(t, RC, version)
			},
		},
		{
			"Fetch patch label",
			func() {
				gock.New("https://api.github.com").
					Get("/repos/antham/versem/issues/1/labels").
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
		t.Run(scenario.name, func(*testing.T) {
			defer gock.Off()
			scenario.setup()
			s := NewSemverLabelService("antham", "versem", "396531004112aa66a7fda31bfdca7d00")
			scenario.test(s.Get(1))
		})
	}
}
