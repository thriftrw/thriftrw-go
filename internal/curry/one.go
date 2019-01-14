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

package curry

import (
	"fmt"
	"reflect"
)

// One wraps the given function to require one less argument, passing in the
// given value instead. The value MUST NOT be nil.
//
// BUG(abg): Variadic functions are not supported.
func One(f interface{}, a interface{}) interface{} {
	if f == nil {
		panic("f is required")
	}
	if a == nil {
		panic("a is required and cannot be nil")
	}

	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		panic(fmt.Sprintf("%v (%v) is not a function", f, fType))
	}

	if fType.IsVariadic() {
		panic(fmt.Sprintf("%v (%v) is variadic", f, fType))
	}

	in := args(fType)
	out := returns(fType)
	if len(in) == 0 {
		panic(fmt.Sprintf("%v (%v) does not accept enough arguments to curry in %v", f, fType, a))
	}

	fVal := reflect.ValueOf(f)
	aVal := reflect.ValueOf(a)
	newFType := reflect.FuncOf(in[1:], out, false)
	return reflect.MakeFunc(newFType, func(args []reflect.Value) []reflect.Value {
		newArgs := make([]reflect.Value, 0, fType.NumIn())
		newArgs = append(newArgs, aVal)
		newArgs = append(newArgs, args...)
		return fVal.Call(newArgs)
	}).Interface()
}

func args(f reflect.Type) []reflect.Type {
	in := make([]reflect.Type, f.NumIn())
	for i := 0; i < f.NumIn(); i++ {
		in[i] = f.In(i)
	}
	return in
}

func returns(f reflect.Type) []reflect.Type {
	out := make([]reflect.Type, f.NumOut())
	for i := 0; i < f.NumOut(); i++ {
		out[i] = f.Out(i)
	}
	return out
}
