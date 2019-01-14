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

package compile

import (
	"go.uber.org/thriftrw/ast"
)

// Annotations maps annotations
type Annotations map[string]string

func compileAnnotations(astAnnotations []*ast.Annotation) (Annotations, error) {
	if len(astAnnotations) == 0 {
		return nil, nil
	}

	namespace := newNamespace(caseSensitive)
	annotations := make(Annotations, len(astAnnotations))

	for _, a := range astAnnotations {
		if err := namespace.claim(a.Name, a.Line); err != nil {
			return nil, annotationConflictError{Reason: err}
		}

		annotations[a.Name] = a.Value
	}

	if len(annotations) == 0 {
		return nil, nil
	}
	return annotations, nil
}
