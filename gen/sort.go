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

package gen

import (
	"fmt"
	"reflect"
	"sort"
)

// sortStringKeys returns a sorted list of strings given a map[string]*.
func sortStringKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	t := v.Type()
	if t.Kind() != reflect.Map || t.Key().Kind() != reflect.String {
		panic(fmt.Sprintf(
			"sortStringKeys may be called with a map[string]* only"))
	}

	keys := v.MapKeys()
	sortedKeys := make([]string, 0, len(keys))

	for _, k := range keys {
		key := k.Interface().(string)
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)
	return sortedKeys
}
