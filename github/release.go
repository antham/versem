package github

import (
	"context"
	"fmt"

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

	switch true {
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
