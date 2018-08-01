// Code generated by thriftrw v1.13.0. DO NOT EDIT.
// @generated

package services

import (
	"errors"
	"fmt"
	"go.uber.org/thriftrw/gen/internal/tests/unions"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/zap/zapcore"
	"strings"
)

// KeyValue_SetValueV2_Args represents the arguments for the KeyValue.setValueV2 function.
//
// The arguments for setValueV2 are sent and received over the wire as this struct.
type KeyValue_SetValueV2_Args struct {
	// Key to change.
	Key Key `json:"key,required"`
	// New value for the key.
	//
	// If the key already has an existing value, it will be overwritten.
	Value *unions.ArbitraryValue `json:"value,required"`
}

// ToWire translates a KeyValue_SetValueV2_Args struct into a Thrift-level intermediate
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
func (v *KeyValue_SetValueV2_Args) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	w, err = v.Key.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	if v.Value == nil {
		return w, errors.New("field Value of KeyValue_SetValueV2_Args is required")
	}
	w, err = v.Value.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a KeyValue_SetValueV2_Args struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a KeyValue_SetValueV2_Args struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v KeyValue_SetValueV2_Args
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *KeyValue_SetValueV2_Args) FromWire(w wire.Value) error {
	var err error

	keyIsSet := false
	valueIsSet := false

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.Key, err = _Key_Read(field.Value)
				if err != nil {
					return err
				}
				keyIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TStruct {
				v.Value, err = _ArbitraryValue_Read(field.Value)
				if err != nil {
					return err
				}
				valueIsSet = true
			}
		}
	}

	if !keyIsSet {
		return errors.New("field Key of KeyValue_SetValueV2_Args is required")
	}

	if !valueIsSet {
		return errors.New("field Value of KeyValue_SetValueV2_Args is required")
	}

	return nil
}

// String returns a readable string representation of a KeyValue_SetValueV2_Args
// struct.
func (v *KeyValue_SetValueV2_Args) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("Key: %v", v.Key)
	i++
	fields[i] = fmt.Sprintf("Value: %v", v.Value)
	i++

	return fmt.Sprintf("KeyValue_SetValueV2_Args{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this KeyValue_SetValueV2_Args match the
// provided KeyValue_SetValueV2_Args.
//
// This function performs a deep comparison.
func (v *KeyValue_SetValueV2_Args) Equals(rhs *KeyValue_SetValueV2_Args) bool {
	if !(v.Key == rhs.Key) {
		return false
	}
	if !v.Value.Equals(rhs.Value) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, allowing
// fast logging of KeyValue_SetValueV2_Args.
func (v *KeyValue_SetValueV2_Args) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("key", (string)(v.Key))
	if err := enc.AddObject("value", v.Value); err != nil {
		return err
	}

	return nil
}

// GetKey returns the value of Key if it is set or its
// zero value if it is unset.
func (v *KeyValue_SetValueV2_Args) GetKey() (o Key) { return v.Key }

// GetValue returns the value of Value if it is set or its
// zero value if it is unset.
func (v *KeyValue_SetValueV2_Args) GetValue() (o *unions.ArbitraryValue) { return v.Value }

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the arguments.
//
// This will always be "setValueV2" for this struct.
func (v *KeyValue_SetValueV2_Args) MethodName() string {
	return "setValueV2"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Call for this struct.
func (v *KeyValue_SetValueV2_Args) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

// KeyValue_SetValueV2_Helper provides functions that aid in handling the
// parameters and return values of the KeyValue.setValueV2
// function.
var KeyValue_SetValueV2_Helper = struct {
	// Args accepts the parameters of setValueV2 in-order and returns
	// the arguments struct for the function.
	Args func(
		key Key,
		value *unions.ArbitraryValue,
	) *KeyValue_SetValueV2_Args

	// IsException returns true if the given error can be thrown
	// by setValueV2.
	//
	// An error can be thrown by setValueV2 only if the
	// corresponding exception type was mentioned in the 'throws'
	// section for it in the Thrift file.
	IsException func(error) bool

	// WrapResponse returns the result struct for setValueV2
	// given the error returned by it. The provided error may
	// be nil if setValueV2 did not fail.
	//
	// This allows mapping errors returned by setValueV2 into a
	// serializable result struct. WrapResponse returns a
	// non-nil error if the provided error cannot be thrown by
	// setValueV2
	//
	//   err := setValueV2(args)
	//   result, err := KeyValue_SetValueV2_Helper.WrapResponse(err)
	//   if err != nil {
	//     return fmt.Errorf("unexpected error from setValueV2: %v", err)
	//   }
	//   serialize(result)
	WrapResponse func(error) (*KeyValue_SetValueV2_Result, error)

	// UnwrapResponse takes the result struct for setValueV2
	// and returns the erorr returned by it (if any).
	//
	// The error is non-nil only if setValueV2 threw an
	// exception.
	//
	//   result := deserialize(bytes)
	//   err := KeyValue_SetValueV2_Helper.UnwrapResponse(result)
	UnwrapResponse func(*KeyValue_SetValueV2_Result) error
}{}

func init() {
	KeyValue_SetValueV2_Helper.Args = func(
		key Key,
		value *unions.ArbitraryValue,
	) *KeyValue_SetValueV2_Args {
		return &KeyValue_SetValueV2_Args{
			Key:   key,
			Value: value,
		}
	}

	KeyValue_SetValueV2_Helper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}

	KeyValue_SetValueV2_Helper.WrapResponse = func(err error) (*KeyValue_SetValueV2_Result, error) {
		if err == nil {
			return &KeyValue_SetValueV2_Result{}, nil
		}

		return nil, err
	}
	KeyValue_SetValueV2_Helper.UnwrapResponse = func(result *KeyValue_SetValueV2_Result) (err error) {
		return
	}

}

// KeyValue_SetValueV2_Result represents the result of a KeyValue.setValueV2 function call.
//
// The result of a setValueV2 execution is sent and received over the wire as this struct.
type KeyValue_SetValueV2_Result struct {
}

// ToWire translates a KeyValue_SetValueV2_Result struct into a Thrift-level intermediate
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
func (v *KeyValue_SetValueV2_Result) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a KeyValue_SetValueV2_Result struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a KeyValue_SetValueV2_Result struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v KeyValue_SetValueV2_Result
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *KeyValue_SetValueV2_Result) FromWire(w wire.Value) error {

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}

	return nil
}

// String returns a readable string representation of a KeyValue_SetValueV2_Result
// struct.
func (v *KeyValue_SetValueV2_Result) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [0]string
	i := 0

	return fmt.Sprintf("KeyValue_SetValueV2_Result{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this KeyValue_SetValueV2_Result match the
// provided KeyValue_SetValueV2_Result.
//
// This function performs a deep comparison.
func (v *KeyValue_SetValueV2_Result) Equals(rhs *KeyValue_SetValueV2_Result) bool {

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, allowing
// fast logging of KeyValue_SetValueV2_Result.
func (v *KeyValue_SetValueV2_Result) MarshalLogObject(enc zapcore.ObjectEncoder) error {

	return nil
}

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the result.
//
// This will always be "setValueV2" for this struct.
func (v *KeyValue_SetValueV2_Result) MethodName() string {
	return "setValueV2"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Reply for this struct.
func (v *KeyValue_SetValueV2_Result) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}
