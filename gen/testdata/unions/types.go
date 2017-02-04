// Code generated by thriftrw v1.1.0
// @generated

package unions

import (
	"fmt"
	"go.uber.org/thriftrw/gen/testdata/typedefs"
	"go.uber.org/thriftrw/wire"
	"strings"
)

type ArbitraryValue struct {
	BoolValue   *bool                      `json:"boolValue,omitempty"`
	Int64Value  *int64                     `json:"int64Value,omitempty"`
	StringValue *string                    `json:"stringValue,omitempty"`
	ListValue   []*ArbitraryValue          `json:"listValue"`
	MapValue    map[string]*ArbitraryValue `json:"mapValue"`
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

func (_List_ArbitraryValue_ValueList) Close() {
}

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

func (_Map_String_ArbitraryValue_MapItemList) Close() {
}

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

func (v *ArbitraryValue) String() string {
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

func _bool_EqualsPtr(lhs, rhs *bool) bool {
	if lhs != nil && rhs != nil {
		x := *lhs
		y := *rhs
		return x == y
	} else if lhs == nil && rhs == nil {
		return true
	} else {
		return false
	}
}

func _i64_EqualsPtr(lhs, rhs *int64) bool {
	if lhs != nil && rhs != nil {
		x := *lhs
		y := *rhs
		return x == y
	} else if lhs == nil && rhs == nil {
		return true
	} else {
		return false
	}
}

func _string_EqualsPtr(lhs, rhs *string) bool {
	if lhs != nil && rhs != nil {
		x := *lhs
		y := *rhs
		return x == y
	} else if lhs == nil && rhs == nil {
		return true
	} else {
		return false
	}
}

func _List_ArbitraryValue_Equals(lhs, rhs []*ArbitraryValue) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, lv := range lhs {
		rv := rhs[i]
		if !(lv.Equals(rv)) {
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
		if !(lv.Equals(rv)) {
			return false
		}
	}
	return true
}

func (lhs *ArbitraryValue) Equals(rhs *ArbitraryValue) bool {
	{
		if !(_bool_EqualsPtr(lhs.BoolValue, rhs.BoolValue)) {
			return false
		}
	}
	{
		if !(_i64_EqualsPtr(lhs.Int64Value, rhs.Int64Value)) {
			return false
		}
	}
	{
		if !(_string_EqualsPtr(lhs.StringValue, rhs.StringValue)) {
			return false
		}
	}
	if (lhs.ListValue == nil && rhs.ListValue != nil) || (lhs.ListValue != nil && rhs.ListValue == nil) {
		return false
	} else if lhs.ListValue != nil && rhs.ListValue != nil {
		if !(_List_ArbitraryValue_Equals(lhs.ListValue, rhs.ListValue)) {
			return false
		}
	}
	if (lhs.MapValue == nil && rhs.MapValue != nil) || (lhs.MapValue != nil && rhs.MapValue == nil) {
		return false
	} else if lhs.MapValue != nil && rhs.MapValue != nil {
		if !(_Map_String_ArbitraryValue_Equals(lhs.MapValue, rhs.MapValue)) {
			return false
		}
	}
	return true
}

type Document struct {
	Pdf       typedefs.PDF `json:"pdf"`
	PlainText *string      `json:"plainText,omitempty"`
}

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

func (v *Document) String() string {
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

func (lhs *Document) Equals(rhs *Document) bool {
	if (lhs.Pdf == nil && rhs.Pdf != nil) || (lhs.Pdf != nil && rhs.Pdf == nil) {
		return false
	} else if lhs.Pdf != nil && rhs.Pdf != nil {
		if !(lhs.Pdf.Equals(rhs.Pdf)) {
			return false
		}
	}
	{
		if !(_string_EqualsPtr(lhs.PlainText, rhs.PlainText)) {
			return false
		}
	}
	return true
}

type EmptyUnion struct{}

func (v *EmptyUnion) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *EmptyUnion) FromWire(w wire.Value) error {
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}
	return nil
}

func (v *EmptyUnion) String() string {
	var fields [0]string
	i := 0
	return fmt.Sprintf("EmptyUnion{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *EmptyUnion) Equals(rhs *EmptyUnion) bool {
	return true
}
