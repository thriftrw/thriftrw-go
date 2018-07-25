// Code generated by thriftrw v1.13.0. DO NOT EDIT.
// @generated

package services

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/zap/zapcore"
	"strings"
)

type ConflictingNamesSetValueArgs struct {
	Key   string `json:"key,required"`
	Value []byte `json:"value,required"`
}

// ToWire translates a ConflictingNamesSetValueArgs struct into a Thrift-level intermediate
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
func (v *ConflictingNamesSetValueArgs) ToWire() (wire.Value, error) {
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
	if v.Value == nil {
		return w, errors.New("field Value of ConflictingNamesSetValueArgs is required")
	}
	w, err = wire.NewValueBinary(v.Value), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a ConflictingNamesSetValueArgs struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a ConflictingNamesSetValueArgs struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v ConflictingNamesSetValueArgs
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *ConflictingNamesSetValueArgs) FromWire(w wire.Value) error {
	var err error

	keyIsSet := false
	valueIsSet := false

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
				v.Value, err = field.Value.GetBinary(), error(nil)
				if err != nil {
					return err
				}
				valueIsSet = true
			}
		}
	}

	if !keyIsSet {
		return errors.New("field Key of ConflictingNamesSetValueArgs is required")
	}

	if !valueIsSet {
		return errors.New("field Value of ConflictingNamesSetValueArgs is required")
	}

	return nil
}

// String returns a readable string representation of a ConflictingNamesSetValueArgs
// struct.
func (v *ConflictingNamesSetValueArgs) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("Key: %v", v.Key)
	i++
	fields[i] = fmt.Sprintf("Value: %v", v.Value)
	i++

	return fmt.Sprintf("ConflictingNamesSetValueArgs{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this ConflictingNamesSetValueArgs match the
// provided ConflictingNamesSetValueArgs.
//
// This function performs a deep comparison.
func (v *ConflictingNamesSetValueArgs) Equals(rhs *ConflictingNamesSetValueArgs) bool {
	if !(v.Key == rhs.Key) {
		return false
	}
	if !bytes.Equal(v.Value, rhs.Value) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler. (TODO)
func (v *ConflictingNamesSetValueArgs) MarshalLogObject(enc zapcore.ObjectEncoder) {

	enc.AddString("key", v.Key)

	enc.AddBinary("value", v.Value)

}

// GetKey returns the value of Key if it is set or its
// zero value if it is unset.
func (v *ConflictingNamesSetValueArgs) GetKey() (o string) { return v.Key }

// GetValue returns the value of Value if it is set or its
// zero value if it is unset.
func (v *ConflictingNamesSetValueArgs) GetValue() (o []byte) { return v.Value }

type InternalError struct {
	Message *string `json:"message,omitempty"`
}

// ToWire translates a InternalError struct into a Thrift-level intermediate
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
func (v *InternalError) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.Message != nil {
		w, err = wire.NewValueString(*(v.Message)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a InternalError struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a InternalError struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v InternalError
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *InternalError) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.Message = &x
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

// String returns a readable string representation of a InternalError
// struct.
func (v *InternalError) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [1]string
	i := 0
	if v.Message != nil {
		fields[i] = fmt.Sprintf("Message: %v", *(v.Message))
		i++
	}

	return fmt.Sprintf("InternalError{%v}", strings.Join(fields[:i], ", "))
}

func _String_EqualsPtr(lhs, rhs *string) bool {
	if lhs != nil && rhs != nil {

		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

// Equals returns true if all the fields of this InternalError match the
// provided InternalError.
//
// This function performs a deep comparison.
func (v *InternalError) Equals(rhs *InternalError) bool {
	if !_String_EqualsPtr(v.Message, rhs.Message) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler. (TODO)
func (v *InternalError) MarshalLogObject(enc zapcore.ObjectEncoder) {

	if v.Message != nil {
		enc.AddString("message", *v.Message)
	}

}

// GetMessage returns the value of Message if it is set or its
// zero value if it is unset.
func (v *InternalError) GetMessage() (o string) {
	if v.Message != nil {
		return *v.Message
	}

	return
}

func (v *InternalError) Error() string {
	return v.String()
}

type Key string

// ToWire translates Key into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v Key) ToWire() (wire.Value, error) {
	x := (string)(v)
	return wire.NewValueString(x), error(nil)
}

// String returns a readable string representation of Key.
func (v Key) String() string {
	x := (string)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes Key from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *Key) FromWire(w wire.Value) error {
	x, err := w.GetString(), error(nil)
	*v = (Key)(x)
	return err
}

// Equals returns true if this Key is equal to the provided
// Key.
func (lhs Key) Equals(rhs Key) bool {
	return (lhs == rhs)
}
