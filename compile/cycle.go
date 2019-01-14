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

// findTypeCycles look for invalid type reference cycles in the given
// TypeSpec.
func findTypeCycles(t TypeSpec) error {
	return make(typeCycleFinder, 0).Visit(t)
}

type typeCycleFinder []TypeSpec

// visited returns true if the given TypeSpec has already been visited.
// Otherwise it returns false, and marks the TypeSpec as visited.
func (f typeCycleFinder) visited(s TypeSpec) bool {
	for _, t := range f {
		if t == s {
			return true
		}
	}
	return false
}

// cloneWithPart creates a copy of this typeCycleFinder with the given
// TypeSpec added to the chain.
func (f typeCycleFinder) cloneWithPart(s TypeSpec) typeCycleFinder {
	newf := make(typeCycleFinder, 0, len(f)+1)
	newf = append(newf, f...)
	newf = append(newf, s)
	return newf
}

func (f typeCycleFinder) Visit(s TypeSpec) error {
	_, isTypedef := s.(*TypedefSpec)

	if f.visited(s) {
		// cycles are errors only for typedefs
		if isTypedef {
			return typeReferenceCycleError{Nodes: append(f, s)}
		}
		return nil
	}

	if _, isStruct := s.(*StructSpec); isStruct {
		// Cycles break at structs. A typedef for a struct which refreences back
		// to the typedef is valid because we support self-referential structs.
		return nil
	}

	return s.ForEachTypeReference(f.cloneWithPart(s).Visit)
}
