// Code generated by thriftrw v1.34.0. DO NOT EDIT.
// @generated

package enum_text_marshal_strict

import (
	bytes "bytes"
	json "encoding/json"
	fmt "fmt"
	stream "go.uber.org/thriftrw/protocol/stream"
	thriftreflect "go.uber.org/thriftrw/thriftreflect"
	wire "go.uber.org/thriftrw/wire"
	zapcore "go.uber.org/zap/zapcore"
	math "math"
	strconv "strconv"
)

type EnumMarshalStrict int32

const (
	EnumMarshalStrictFoo EnumMarshalStrict = 0
	EnumMarshalStrictBar EnumMarshalStrict = 1
	EnumMarshalStrictBaz EnumMarshalStrict = 2
	EnumMarshalStrictBat EnumMarshalStrict = 3
)

// EnumMarshalStrict_Values returns all recognized values of EnumMarshalStrict.
func EnumMarshalStrict_Values() []EnumMarshalStrict {
	return []EnumMarshalStrict{
		EnumMarshalStrictFoo,
		EnumMarshalStrictBar,
		EnumMarshalStrictBaz,
		EnumMarshalStrictBat,
	}
}

// UnmarshalText tries to decode EnumMarshalStrict from a byte slice
// containing its name.
//
//	var v EnumMarshalStrict
//	err := v.UnmarshalText([]byte("Foo"))
func (v *EnumMarshalStrict) UnmarshalText(value []byte) error {
	switch s := string(value); s {
	case "Foo":
		*v = EnumMarshalStrictFoo
		return nil
	case "Bar":
		*v = EnumMarshalStrictBar
		return nil
	case "Baz":
		*v = EnumMarshalStrictBaz
		return nil
	case "Bat":
		*v = EnumMarshalStrictBat
		return nil
	default:
		val, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("unknown enum value %q for %q: %v", s, "EnumMarshalStrict", err)
		}
		*v = EnumMarshalStrict(val)
		return nil
	}
}

// MarshalText encodes EnumMarshalStrict to text.
//
// If the enum value is recognized, its name is returned.
// Otherwise, an error is returned.
//
// This implements the TextMarshaler interface.
func (v EnumMarshalStrict) MarshalText() ([]byte, error) {
	switch int32(v) {
	case 0:
		return []byte("Foo"), nil
	case 1:
		return []byte("Bar"), nil
	case 2:
		return []byte("Baz"), nil
	case 3:
		return []byte("Bat"), nil
	}
	return nil, fmt.Errorf("unknown enum value %q for %q", v, "EnumMarshalStrict")
}

// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
// fast logging of EnumMarshalStrict.
// Enums are logged as objects, where the value is logged with key "value", and
// if this value's name is known, the name is logged with key "name".
func (v EnumMarshalStrict) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt32("value", int32(v))
	switch int32(v) {
	case 0:
		enc.AddString("name", "Foo")
	case 1:
		enc.AddString("name", "Bar")
	case 2:
		enc.AddString("name", "Baz")
	case 3:
		enc.AddString("name", "Bat")
	}
	return nil
}

// Ptr returns a pointer to this enum value.
func (v EnumMarshalStrict) Ptr() *EnumMarshalStrict {
	return &v
}

// Encode encodes EnumMarshalStrict directly to bytes.
//
//	sWriter := BinaryStreamer.Writer(writer)
//
//	var v EnumMarshalStrict
//	return v.Encode(sWriter)
func (v EnumMarshalStrict) Encode(sw stream.Writer) error {
	return sw.WriteInt32(int32(v))
}

// ToWire translates EnumMarshalStrict into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
//
// Enums are represented as 32-bit integers over the wire.
func (v EnumMarshalStrict) ToWire() (wire.Value, error) {
	return wire.NewValueI32(int32(v)), nil
}

// FromWire deserializes EnumMarshalStrict from its Thrift-level
// representation.
//
//	x, err := binaryProtocol.Decode(reader, wire.TI32)
//	if err != nil {
//	    return EnumMarshalStrict(0), err
//	}
//
//	var v EnumMarshalStrict
//	if err := v.FromWire(x); err != nil {
//	    return EnumMarshalStrict(0), err
//	}
//	return v, nil
func (v *EnumMarshalStrict) FromWire(w wire.Value) error {
	*v = (EnumMarshalStrict)(w.GetI32())
	return nil
}

// Decode reads off the encoded EnumMarshalStrict directly off of the wire.
//
//	sReader := BinaryStreamer.Reader(reader)
//
//	var v EnumMarshalStrict
//	if err := v.Decode(sReader); err != nil {
//	    return EnumMarshalStrict(0), err
//	}
//	return v, nil
func (v *EnumMarshalStrict) Decode(sr stream.Reader) error {
	i, err := sr.ReadInt32()
	if err != nil {
		return err
	}
	*v = (EnumMarshalStrict)(i)
	return nil
}

// String returns a readable string representation of EnumMarshalStrict.
func (v EnumMarshalStrict) String() string {
	w := int32(v)
	switch w {
	case 0:
		return "Foo"
	case 1:
		return "Bar"
	case 2:
		return "Baz"
	case 3:
		return "Bat"
	}
	return fmt.Sprintf("EnumMarshalStrict(%d)", w)
}

// Equals returns true if this EnumMarshalStrict value matches the provided
// value.
func (v EnumMarshalStrict) Equals(rhs EnumMarshalStrict) bool {
	return v == rhs
}

// MarshalJSON serializes EnumMarshalStrict into JSON.
//
// If the enum value is recognized, its name is returned.
// Otherwise, an error is returned.
//
// This implements json.Marshaler.
func (v EnumMarshalStrict) MarshalJSON() ([]byte, error) {
	switch int32(v) {
	case 0:
		return ([]byte)("\"Foo\""), nil
	case 1:
		return ([]byte)("\"Bar\""), nil
	case 2:
		return ([]byte)("\"Baz\""), nil
	case 3:
		return ([]byte)("\"Bat\""), nil
	}
	return nil, fmt.Errorf("unknown enum value %q for %q", v, "EnumMarshalStrict")
}

// UnmarshalJSON attempts to decode EnumMarshalStrict from its JSON
// representation.
//
// This implementation supports both, numeric and string inputs. If a
// string is provided, it must be a known enum name.
//
// This implements json.Unmarshaler.
func (v *EnumMarshalStrict) UnmarshalJSON(text []byte) error {
	d := json.NewDecoder(bytes.NewReader(text))
	d.UseNumber()
	t, err := d.Token()
	if err != nil {
		return err
	}

	switch w := t.(type) {
	case json.Number:
		x, err := w.Int64()
		if err != nil {
			return err
		}
		if x > math.MaxInt32 {
			return fmt.Errorf("enum overflow from JSON %q for %q", text, "EnumMarshalStrict")
		}
		if x < math.MinInt32 {
			return fmt.Errorf("enum underflow from JSON %q for %q", text, "EnumMarshalStrict")
		}
		*v = (EnumMarshalStrict)(x)
		return nil
	case string:
		return v.UnmarshalText([]byte(w))
	default:
		return fmt.Errorf("invalid JSON value %q (%T) to unmarshal into %q", t, t, "EnumMarshalStrict")
	}
}

// ThriftModule represents the IDL file used to generate this package.
var ThriftModule = &thriftreflect.ThriftModule{
	Name:     "enum-text-marshal-strict",
	Package:  "go.uber.org/thriftrw/gen/internal/tests/enum-text-marshal-strict",
	FilePath: "enum-text-marshal-strict.thrift",
	SHA1:     "7d9566a0ff9eccda2ed5be518f321cfcba028dcb",
	Raw:      rawIDL,
}

const rawIDL = "enum EnumMarshalStrict {\n    Foo, Bar, Baz, Bat\n}"
