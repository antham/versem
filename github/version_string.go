// Code generated by "stringer -type=Version"; DO NOT EDIT.

package github

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UNVALID_VERSION-0]
	_ = x[NORELEASE-1]
	_ = x[ALPHA-2]
	_ = x[BETA-3]
	_ = x[RC-4]
	_ = x[PATCH-5]
	_ = x[MINOR-6]
	_ = x[MAJOR-7]
}

const _Version_name = "UNVALID_VERSIONNORELEASEALPHABETARCPATCHMINORMAJOR"

var _Version_index = [...]uint8{0, 15, 24, 29, 33, 35, 40, 45, 50}

func (i Version) String() string {
	if i < 0 || i >= Version(len(_Version_index)-1) {
		return "Version(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Version_name[_Version_index[i]:_Version_index[i+1]]
}
