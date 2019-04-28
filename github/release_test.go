package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
