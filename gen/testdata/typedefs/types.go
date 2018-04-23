// Code generated by thriftrw v1.12.0. DO NOT EDIT.
// @generated

package typedefs

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/thriftrw/gen/testdata/enums"
	"go.uber.org/thriftrw/gen/testdata/structs"
	"go.uber.org/thriftrw/wire"
	"strings"
)

type _Set_Binary_ValueList [][]byte

func (v _Set_Binary_ValueList) ForEach(f func(wire.Value) error) error {
	for _, x := range v {
		if x == nil {
			return fmt.Errorf("invalid set item: value is nil")
		}
		w, err := wire.NewValueBinary(x), error(nil)
		if err != nil {
			return err
		}

		if err := f(w); err != nil {
			return err
		}
	}
	return nil
}

func (v _Set_Binary_ValueList) Size() int {
	return len(v)
}

func (_Set_Binary_ValueList) ValueType() wire.Type {
	return wire.TBinary
}

func (_Set_Binary_ValueList) Close() {}

func _Set_Binary_Read(s wire.ValueList) ([][]byte, error) {
	if s.ValueType() != wire.TBinary {
		return nil, nil
	}

	o := make([][]byte, 0, s.Size())
	err := s.ForEach(func(x wire.Value) error {
		i, err := x.GetBinary(), error(nil)
		if err != nil {
			return err
		}

		o = append(o, i)
		return nil
	})
	s.Close()
	return o, err
}

func _Set_Binary_Equals(lhs, rhs [][]byte) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for _, x := range lhs {
		ok := false
		for _, y := range rhs {
			if bytes.Equal(x, y) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}

	return true
}

type BinarySet [][]byte

// ToWire translates BinarySet into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v BinarySet) ToWire() (wire.Value, error) {
	x := ([][]byte)(v)
	return wire.NewValueSet(_Set_Binary_ValueList(x)), error(nil)
}

// String returns a readable string representation of BinarySet.
func (v BinarySet) String() string {
	x := ([][]byte)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes BinarySet from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *BinarySet) FromWire(w wire.Value) error {
	x, err := _Set_Binary_Read(w.GetSet())
	*v = (BinarySet)(x)
	return err
}

// Equals returns true if this BinarySet is equal to the provided
// BinarySet.
func (lhs BinarySet) Equals(rhs BinarySet) bool {
	return _Set_Binary_Equals(lhs, rhs)
}

type DefaultPrimitiveTypedef struct {
	State *State `json:"state,omitempty"`
}

func _State_ptr(v State) *State {
	return &v
}

// ToWire translates a DefaultPrimitiveTypedef struct into a Thrift-level intermediate
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
func (v *DefaultPrimitiveTypedef) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.State == nil {
		v.State = _State_ptr("hello")
	}
	{
		w, err = v.State.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _State_Read(w wire.Value) (State, error) {
	var x State
	err := x.FromWire(w)
	return x, err
}

// FromWire deserializes a DefaultPrimitiveTypedef struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a DefaultPrimitiveTypedef struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v DefaultPrimitiveTypedef
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *DefaultPrimitiveTypedef) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				var x State
				x, err = _State_Read(field.Value)
				v.State = &x
				if err != nil {
					return err
				}

			}
		}
	}

	if v.State == nil {
		v.State = _State_ptr("hello")
	}

	return nil
}

// String returns a readable string representation of a DefaultPrimitiveTypedef
// struct.
func (v *DefaultPrimitiveTypedef) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [1]string
	i := 0
	if v.State != nil {
		fields[i] = fmt.Sprintf("State: %v", *(v.State))
		i++
	}

	return fmt.Sprintf("DefaultPrimitiveTypedef{%v}", strings.Join(fields[:i], ", "))
}

func _State_EqualsPtr(lhs, rhs *State) bool {
	if lhs != nil && rhs != nil {

		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

// Equals returns true if all the fields of this DefaultPrimitiveTypedef match the
// provided DefaultPrimitiveTypedef.
//
// This function performs a deep comparison.
func (v *DefaultPrimitiveTypedef) Equals(rhs *DefaultPrimitiveTypedef) bool {
	if !_State_EqualsPtr(v.State, rhs.State) {
		return false
	}

	return true
}

// GetState returns the value of State if it is set or its
// zero value if it is unset.
func (v *DefaultPrimitiveTypedef) GetState() (o State) {
	if v.State != nil {
		return *v.State
	}
	o = "hello"
	return
}

type _Map_Edge_Edge_MapItemList []struct {
	Key   *structs.Edge
	Value *structs.Edge
}

func (m _Map_Edge_Edge_MapItemList) ForEach(f func(wire.MapItem) error) error {
	for _, i := range m {
		k := i.Key
		v := i.Value
		if k == nil {
			return fmt.Errorf("invalid map key: value is nil")
		}
		if v == nil {
			return fmt.Errorf("invalid [%v]: value is nil", k)
		}
		kw, err := k.ToWire()
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

func (m _Map_Edge_Edge_MapItemList) Size() int {
	return len(m)
}

func (_Map_Edge_Edge_MapItemList) KeyType() wire.Type {
	return wire.TStruct
}

func (_Map_Edge_Edge_MapItemList) ValueType() wire.Type {
	return wire.TStruct
}

func (_Map_Edge_Edge_MapItemList) Close() {}

func _Edge_Read(w wire.Value) (*structs.Edge, error) {
	var v structs.Edge
	err := v.FromWire(w)
	return &v, err
}

func _Map_Edge_Edge_Read(m wire.MapItemList) ([]struct {
	Key   *structs.Edge
	Value *structs.Edge
}, error) {
	if m.KeyType() != wire.TStruct {
		return nil, nil
	}

	if m.ValueType() != wire.TStruct {
		return nil, nil
	}

	o := make([]struct {
		Key   *structs.Edge
		Value *structs.Edge
	}, 0, m.Size())
	err := m.ForEach(func(x wire.MapItem) error {
		k, err := _Edge_Read(x.Key)
		if err != nil {
			return err
		}

		v, err := _Edge_Read(x.Value)
		if err != nil {
			return err
		}

		o = append(o, struct {
			Key   *structs.Edge
			Value *structs.Edge
		}{k, v})
		return nil
	})
	m.Close()
	return o, err
}

func _Map_Edge_Edge_Equals(lhs, rhs []struct {
	Key   *structs.Edge
	Value *structs.Edge
}) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for _, i := range lhs {
		lk := i.Key
		lv := i.Value
		ok := false
		for _, j := range rhs {
			rk := j.Key
			rv := j.Value
			if !lk.Equals(rk) {
				continue
			}

			if !lv.Equals(rv) {
				return false
			}
			ok = true
			break
		}

		if !ok {
			return false
		}
	}
	return true
}

type EdgeMap []struct {
	Key   *structs.Edge
	Value *structs.Edge
}

// ToWire translates EdgeMap into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v EdgeMap) ToWire() (wire.Value, error) {
	x := ([]struct {
		Key   *structs.Edge
		Value *structs.Edge
	})(v)
	return wire.NewValueMap(_Map_Edge_Edge_MapItemList(x)), error(nil)
}

// String returns a readable string representation of EdgeMap.
func (v EdgeMap) String() string {
	x := ([]struct {
		Key   *structs.Edge
		Value *structs.Edge
	})(v)
	return fmt.Sprint(x)
}

// FromWire deserializes EdgeMap from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *EdgeMap) FromWire(w wire.Value) error {
	x, err := _Map_Edge_Edge_Read(w.GetMap())
	*v = (EdgeMap)(x)
	return err
}

// Equals returns true if this EdgeMap is equal to the provided
// EdgeMap.
func (lhs EdgeMap) Equals(rhs EdgeMap) bool {
	return _Map_Edge_Edge_Equals(lhs, rhs)
}

type Event struct {
	UUID *UUID      `json:"uuid,required"`
	Time *Timestamp `json:"time,omitempty"`
}

// ToWire translates a Event struct into a Thrift-level intermediate
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
func (v *Event) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.UUID == nil {
		return w, errors.New("field UUID of Event is required")
	}
	w, err = v.UUID.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	if v.Time != nil {
		w, err = v.Time.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _UUID_Read(w wire.Value) (*UUID, error) {
	var x UUID
	err := x.FromWire(w)
	return &x, err
}

func _Timestamp_Read(w wire.Value) (Timestamp, error) {
	var x Timestamp
	err := x.FromWire(w)
	return x, err
}

// FromWire deserializes a Event struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a Event struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v Event
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *Event) FromWire(w wire.Value) error {
	var err error

	uuidIsSet := false

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.UUID, err = _UUID_Read(field.Value)
				if err != nil {
					return err
				}
				uuidIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TI64 {
				var x Timestamp
				x, err = _Timestamp_Read(field.Value)
				v.Time = &x
				if err != nil {
					return err
				}

			}
		}
	}

	if !uuidIsSet {
		return errors.New("field UUID of Event is required")
	}

	return nil
}

// String returns a readable string representation of a Event
// struct.
func (v *Event) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("UUID: %v", v.UUID)
	i++
	if v.Time != nil {
		fields[i] = fmt.Sprintf("Time: %v", *(v.Time))
		i++
	}

	return fmt.Sprintf("Event{%v}", strings.Join(fields[:i], ", "))
}

func _Timestamp_EqualsPtr(lhs, rhs *Timestamp) bool {
	if lhs != nil && rhs != nil {

		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

// Equals returns true if all the fields of this Event match the
// provided Event.
//
// This function performs a deep comparison.
func (v *Event) Equals(rhs *Event) bool {
	if !v.UUID.Equals(rhs.UUID) {
		return false
	}
	if !_Timestamp_EqualsPtr(v.Time, rhs.Time) {
		return false
	}

	return true
}

// GetUUID returns the value of UUID if it is set or its
// zero value if it is unset.
func (v *Event) GetUUID() (o *UUID) { return v.UUID }

// GetTime returns the value of Time if it is set or its
// zero value if it is unset.
func (v *Event) GetTime() (o Timestamp) {
	if v.Time != nil {
		return *v.Time
	}

	return
}

type _List_Event_ValueList []*Event

func (v _List_Event_ValueList) ForEach(f func(wire.Value) error) error {
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

func (v _List_Event_ValueList) Size() int {
	return len(v)
}

func (_List_Event_ValueList) ValueType() wire.Type {
	return wire.TStruct
}

func (_List_Event_ValueList) Close() {}

func _Event_Read(w wire.Value) (*Event, error) {
	var v Event
	err := v.FromWire(w)
	return &v, err
}

func _List_Event_Read(l wire.ValueList) ([]*Event, error) {
	if l.ValueType() != wire.TStruct {
		return nil, nil
	}

	o := make([]*Event, 0, l.Size())
	err := l.ForEach(func(x wire.Value) error {
		i, err := _Event_Read(x)
		if err != nil {
			return err
		}
		o = append(o, i)
		return nil
	})
	l.Close()
	return o, err
}

func _List_Event_Equals(lhs, rhs []*Event) bool {
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

type EventGroup []*Event

// ToWire translates EventGroup into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v EventGroup) ToWire() (wire.Value, error) {
	x := ([]*Event)(v)
	return wire.NewValueList(_List_Event_ValueList(x)), error(nil)
}

// String returns a readable string representation of EventGroup.
func (v EventGroup) String() string {
	x := ([]*Event)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes EventGroup from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *EventGroup) FromWire(w wire.Value) error {
	x, err := _List_Event_Read(w.GetList())
	*v = (EventGroup)(x)
	return err
}

// Equals returns true if this EventGroup is equal to the provided
// EventGroup.
func (lhs EventGroup) Equals(rhs EventGroup) bool {
	return _List_Event_Equals(lhs, rhs)
}

type _Set_Frame_ValueList []*structs.Frame

func (v _Set_Frame_ValueList) ForEach(f func(wire.Value) error) error {
	for _, x := range v {
		if x == nil {
			return fmt.Errorf("invalid set item: value is nil")
		}
		w, err := x.ToWire()
		if err != nil {
			return err
		}

		if err := f(w); err != nil {
			return err
		}
	}
	return nil
}

func (v _Set_Frame_ValueList) Size() int {
	return len(v)
}

func (_Set_Frame_ValueList) ValueType() wire.Type {
	return wire.TStruct
}

func (_Set_Frame_ValueList) Close() {}

func _Frame_Read(w wire.Value) (*structs.Frame, error) {
	var v structs.Frame
	err := v.FromWire(w)
	return &v, err
}

func _Set_Frame_Read(s wire.ValueList) ([]*structs.Frame, error) {
	if s.ValueType() != wire.TStruct {
		return nil, nil
	}

	o := make([]*structs.Frame, 0, s.Size())
	err := s.ForEach(func(x wire.Value) error {
		i, err := _Frame_Read(x)
		if err != nil {
			return err
		}

		o = append(o, i)
		return nil
	})
	s.Close()
	return o, err
}

func _Set_Frame_Equals(lhs, rhs []*structs.Frame) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for _, x := range lhs {
		ok := false
		for _, y := range rhs {
			if x.Equals(y) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}

	return true
}

type FrameGroup []*structs.Frame

// ToWire translates FrameGroup into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v FrameGroup) ToWire() (wire.Value, error) {
	x := ([]*structs.Frame)(v)
	return wire.NewValueSet(_Set_Frame_ValueList(x)), error(nil)
}

// String returns a readable string representation of FrameGroup.
func (v FrameGroup) String() string {
	x := ([]*structs.Frame)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes FrameGroup from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *FrameGroup) FromWire(w wire.Value) error {
	x, err := _Set_Frame_Read(w.GetSet())
	*v = (FrameGroup)(x)
	return err
}

// Equals returns true if this FrameGroup is equal to the provided
// FrameGroup.
func (lhs FrameGroup) Equals(rhs FrameGroup) bool {
	return _Set_Frame_Equals(lhs, rhs)
}

func _EnumWithValues_Read(w wire.Value) (enums.EnumWithValues, error) {
	var v enums.EnumWithValues
	err := v.FromWire(w)
	return v, err
}

type MyEnum enums.EnumWithValues

// ToWire translates MyEnum into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v MyEnum) ToWire() (wire.Value, error) {
	x := (enums.EnumWithValues)(v)
	return x.ToWire()
}

// String returns a readable string representation of MyEnum.
func (v MyEnum) String() string {
	x := (enums.EnumWithValues)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes MyEnum from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *MyEnum) FromWire(w wire.Value) error {
	x, err := _EnumWithValues_Read(w)
	*v = (MyEnum)(x)
	return err
}

// Equals returns true if this MyEnum is equal to the provided
// MyEnum.
func (lhs MyEnum) Equals(rhs MyEnum) bool {
	return lhs.Equals(rhs)
}

type PDF []byte

// ToWire translates PDF into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v PDF) ToWire() (wire.Value, error) {
	x := ([]byte)(v)
	return wire.NewValueBinary(x), error(nil)
}

// String returns a readable string representation of PDF.
func (v PDF) String() string {
	x := ([]byte)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes PDF from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *PDF) FromWire(w wire.Value) error {
	x, err := w.GetBinary(), error(nil)
	*v = (PDF)(x)
	return err
}

// Equals returns true if this PDF is equal to the provided
// PDF.
func (lhs PDF) Equals(rhs PDF) bool {
	return bytes.Equal(lhs, rhs)
}

type _Map_Point_Point_MapItemList []struct {
	Key   *structs.Point
	Value *structs.Point
}

func (m _Map_Point_Point_MapItemList) ForEach(f func(wire.MapItem) error) error {
	for _, i := range m {
		k := i.Key
		v := i.Value
		if k == nil {
			return fmt.Errorf("invalid map key: value is nil")
		}
		if v == nil {
			return fmt.Errorf("invalid [%v]: value is nil", k)
		}
		kw, err := k.ToWire()
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

func (m _Map_Point_Point_MapItemList) Size() int {
	return len(m)
}

func (_Map_Point_Point_MapItemList) KeyType() wire.Type {
	return wire.TStruct
}

func (_Map_Point_Point_MapItemList) ValueType() wire.Type {
	return wire.TStruct
}

func (_Map_Point_Point_MapItemList) Close() {}

func _Point_Read(w wire.Value) (*structs.Point, error) {
	var v structs.Point
	err := v.FromWire(w)
	return &v, err
}

func _Map_Point_Point_Read(m wire.MapItemList) ([]struct {
	Key   *structs.Point
	Value *structs.Point
}, error) {
	if m.KeyType() != wire.TStruct {
		return nil, nil
	}

	if m.ValueType() != wire.TStruct {
		return nil, nil
	}

	o := make([]struct {
		Key   *structs.Point
		Value *structs.Point
	}, 0, m.Size())
	err := m.ForEach(func(x wire.MapItem) error {
		k, err := _Point_Read(x.Key)
		if err != nil {
			return err
		}

		v, err := _Point_Read(x.Value)
		if err != nil {
			return err
		}

		o = append(o, struct {
			Key   *structs.Point
			Value *structs.Point
		}{k, v})
		return nil
	})
	m.Close()
	return o, err
}

func _Map_Point_Point_Equals(lhs, rhs []struct {
	Key   *structs.Point
	Value *structs.Point
}) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for _, i := range lhs {
		lk := i.Key
		lv := i.Value
		ok := false
		for _, j := range rhs {
			rk := j.Key
			rv := j.Value
			if !lk.Equals(rk) {
				continue
			}

			if !lv.Equals(rv) {
				return false
			}
			ok = true
			break
		}

		if !ok {
			return false
		}
	}
	return true
}

type PointMap []struct {
	Key   *structs.Point
	Value *structs.Point
}

// ToWire translates PointMap into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v PointMap) ToWire() (wire.Value, error) {
	x := ([]struct {
		Key   *structs.Point
		Value *structs.Point
	})(v)
	return wire.NewValueMap(_Map_Point_Point_MapItemList(x)), error(nil)
}

// String returns a readable string representation of PointMap.
func (v PointMap) String() string {
	x := ([]struct {
		Key   *structs.Point
		Value *structs.Point
	})(v)
	return fmt.Sprint(x)
}

// FromWire deserializes PointMap from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *PointMap) FromWire(w wire.Value) error {
	x, err := _Map_Point_Point_Read(w.GetMap())
	*v = (PointMap)(x)
	return err
}

// Equals returns true if this PointMap is equal to the provided
// PointMap.
func (lhs PointMap) Equals(rhs PointMap) bool {
	return _Map_Point_Point_Equals(lhs, rhs)
}

type State string

// ToWire translates State into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v State) ToWire() (wire.Value, error) {
	x := (string)(v)
	return wire.NewValueString(x), error(nil)
}

// String returns a readable string representation of State.
func (v State) String() string {
	x := (string)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes State from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *State) FromWire(w wire.Value) error {
	x, err := w.GetString(), error(nil)
	*v = (State)(x)
	return err
}

// Equals returns true if this State is equal to the provided
// State.
func (lhs State) Equals(rhs State) bool {
	return (lhs == rhs)
}

// Number of seconds since epoch.
//
// Deprecated: Use ISOTime instead.
type Timestamp int64

// ToWire translates Timestamp into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v Timestamp) ToWire() (wire.Value, error) {
	x := (int64)(v)
	return wire.NewValueI64(x), error(nil)
}

// String returns a readable string representation of Timestamp.
func (v Timestamp) String() string {
	x := (int64)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes Timestamp from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *Timestamp) FromWire(w wire.Value) error {
	x, err := w.GetI64(), error(nil)
	*v = (Timestamp)(x)
	return err
}

// Equals returns true if this Timestamp is equal to the provided
// Timestamp.
func (lhs Timestamp) Equals(rhs Timestamp) bool {
	return (lhs == rhs)
}

type Transition struct {
	FromState State      `json:"fromState,required"`
	ToState   State      `json:"toState,required"`
	Events    EventGroup `json:"events,omitempty"`
}

// ToWire translates a Transition struct into a Thrift-level intermediate
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
func (v *Transition) ToWire() (wire.Value, error) {
	var (
		fields [3]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	w, err = v.FromState.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++

	w, err = v.ToState.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++
	if v.Events != nil {
		w, err = v.Events.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 3, Value: w}
		i++
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _EventGroup_Read(w wire.Value) (EventGroup, error) {
	var x EventGroup
	err := x.FromWire(w)
	return x, err
}

// FromWire deserializes a Transition struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a Transition struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v Transition
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *Transition) FromWire(w wire.Value) error {
	var err error

	fromStateIsSet := false
	toStateIsSet := false

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.FromState, err = _State_Read(field.Value)
				if err != nil {
					return err
				}
				fromStateIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TBinary {
				v.ToState, err = _State_Read(field.Value)
				if err != nil {
					return err
				}
				toStateIsSet = true
			}
		case 3:
			if field.Value.Type() == wire.TList {
				v.Events, err = _EventGroup_Read(field.Value)
				if err != nil {
					return err
				}

			}
		}
	}

	if !fromStateIsSet {
		return errors.New("field FromState of Transition is required")
	}

	if !toStateIsSet {
		return errors.New("field ToState of Transition is required")
	}

	return nil
}

// String returns a readable string representation of a Transition
// struct.
func (v *Transition) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [3]string
	i := 0
	fields[i] = fmt.Sprintf("FromState: %v", v.FromState)
	i++
	fields[i] = fmt.Sprintf("ToState: %v", v.ToState)
	i++
	if v.Events != nil {
		fields[i] = fmt.Sprintf("Events: %v", v.Events)
		i++
	}

	return fmt.Sprintf("Transition{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this Transition match the
// provided Transition.
//
// This function performs a deep comparison.
func (v *Transition) Equals(rhs *Transition) bool {
	if !(v.FromState == rhs.FromState) {
		return false
	}
	if !(v.ToState == rhs.ToState) {
		return false
	}
	if !((v.Events == nil && rhs.Events == nil) || (v.Events != nil && rhs.Events != nil && v.Events.Equals(rhs.Events))) {
		return false
	}

	return true
}

// GetFromState returns the value of FromState if it is set or its
// zero value if it is unset.
func (v *Transition) GetFromState() (o State) { return v.FromState }

// GetToState returns the value of ToState if it is set or its
// zero value if it is unset.
func (v *Transition) GetToState() (o State) { return v.ToState }

// GetEvents returns the value of Events if it is set or its
// zero value if it is unset.
func (v *Transition) GetEvents() (o EventGroup) {
	if v.Events != nil {
		return v.Events
	}

	return
}

type UUID I128

// ToWire translates UUID into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
func (v *UUID) ToWire() (wire.Value, error) {
	x := (*I128)(v)
	return x.ToWire()
}

// String returns a readable string representation of UUID.
func (v *UUID) String() string {
	x := (*I128)(v)
	return fmt.Sprint(x)
}

// FromWire deserializes UUID from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
func (v *UUID) FromWire(w wire.Value) error {
	return (*I128)(v).FromWire(w)
}

// Equals returns true if this UUID is equal to the provided
// UUID.
func (lhs *UUID) Equals(rhs *UUID) bool {
	return (*I128)(lhs).Equals((*I128)(rhs))
}

type I128 struct {
	High int64 `json:"high,required"`
	Low  int64 `json:"low,required"`
}

// ToWire translates a I128 struct into a Thrift-level intermediate
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
func (v *I128) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	w, err = wire.NewValueI64(v.High), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++

	w, err = wire.NewValueI64(v.Low), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a I128 struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a I128 struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v I128
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *I128) FromWire(w wire.Value) error {
	var err error

	highIsSet := false
	lowIsSet := false

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TI64 {
				v.High, err = field.Value.GetI64(), error(nil)
				if err != nil {
					return err
				}
				highIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TI64 {
				v.Low, err = field.Value.GetI64(), error(nil)
				if err != nil {
					return err
				}
				lowIsSet = true
			}
		}
	}

	if !highIsSet {
		return errors.New("field High of I128 is required")
	}

	if !lowIsSet {
		return errors.New("field Low of I128 is required")
	}

	return nil
}

// String returns a readable string representation of a I128
// struct.
func (v *I128) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("High: %v", v.High)
	i++
	fields[i] = fmt.Sprintf("Low: %v", v.Low)
	i++

	return fmt.Sprintf("I128{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this I128 match the
// provided I128.
//
// This function performs a deep comparison.
func (v *I128) Equals(rhs *I128) bool {
	if !(v.High == rhs.High) {
		return false
	}
	if !(v.Low == rhs.Low) {
		return false
	}

	return true
}

// GetHigh returns the value of High if it is set or its
// zero value if it is unset.
func (v *I128) GetHigh() (o int64) { return v.High }

// GetLow returns the value of Low if it is set or its
// zero value if it is unset.
func (v *I128) GetLow() (o int64) { return v.Low }
