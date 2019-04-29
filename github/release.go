package github

import (
	"fmt"
	"regexp"
	"strconv"
)

// Tag represents a semver tag
type Tag struct {
	Major int
	Minor int
	Patch int
	RC    *int
	Beta  *int
	Alpha *int
}

// String converts a Tag to its string representation
func (t Tag) String() string {
	main := fmt.Sprintf("%d.%d.%d", t.Major, t.Minor, t.Patch)

	switch {
	case t.Alpha != nil && (*t.Alpha) == 0:
		return fmt.Sprintf("%s-alpha", main)
	case t.Alpha != nil && (*t.Alpha) > 0:
		return fmt.Sprintf("%s-alpha.%d", main, *t.Alpha)
	case t.Beta != nil && (*t.Beta) == 0:
		return fmt.Sprintf("%s-beta", main)
	case t.Beta != nil && (*t.Beta) > 0:
		return fmt.Sprintf("%s-beta.%d", main, *t.Beta)
	case t.RC != nil && (*t.RC) == 0:
		return fmt.Sprintf("%s-rc", main)
	case t.RC != nil && (*t.RC) > 0:
		return fmt.Sprintf("%s-rc.%d", main, *t.RC)
	}

	return main
}

// NewTagFromString creates a Tag instance from a tag string representation
func NewTagFromString(tag string) (Tag, error) {
	return parseStringTag(tag)
}
func getNextTag(previousTag Tag, version Version) Tag {
	switch version {
	case ALPHA:
		var v int
		if previousTag.Alpha != nil {
			v = *previousTag.Alpha + 1
		}

		return Tag{
			Major: previousTag.Major,
			Minor: previousTag.Minor,
			Patch: previousTag.Patch,
			Alpha: &v,
		}
	case BETA:
		var v int
		if previousTag.Beta != nil {
			v = *previousTag.Beta + 1
		}

		return Tag{
			Major: previousTag.Major,
			Minor: previousTag.Minor,
			Patch: previousTag.Patch,
			Beta:  &v,
		}
	case RC:
		var v int
		if previousTag.RC != nil {
			v = *previousTag.RC + 1
		}

		return Tag{
			Major: previousTag.Major,
			Minor: previousTag.Minor,
			Patch: previousTag.Patch,
			RC:    &v,
		}
	case PATCH:
		return Tag{
			Major: previousTag.Major,
			Minor: previousTag.Minor,
			Patch: previousTag.Patch + 1,
		}
	case MINOR:
		return Tag{
			Major: previousTag.Major,
			Minor: previousTag.Minor + 1,
			Patch: 0,
		}
	case MAJOR:
		return Tag{
			Major: previousTag.Major + 1,
			Minor: 0,
			Patch: 0,
		}
	}

	return Tag{}
}

func parseStringTag(tag string) (Tag, error) {
	semverRe := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)((?:\-[0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*)?)((?:\+[0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*)?)$`)
	prereleaseRe := regexp.MustCompile(`^\-(alpha|beta|rc)(?:\.(\d+))?$`)

	if !semverRe.MatchString(tag) {
		return Tag{}, fmt.Errorf("%s is not a valid semver tag", tag)
	}

	extractedTag := Tag{}

	matches := semverRe.FindStringSubmatch(tag)

	for i := 1; i < 4; i++ {
		if n, err := strconv.Atoi(matches[i]); err == nil {
			switch i {
			case 1:
				extractedTag.Major = n
			case 2:
				extractedTag.Minor = n
			case 3:
				extractedTag.Patch = n
			}

			continue
		}

		return Tag{}, fmt.Errorf("%s is not a valid integer", matches[i])
	}

	if prereleaseRe.MatchString(matches[4]) {
		prereleseMatches := prereleaseRe.FindStringSubmatch(matches[4])
		var n int

		if len(prereleseMatches) == 3 {
			n, _ = strconv.Atoi(prereleseMatches[2])
		}

		switch prereleseMatches[1] {
		case "alpha":
			extractedTag.Alpha = &n
		case "beta":
			extractedTag.Beta = &n
		case "rc":
			extractedTag.RC = &n
		}
	}

	return extractedTag, nil
}
