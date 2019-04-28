package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero}, tag)
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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &one}, tag)
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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero}, tag)
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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &one}, tag)
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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero}, tag)
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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &one}, tag)
			},
		},
		{
			"Parse a tag with custom prerelease",
			func() string {
				return "1.2.3-1.2.3"
			},
			func(tag Tag, err error) {
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3}, tag)
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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &one}, tag)
			},
		},
		{
			"Parse a tag with build metadatas",
			func() string {
				return "1.2.3+20150901.sha.5114f85"
			},
			func(tag Tag, err error) {
				assert.NoError(t, err)
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3}, tag)
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
	}

	for _, scenario := range scenarios {
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
			"Print alpha tag version 0",
			func() Tag {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3-alpha", tag)
			},
		},
		{
			"Print alpha tag version 1",
			func() Tag {
				one := 1
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &one}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3-alpha.1", tag)
			},
		},
		{
			"Print beta tag version 0",
			func() Tag {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3-beta", tag)
			},
		},
		{
			"Print beta tag version 1",
			func() Tag {
				one := 1
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &one}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3-beta.1", tag)
			},
		},
		{
			"Print rc tag version 0",
			func() Tag {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3-rc", tag)
			},
		},
		{
			"Print rc tag version 1",
			func() Tag {
				one := 1
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &one}
			},
			func(tag string) {
				assert.Equal(t, "1.2.3-rc.1", tag)
			},
		},
	}

	for _, scenario := range scenarios {
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
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero}, tag)
			},
		},
		{
			"Get next tag from alpha version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero}, ALPHA
			},
			func(tag Tag) {
				var one = 1
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &one}, tag)
			},
		},
		{
			"Get next tag from beta version",
			func() (Tag, Version) {
				return Tag{Major: 1, Minor: 2, Patch: 3}, BETA
			},
			func(tag Tag) {
				var zero int
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero}, tag)
			},
		},
		{
			"Get next tag from beta version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero}, BETA
			},
			func(tag Tag) {
				var one = 1
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, Beta: &one}, tag)
			},
		},
		{
			"Get next tag from rc version",
			func() (Tag, Version) {
				return Tag{Major: 1, Minor: 2, Patch: 3}, RC
			},
			func(tag Tag) {
				var zero int
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero}, tag)
			},
		},
		{
			"Get next tag from rc version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero}, RC
			},
			func(tag Tag) {
				var one = 1
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 3, RC: &one}, tag)
			},
		},
		{
			"Get next tag from patch version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero}, PATCH
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 4}, tag)
			},
		},
		{
			"Get next tag from patch version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero}, PATCH
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 4}, tag)
			},
		},
		{
			"Get next tag from patch version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero}, PATCH
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 2, Patch: 4}, tag)
			},
		},
		{
			"Get next tag from minor version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero}, MINOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 3, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from minor version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero}, MINOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 3, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from minor version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero}, MINOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 1, Minor: 3, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from major version and previous alpha tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Alpha: &zero}, MAJOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 2, Minor: 0, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from major version and previous beta tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, Beta: &zero}, MAJOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 2, Minor: 0, Patch: 0}, tag)
			},
		},
		{
			"Get next tag from major version and previous rc tag",
			func() (Tag, Version) {
				var zero int
				return Tag{Major: 1, Minor: 2, Patch: 3, RC: &zero}, MAJOR
			},
			func(tag Tag) {
				assert.Equal(t, Tag{Major: 2, Minor: 0, Patch: 0}, tag)
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(*testing.T) {
			scenario.test(getNextTag(scenario.getArguments()))
		})
	}
}
