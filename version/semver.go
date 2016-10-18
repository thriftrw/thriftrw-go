// Copyright (c) 2016 Uber Technologies, Inc.
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

package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type semVer struct {
	Major uint
	Minor uint
	Patch uint
	Pre   []string
	Meta  string
}

// Following http://semver.org/spec/v2.0.0.html
const numPart = `([0-9]+)\.([0-9]+)\.([0-9]+)`
const preReleasePart = `-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)`
const metaPart = `\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)`

var semVerRegex = regexp.MustCompile(`^` + numPart + `(?:` + preReleasePart + `)?(?:` + metaPart + `)?$`)

func parseSemVer(v string) (semVer, error) {
	r := semVer{}
	parts := semVerRegex.FindStringSubmatch(v)
	if parts == nil {
		return r, fmt.Errorf(`cannot parse as semantic version: "%s"`, v)
	}

	var err error
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
	var v uint64
	var err error
	if v, err = strconv.ParseUint(s, 10, 31); err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (v *semVer) String() string {
	r := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != nil {
		r = fmt.Sprintf("%s-%s", r, strings.Join(v.Pre, "."))
	}
	if v.Meta != "" {
		r = fmt.Sprintf("%s+%s", r, v.Meta)
	}
	return r
}

// Compare returns:
//  0 if a == b
// -1 if a < b
// +1 if a > b
func (v *semVer) Compare(b *semVer) int {
	a := v
	r := uintCmp(a.Major, b.Major)
	if r != 0 {
		return r
	}
	r = uintCmp(a.Minor, b.Minor)
	if r != 0 {
		return r
	}
	r = uintCmp(a.Patch, b.Patch)
	if r != 0 {
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
		r = preCmp(aPre[0], bPre[0])
		if r != 0 {
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
