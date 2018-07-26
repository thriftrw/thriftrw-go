// Code generated by thriftrw v1.13.0. DO NOT EDIT.
// @generated

// Copyright (c) 2018 Uber Technologies, Inc.
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

package exceptions

import (
	"errors"
	"fmt"
	"go.uber.org/thriftrw/wire"
	"strings"
)

// Raised when something doesn't exist.
type DoesNotExistException struct {
	// Key that was missing.
	Key    string  `json:"key,required"`
	Error2 *string `json:"Error,omitempty"`
}

// ToWire translates a DoesNotExistException struct into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
//
// An error is returned if the struct or any of its fields failed to
// validate.
//
//   x, err := v.ToWire()
//   if err != nil {
//     return err
//   }
//
//   if err := binaryProtocol.Encode(x, writer); err != nil {
//     return err
//   }
func (v *DoesNotExistException) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	w, err = wire.NewValueString(v.Key), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	if v.Error2 != nil {
		w, err = wire.NewValueString(*(v.Error2)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a DoesNotExistException struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a DoesNotExistException struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v DoesNotExistException
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *DoesNotExistException) FromWire(w wire.Value) error {
	var err error

	keyIsSet := false

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.Key, err = field.Value.GetString(), error(nil)
				if err != nil {
					return err
				}
				keyIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.Error2 = &x
				if err != nil {
					return err
				}

			}
		}
	}

	if !keyIsSet {
		return errors.New("field Key of DoesNotExistException is required")
	}

	return nil
}

// String returns a readable string representation of a DoesNotExistException
// struct.
func (v *DoesNotExistException) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("Key: %v", v.Key)
	i++
	if v.Error2 != nil {
		fields[i] = fmt.Sprintf("Error2: %v", *(v.Error2))
		i++
	}

	return fmt.Sprintf("DoesNotExistException{%v}", strings.Join(fields[:i], ", "))
}

func _String_EqualsPtr(lhs, rhs *string) bool {
	if lhs != nil && rhs != nil {

		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

// Equals returns true if all the fields of this DoesNotExistException match the
// provided DoesNotExistException.
//
// This function performs a deep comparison.
func (v *DoesNotExistException) Equals(rhs *DoesNotExistException) bool {
	if !(v.Key == rhs.Key) {
		return false
	}
	if !_String_EqualsPtr(v.Error2, rhs.Error2) {
		return false
	}

	return true
}

// GetKey returns the value of Key if it is set or its
// zero value if it is unset.
func (v *DoesNotExistException) GetKey() (o string) { return v.Key }

// GetError2 returns the value of Error2 if it is set or its
// zero value if it is unset.
func (v *DoesNotExistException) GetError2() (o string) {
	if v.Error2 != nil {
		return *v.Error2
	}

	return
}

func (v *DoesNotExistException) Error() string {
	return v.String()
}

type EmptyException struct {
}

// ToWire translates a EmptyException struct into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
//
// An error is returned if the struct or any of its fields failed to
// validate.
//
//   x, err := v.ToWire()
//   if err != nil {
//     return err
//   }
//
//   if err := binaryProtocol.Encode(x, writer); err != nil {
//     return err
//   }
func (v *EmptyException) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a EmptyException struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a EmptyException struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v EmptyException
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *EmptyException) FromWire(w wire.Value) error {

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}

	return nil
}

// String returns a readable string representation of a EmptyException
// struct.
func (v *EmptyException) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [0]string
	i := 0

	return fmt.Sprintf("EmptyException{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this EmptyException match the
// provided EmptyException.
//
// This function performs a deep comparison.
func (v *EmptyException) Equals(rhs *EmptyException) bool {

	return true
}

func (v *EmptyException) Error() string {
	return v.String()
}
