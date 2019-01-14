// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version is a parsed semantic version representation.
type Version struct {
	Major uint
	Minor uint
	Patch uint
	Pre   []string
	Meta  string
}

// Following http://semver.org/spec/v2.0.0.html
const (
	numPart        = `(\d+)\.(\d+)\.(\d+)`
	preReleasePart = `-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)`
	metaPart       = `\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)`
)

var semVerRegex = regexp.MustCompile(`^` + numPart + `(?:` + preReleasePart + `)?(?:` + metaPart + `)?$`)

// Parse a semantic version string.
func Parse(v string) (r Version, err error) {
	parts := semVerRegex.FindStringSubmatch(v)
	if parts == nil {
		return r, fmt.Errorf(`cannot parse as semantic version: %q`, v)
	}

	if r.Major, err = parseUint(parts[1]); err != nil {
		return r, err
	}
	if r.Minor, err = parseUint(parts[2]); err != nil {
		return r, err
	}
	if r.Patch, err = parseUint(parts[3]); err != nil {
		return r, err
	}

	if parts[4] != "" {
		r.Pre = strings.Split(parts[4], ".")
	}
	r.Meta = parts[5]
	return r, nil
}

func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 31)
	return uint(v), err
}

func (v *Version) String() string {
	r := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if len(v.Pre) > 0 {
		r += "-" + strings.Join(v.Pre, ".")
	}
	if v.Meta != "" {
		r += "+" + v.Meta
	}
	return r
}

// Compare returns:
//  0 if a == b
// -1 if a < b
// +1 if a > b
func (v *Version) Compare(b *Version) int {
	a := v
	if r := uintCmp(a.Major, b.Major); r != 0 {
		return r
	}
	if r := uintCmp(a.Minor, b.Minor); r != 0 {
		return r
	}
	if r := uintCmp(a.Patch, b.Patch); r != 0 {
		return r
	}

	aPre := a.Pre
	bPre := b.Pre

	switch {
	case len(aPre) == 0 && len(bPre) > 0:
		return +1
	case len(aPre) > 0 && len(bPre) == 0:
		return -1
	}

	for len(aPre) > 0 && len(bPre) > 0 {
		if r := preCmp(aPre[0], bPre[0]); r != 0 {
			return r
		}
		aPre = aPre[1:]
		bPre = bPre[1:]
	}
	return uintCmp(uint(len(aPre)), uint(len(bPre)))
}

func uintCmp(a, b uint) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return +1
	default:
		return 0
	}
}

// preCmp compares the pre release items `a` and `b` following semver rules.
func preCmp(a, b string) int {
	auint, auinterr := parseUint(a)
	buint, buinterr := parseUint(b)
	switch {
	case auinterr == nil && buinterr == nil:
		return uintCmp(auint, buint)
	case auinterr == nil && buinterr != nil:
		return -1
	case auinterr != nil && buinterr == nil:
		return +1
	}
	return strings.Compare(a, b)
}
