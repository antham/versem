package github

import (
	"fmt"
	"regexp"
	"strconv"
)

// Tag represents a semver tag
type Tag struct {
	LeadingV      bool
	Major         int
	Minor         int
	Patch         int
	RC            *int
	Beta          *int
	Alpha         *int
	PreRelease    string
	BuildMetadata string
}

// String converts a Tag to its string representation
func (t Tag) String() string {
	main := fmt.Sprintf("%d.%d.%d", t.Major, t.Minor, t.Patch)

	if t.PreRelease != "" {
		main = fmt.Sprintf("%s-%s", main, t.PreRelease)
	}

	if t.BuildMetadata != "" {
		main = fmt.Sprintf("%s+%s", main, t.BuildMetadata)
	}

	return main
}

// NewTagFromString creates a Tag instance from a tag string representation
func NewTagFromString(tag string) (Tag, error) {
	return parseStringTag(tag)
}
func getNextTag(previousTag Tag, version Version) Tag {
	nextTag := Tag{
		LeadingV: previousTag.LeadingV,
		Major:    previousTag.Major,
		Minor:    previousTag.Minor,
		Patch:    previousTag.Patch,
	}

	switch version {
	case ALPHA:
		var v int
		if previousTag.Alpha != nil {
			v = *previousTag.Alpha + 1
		}
		nextTag.Alpha = &v
		nextTag.PreRelease = "alpha"
		if v > 0 {
			nextTag.PreRelease = fmt.Sprintf("%s.%d", nextTag.PreRelease, v)
		}
	case BETA:
		var v int
		if previousTag.Beta != nil {
			v = *previousTag.Beta + 1
		}
		nextTag.Beta = &v
		nextTag.PreRelease = "beta"
		if v > 0 {
			nextTag.PreRelease = fmt.Sprintf("%s.%d", nextTag.PreRelease, v)
		}
	case RC:
		var v int
		if previousTag.RC != nil {
			v = *previousTag.RC + 1
		}
		nextTag.RC = &v
		nextTag.PreRelease = "rc"
		if v > 0 {
			nextTag.PreRelease = fmt.Sprintf("%s.%d", nextTag.PreRelease, v)
		}
	case PATCH:
		nextTag.Patch = previousTag.Patch + 1
	case MINOR:
		nextTag.Minor = previousTag.Minor + 1
		nextTag.Patch = 0
	case MAJOR:
		nextTag.Major = previousTag.Major + 1
		nextTag.Minor = 0
		nextTag.Patch = 0
	default:
		return Tag{}
	}

	return nextTag
}

func parseStringTag(tag string) (Tag, error) {
	semverRe := regexp.MustCompile(`^(v?)(\d+)\.(\d+)\.(\d+)((?:\-[0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*)?)((?:\+[0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*)?)$`)
	prereleaseRe := regexp.MustCompile(`^\-(alpha|beta|rc)(?:\.(\d+))?$`)

	if !semverRe.MatchString(tag) {
		return Tag{}, fmt.Errorf("%s is not a valid semver tag", tag)
	}

	extractedTag := Tag{}

	matches := semverRe.FindStringSubmatch(tag)

	if matches[1] == "v" {
		extractedTag.LeadingV = true
	}

	for i := 2; i < 5; i++ {
		if n, err := strconv.Atoi(matches[i]); err == nil {
			switch i {
			case 2:
				extractedTag.Major = n
			case 3:
				extractedTag.Minor = n
			case 4:
				extractedTag.Patch = n
			}

			continue
		}

		return Tag{}, fmt.Errorf("%s is not a valid integer", matches[i])
	}

	if prereleaseRe.MatchString(matches[5]) {
		prereleseMatches := prereleaseRe.FindStringSubmatch(matches[5])
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

	if matches[5] != "" {
		extractedTag.PreRelease = matches[5][1:]
	}

	if matches[6] != "" {
		extractedTag.BuildMetadata = matches[6][1:]
	}

	return extractedTag, nil
}
