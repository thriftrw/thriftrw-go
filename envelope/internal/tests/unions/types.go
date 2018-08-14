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

package unions

import (
	"encoding/base64"
	"fmt"
	"go.uber.org/thriftrw/envelope/internal/tests/typedefs"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/zap/zapcore"
	"strings"
)

// ArbitraryValue allows constructing complex values without a schema.
//
// A value is one of,
//
// * Boolean
// * Integer
// * String
// * A list of other values
// * A dictionary of other values
type ArbitraryValue struct {
	BoolValue   *bool                      `json:"boolValue,omitempty"`
	Int64Value  *int64                     `json:"int64Value,omitempty"`
	StringValue *string                    `json:"stringValue,omitempty"`
	ListValue   []*ArbitraryValue          `json:"listValue,omitempty"`
	MapValue    map[string]*ArbitraryValue `json:"mapValue,omitempty"`
}

type _List_ArbitraryValue_ValueList []*ArbitraryValue

func (v _List_ArbitraryValue_ValueList) ForEach(f func(wire.Value) error) error {
	for i, x := range v {
		if x == nil {
			return fmt.Errorf("invalid [%v]: value is nil", i)
		}
		w, err := x.ToWire()
		if err != nil {
			return err
		}
		err = f(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v _List_ArbitraryValue_ValueList) Size() int {
	return len(v)
}

func (_List_ArbitraryValue_ValueList) ValueType() wire.Type {
	return wire.TStruct
}

func (_List_ArbitraryValue_ValueList) Close() {}

type _Map_String_ArbitraryValue_MapItemList map[string]*ArbitraryValue

func (m _Map_String_ArbitraryValue_MapItemList) ForEach(f func(wire.MapItem) error) error {
	for k, v := range m {
		if v == nil {
			return fmt.Errorf("invalid [%v]: value is nil", k)
		}
		kw, err := wire.NewValueString(k), error(nil)
		if err != nil {
			return err
		}

		vw, err := v.ToWire()
		if err != nil {
			return err
		}
		err = f(wire.MapItem{Key: kw, Value: vw})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m _Map_String_ArbitraryValue_MapItemList) Size() int {
	return len(m)
}

func (_Map_String_ArbitraryValue_MapItemList) KeyType() wire.Type {
	return wire.TBinary
}

func (_Map_String_ArbitraryValue_MapItemList) ValueType() wire.Type {
	return wire.TStruct
}

func (_Map_String_ArbitraryValue_MapItemList) Close() {}

// ToWire translates a ArbitraryValue struct into a Thrift-level intermediate
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
func (v *ArbitraryValue) ToWire() (wire.Value, error) {
	var (
		fields [5]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.BoolValue != nil {
		w, err = wire.NewValueBool(*(v.BoolValue)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	if v.Int64Value != nil {
		w, err = wire.NewValueI64(*(v.Int64Value)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}
	if v.StringValue != nil {
		w, err = wire.NewValueString(*(v.StringValue)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 3, Value: w}
		i++
	}
	if v.ListValue != nil {
		w, err = wire.NewValueList(_List_ArbitraryValue_ValueList(v.ListValue)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 4, Value: w}
		i++
	}
	if v.MapValue != nil {
		w, err = wire.NewValueMap(_Map_String_ArbitraryValue_MapItemList(v.MapValue)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 5, Value: w}
		i++
	}

	if i != 1 {
		return wire.Value{}, fmt.Errorf("ArbitraryValue should have exactly one field: got %v fields", i)
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _ArbitraryValue_Read(w wire.Value) (*ArbitraryValue, error) {
	var v ArbitraryValue
	err := v.FromWire(w)
	return &v, err
}

func _List_ArbitraryValue_Read(l wire.ValueList) ([]*ArbitraryValue, error) {
	if l.ValueType() != wire.TStruct {
		return nil, nil
	}

	o := make([]*ArbitraryValue, 0, l.Size())
	err := l.ForEach(func(x wire.Value) error {
		i, err := _ArbitraryValue_Read(x)
		if err != nil {
			return err
		}
		o = append(o, i)
		return nil
	})
	l.Close()
	return o, err
}

func _Map_String_ArbitraryValue_Read(m wire.MapItemList) (map[string]*ArbitraryValue, error) {
	if m.KeyType() != wire.TBinary {
		return nil, nil
	}

	if m.ValueType() != wire.TStruct {
		return nil, nil
	}

	o := make(map[string]*ArbitraryValue, m.Size())
	err := m.ForEach(func(x wire.MapItem) error {
		k, err := x.Key.GetString(), error(nil)
		if err != nil {
			return err
		}

		v, err := _ArbitraryValue_Read(x.Value)
		if err != nil {
			return err
		}

		o[k] = v
		return nil
	})
	m.Close()
	return o, err
}

// FromWire deserializes a ArbitraryValue struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a ArbitraryValue struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v ArbitraryValue
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *ArbitraryValue) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBool {
				var x bool
				x, err = field.Value.GetBool(), error(nil)
				v.BoolValue = &x
				if err != nil {
					return err
				}

			}
		case 2:
			if field.Value.Type() == wire.TI64 {
				var x int64
				x, err = field.Value.GetI64(), error(nil)
				v.Int64Value = &x
				if err != nil {
					return err
				}

			}
		case 3:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.StringValue = &x
				if err != nil {
					return err
				}

			}
		case 4:
			if field.Value.Type() == wire.TList {
				v.ListValue, err = _List_ArbitraryValue_Read(field.Value.GetList())
				if err != nil {
					return err
				}

			}
		case 5:
			if field.Value.Type() == wire.TMap {
				v.MapValue, err = _Map_String_ArbitraryValue_Read(field.Value.GetMap())
				if err != nil {
					return err
				}

			}
		}
	}

	count := 0
	if v.BoolValue != nil {
		count++
	}
	if v.Int64Value != nil {
		count++
	}
	if v.StringValue != nil {
		count++
	}
	if v.ListValue != nil {
		count++
	}
	if v.MapValue != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("ArbitraryValue should have exactly one field: got %v fields", count)
	}

	return nil
}

// String returns a readable string representation of a ArbitraryValue
// struct.
func (v *ArbitraryValue) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [5]string
	i := 0
	if v.BoolValue != nil {
		fields[i] = fmt.Sprintf("BoolValue: %v", *(v.BoolValue))
		i++
	}
	if v.Int64Value != nil {
		fields[i] = fmt.Sprintf("Int64Value: %v", *(v.Int64Value))
		i++
	}
	if v.StringValue != nil {
		fields[i] = fmt.Sprintf("StringValue: %v", *(v.StringValue))
		i++
	}
	if v.ListValue != nil {
		fields[i] = fmt.Sprintf("ListValue: %v", v.ListValue)
		i++
	}
	if v.MapValue != nil {
		fields[i] = fmt.Sprintf("MapValue: %v", v.MapValue)
		i++
	}

	return fmt.Sprintf("ArbitraryValue{%v}", strings.Join(fields[:i], ", "))
}

func _Bool_EqualsPtr(lhs, rhs *bool) bool {
	if lhs != nil && rhs != nil {

		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

func _I64_EqualsPtr(lhs, rhs *int64) bool {
	if lhs != nil && rhs != nil {

		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

func _String_EqualsPtr(lhs, rhs *string) bool {
	if lhs != nil && rhs != nil {

		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

func _List_ArbitraryValue_Equals(lhs, rhs []*ArbitraryValue) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i, lv := range lhs {
		rv := rhs[i]
		if !lv.Equals(rv) {
			return false
		}
	}

	return true
}

func _Map_String_ArbitraryValue_Equals(lhs, rhs map[string]*ArbitraryValue) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for lk, lv := range lhs {
		rv, ok := rhs[lk]
		if !ok {
			return false
		}
		if !lv.Equals(rv) {
			return false
		}
	}
	return true
}

// Equals returns true if all the fields of this ArbitraryValue match the
// provided ArbitraryValue.
//
// This function performs a deep comparison.
func (v *ArbitraryValue) Equals(rhs *ArbitraryValue) bool {
	if !_Bool_EqualsPtr(v.BoolValue, rhs.BoolValue) {
		return false
	}
	if !_I64_EqualsPtr(v.Int64Value, rhs.Int64Value) {
		return false
	}
	if !_String_EqualsPtr(v.StringValue, rhs.StringValue) {
		return false
	}
	if !((v.ListValue == nil && rhs.ListValue == nil) || (v.ListValue != nil && rhs.ListValue != nil && _List_ArbitraryValue_Equals(v.ListValue, rhs.ListValue))) {
		return false
	}
	if !((v.MapValue == nil && rhs.MapValue == nil) || (v.MapValue != nil && rhs.MapValue != nil && _Map_String_ArbitraryValue_Equals(v.MapValue, rhs.MapValue))) {
		return false
	}

	return true
}

type _List_ArbitraryValue_Zapper []*ArbitraryValue

// MarshalLogArray implements zapcore.ArrayMarshaler, enabling
// fast logging of _List_ArbitraryValue_Zapper.
func (l _List_ArbitraryValue_Zapper) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range l {
		if err := enc.AppendObject(v); err != nil {
			return err
		}
	}
	return nil
}

type _Map_String_ArbitraryValue_Zapper map[string]*ArbitraryValue

// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
// fast logging of _Map_String_ArbitraryValue_Zapper.
func (m _Map_String_ArbitraryValue_Zapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range m {
		if err := enc.AddObject((string)(k), v); err != nil {
			return err
		}
	}
	return nil
}

// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
// fast logging of ArbitraryValue.
func (v *ArbitraryValue) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if v.BoolValue != nil {
		enc.AddBool("boolValue", *v.BoolValue)
	}
	if v.Int64Value != nil {
		enc.AddInt64("int64Value", *v.Int64Value)
	}
	if v.StringValue != nil {
		enc.AddString("stringValue", *v.StringValue)
	}
	if v.ListValue != nil {
		if err := enc.AddArray("listValue", (_List_ArbitraryValue_Zapper)(v.ListValue)); err != nil {
			return err
		}
	}
	if v.MapValue != nil {
		if err := enc.AddObject("mapValue", (_Map_String_ArbitraryValue_Zapper)(v.MapValue)); err != nil {
			return err
		}
	}
	return nil
}

// GetBoolValue returns the value of BoolValue if it is set or its
// zero value if it is unset.
func (v *ArbitraryValue) GetBoolValue() (o bool) {
	if v.BoolValue != nil {
		return *v.BoolValue
	}

	return
}

// GetInt64Value returns the value of Int64Value if it is set or its
// zero value if it is unset.
func (v *ArbitraryValue) GetInt64Value() (o int64) {
	if v.Int64Value != nil {
		return *v.Int64Value
	}

	return
}

// GetStringValue returns the value of StringValue if it is set or its
// zero value if it is unset.
func (v *ArbitraryValue) GetStringValue() (o string) {
	if v.StringValue != nil {
		return *v.StringValue
	}

	return
}

// GetListValue returns the value of ListValue if it is set or its
// zero value if it is unset.
func (v *ArbitraryValue) GetListValue() (o []*ArbitraryValue) {
	if v.ListValue != nil {
		return v.ListValue
	}

	return
}

// GetMapValue returns the value of MapValue if it is set or its
// zero value if it is unset.
func (v *ArbitraryValue) GetMapValue() (o map[string]*ArbitraryValue) {
	if v.MapValue != nil {
		return v.MapValue
	}

	return
}

type Document struct {
	Pdf       typedefs.PDF `json:"pdf,omitempty"`
	PlainText *string      `json:"plainText,omitempty"`
}

// ToWire translates a Document struct into a Thrift-level intermediate
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
func (v *Document) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.Pdf != nil {
		w, err = v.Pdf.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	if v.PlainText != nil {
		w, err = wire.NewValueString(*(v.PlainText)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}

	if i != 1 {
		return wire.Value{}, fmt.Errorf("Document should have exactly one field: got %v fields", i)
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _PDF_Read(w wire.Value) (typedefs.PDF, error) {
	var x typedefs.PDF
	err := x.FromWire(w)
	return x, err
}

// FromWire deserializes a Document struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a Document struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v Document
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *Document) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.Pdf, err = _PDF_Read(field.Value)
				if err != nil {
					return err
				}

			}
		case 2:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.PlainText = &x
				if err != nil {
					return err
				}

			}
		}
	}

	count := 0
	if v.Pdf != nil {
		count++
	}
	if v.PlainText != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("Document should have exactly one field: got %v fields", count)
	}

	return nil
}

// String returns a readable string representation of a Document
// struct.
func (v *Document) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [2]string
	i := 0
	if v.Pdf != nil {
		fields[i] = fmt.Sprintf("Pdf: %v", v.Pdf)
		i++
	}
	if v.PlainText != nil {
		fields[i] = fmt.Sprintf("PlainText: %v", *(v.PlainText))
		i++
	}

	return fmt.Sprintf("Document{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this Document match the
// provided Document.
//
// This function performs a deep comparison.
func (v *Document) Equals(rhs *Document) bool {
	if !((v.Pdf == nil && rhs.Pdf == nil) || (v.Pdf != nil && rhs.Pdf != nil && v.Pdf.Equals(rhs.Pdf))) {
		return false
	}
	if !_String_EqualsPtr(v.PlainText, rhs.PlainText) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
// fast logging of Document.
func (v *Document) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if v.Pdf != nil {
		enc.AddString("pdf", base64.StdEncoding.EncodeToString(([]byte)(v.Pdf)))
	}
	if v.PlainText != nil {
		enc.AddString("plainText", *v.PlainText)
	}
	return nil
}

// GetPdf returns the value of Pdf if it is set or its
// zero value if it is unset.
func (v *Document) GetPdf() (o typedefs.PDF) {
	if v.Pdf != nil {
		return v.Pdf
	}

	return
}

// GetPlainText returns the value of PlainText if it is set or its
// zero value if it is unset.
func (v *Document) GetPlainText() (o string) {
	if v.PlainText != nil {
		return *v.PlainText
	}

	return
}

type EmptyUnion struct {
}

// ToWire translates a EmptyUnion struct into a Thrift-level intermediate
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
func (v *EmptyUnion) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a EmptyUnion struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a EmptyUnion struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v EmptyUnion
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *EmptyUnion) FromWire(w wire.Value) error {

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}

	return nil
}

// String returns a readable string representation of a EmptyUnion
// struct.
func (v *EmptyUnion) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [0]string
	i := 0

	return fmt.Sprintf("EmptyUnion{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this EmptyUnion match the
// provided EmptyUnion.
//
// This function performs a deep comparison.
func (v *EmptyUnion) Equals(rhs *EmptyUnion) bool {

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
// fast logging of EmptyUnion.
func (v *EmptyUnion) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}
