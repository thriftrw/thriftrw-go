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

package concurrent

import (
	"log"
	"reflect"
	"sync"

	"go.uber.org/multierr"
)

var (
	_typeOfError = reflect.TypeOf((*error)(nil)).Elem()
	_typeOfInt   = reflect.TypeOf(int(0))
)

// Range calls the function fn on all items in coll concurrently and waits for
// all calls to finish.
//
// coll may be a slice or a map. If coll is a map, fn must accept the key and
// the value as its arguments, and otherwise it must accept an int index and
// the value as its arguments.
//
// fn may return nothing or error.
func Range(coll, fn interface{}) error {
	if coll == nil || fn == nil {
		log.Panicf("ConcurrentRange(%T, %T): both arguments must be non-nil", coll, fn)
	}

	cv := reflect.ValueOf(coll)
	fv := reflect.ValueOf(fn)
	ct := cv.Type()
	ft := fv.Type()

	if ft.NumIn() != 2 {
		log.Panicf("ConcurrentRange(%T, %T): fn must accept exactly two arguments", coll, fn)
	}

	switch ft.NumOut() {
	case 0:
		// adapt into a function that always returns a nil error
		fv = alwaysReturnNoError(fv)
		ft = fv.Type()
	case 1:
		if ft.Out(0) != _typeOfError {
			log.Panicf("ConcurrentRange(%T, %T): fn may only return error or nothing", coll, fn)
		}
	case 2:
		log.Panicf("ConcurrentRange(%T, %T): fn may only return error or nothing", coll, fn)
	}

	var (
		wg     sync.WaitGroup
		lock   sync.Mutex
		errors []error
	)

	switch ct.Kind() {
	case reflect.Map:
		if ft.In(0) != ct.Key() {
			log.Panicf("ConcurrentRange(%T, %T): fn's first argument must be a %v", coll, fn, ct.Key())
		}

		if ft.In(1) != ct.Elem() {
			log.Panicf("ConcurrentRange(%T, %T): fn's second argument must be a %v", coll, fn, ct.Elem())
		}

		for _, key := range cv.MapKeys() {
			value := cv.MapIndex(key)
			wg.Add(1)
			go func(key, value reflect.Value) {
				defer wg.Done()
				err, ok := fv.Call([]reflect.Value{key, value})[0].Interface().(error)
				if ok && err != nil {
					lock.Lock()
					errors = append(errors, err)
					lock.Unlock()
				}
			}(key, value)
		}

	case reflect.Slice:
		if ft.In(0) != _typeOfInt {
			log.Panicf("ConcurrentRange(%T, %T): fn's first argument must be an int", coll, fn)
		}

		if ft.In(1) != ct.Elem() {
			log.Panicf("ConcurrentRange(%T, %T): fn's second argument must be a %v", coll, fn, ct.Elem())
		}

		for i := 0; i < cv.Len(); i++ {
			value := cv.Index(i)
			wg.Add(1)
			go func(key, value reflect.Value) {
				defer wg.Done()
				err, ok := fv.Call([]reflect.Value{key, value})[0].Interface().(error)
				if ok && err != nil {
					lock.Lock()
					errors = append(errors, err)
					lock.Unlock()
				}
			}(reflect.ValueOf(i), value)
		}

	default:
		log.Panicf("ConcurrentRange(%T, %T): called with a type that is not a slice or a map", coll, fn)
	}

	wg.Wait()
	return multierr.Combine(errors...)
}

func alwaysReturnNoError(fn reflect.Value) reflect.Value {
	var (
		ft       = fn.Type()
		in       []reflect.Type
		variadic = ft.IsVariadic()
	)

	for i := 0; i < ft.NumIn(); i++ {
		in = append(in, ft.In(i))
	}

	newFt := reflect.FuncOf(in, []reflect.Type{_typeOfError}, variadic)
	return reflect.MakeFunc(newFt, func(args []reflect.Value) []reflect.Value {
		fn.Call(args)
		return []reflect.Value{reflect.Zero(_typeOfError)}
	})
}
