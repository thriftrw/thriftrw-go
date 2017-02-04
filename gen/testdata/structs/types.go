// Code generated by thriftrw v1.1.0
// @generated

package structs

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/thriftrw/gen/testdata/enums"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/wire"
	"strings"
)

type ContactInfo struct {
	EmailAddress string `json:"emailAddress"`
}

func (v *ContactInfo) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	w, err = wire.NewValueString(v.EmailAddress), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *ContactInfo) FromWire(w wire.Value) error {
	var err error
	emailAddressIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.EmailAddress, err = field.Value.GetString(), error(nil)
				if err != nil {
					return err
				}
				emailAddressIsSet = true
			}
		}
	}
	if !emailAddressIsSet {
		return errors.New("field EmailAddress of ContactInfo is required")
	}
	return nil
}

func (v *ContactInfo) String() string {
	var fields [1]string
	i := 0
	fields[i] = fmt.Sprintf("EmailAddress: %v", v.EmailAddress)
	i++
	return fmt.Sprintf("ContactInfo{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *ContactInfo) Equals(rhs *ContactInfo) bool {
	if !(lhs.EmailAddress == rhs.EmailAddress) {
		return false
	}
	return true
}

type DefaultsStruct struct {
	RequiredPrimitive *int32             `json:"requiredPrimitive,omitempty"`
	OptionalPrimitive *int32             `json:"optionalPrimitive,omitempty"`
	RequiredEnum      *enums.EnumDefault `json:"requiredEnum,omitempty"`
	OptionalEnum      *enums.EnumDefault `json:"optionalEnum,omitempty"`
	RequiredList      []string           `json:"requiredList"`
	OptionalList      []float64          `json:"optionalList"`
	RequiredStruct    *Frame             `json:"requiredStruct,omitempty"`
	OptionalStruct    *Edge              `json:"optionalStruct,omitempty"`
}

func _EnumDefault_ptr(v enums.EnumDefault) *enums.EnumDefault {
	return &v
}

type _List_String_ValueList []string

func (v _List_String_ValueList) ForEach(f func(wire.Value) error) error {
	for _, x := range v {
		w, err := wire.NewValueString(x), error(nil)
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

func (v _List_String_ValueList) Size() int {
	return len(v)
}

func (_List_String_ValueList) ValueType() wire.Type {
	return wire.TBinary
}

func (_List_String_ValueList) Close() {
}

type _List_Double_ValueList []float64

func (v _List_Double_ValueList) ForEach(f func(wire.Value) error) error {
	for _, x := range v {
		w, err := wire.NewValueDouble(x), error(nil)
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

func (v _List_Double_ValueList) Size() int {
	return len(v)
}

func (_List_Double_ValueList) ValueType() wire.Type {
	return wire.TDouble
}

func (_List_Double_ValueList) Close() {
}

func (v *DefaultsStruct) ToWire() (wire.Value, error) {
	var (
		fields [8]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.RequiredPrimitive == nil {
		v.RequiredPrimitive = ptr.Int32(100)
	}
	{
		w, err = wire.NewValueI32(*(v.RequiredPrimitive)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	if v.OptionalPrimitive == nil {
		v.OptionalPrimitive = ptr.Int32(200)
	}
	{
		w, err = wire.NewValueI32(*(v.OptionalPrimitive)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}
	if v.RequiredEnum == nil {
		v.RequiredEnum = _EnumDefault_ptr(enums.EnumDefaultBar)
	}
	{
		w, err = v.RequiredEnum.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 3, Value: w}
		i++
	}
	if v.OptionalEnum == nil {
		v.OptionalEnum = _EnumDefault_ptr(enums.EnumDefaultBaz)
	}
	{
		w, err = v.OptionalEnum.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 4, Value: w}
		i++
	}
	if v.RequiredList == nil {
		v.RequiredList = []string{"hello", "world"}
	}
	{
		w, err = wire.NewValueList(_List_String_ValueList(v.RequiredList)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 5, Value: w}
		i++
	}
	if v.OptionalList == nil {
		v.OptionalList = []float64{1, 2, 3}
	}
	{
		w, err = wire.NewValueList(_List_Double_ValueList(v.OptionalList)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 6, Value: w}
		i++
	}
	if v.RequiredStruct == nil {
		v.RequiredStruct = &Frame{Size: &Size{Height: 200, Width: 100}, TopLeft: &Point{X: 1, Y: 2}}
	}
	{
		w, err = v.RequiredStruct.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 7, Value: w}
		i++
	}
	if v.OptionalStruct == nil {
		v.OptionalStruct = &Edge{EndPoint: &Point{X: 3, Y: 4}, StartPoint: &Point{X: 1, Y: 2}}
	}
	{
		w, err = v.OptionalStruct.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 8, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _EnumDefault_Read(w wire.Value) (enums.EnumDefault, error) {
	var v enums.EnumDefault
	err := v.FromWire(w)
	return v, err
}

func _List_String_Read(l wire.ValueList) ([]string, error) {
	if l.ValueType() != wire.TBinary {
		return nil, nil
	}
	o := make([]string, 0, l.Size())
	err := l.ForEach(func(x wire.Value) error {
		i, err := x.GetString(), error(nil)
		if err != nil {
			return err
		}
		o = append(o, i)
		return nil
	})
	l.Close()
	return o, err
}

func _List_Double_Read(l wire.ValueList) ([]float64, error) {
	if l.ValueType() != wire.TDouble {
		return nil, nil
	}
	o := make([]float64, 0, l.Size())
	err := l.ForEach(func(x wire.Value) error {
		i, err := x.GetDouble(), error(nil)
		if err != nil {
			return err
		}
		o = append(o, i)
		return nil
	})
	l.Close()
	return o, err
}

func _Frame_Read(w wire.Value) (*Frame, error) {
	var v Frame
	err := v.FromWire(w)
	return &v, err
}

func _Edge_Read(w wire.Value) (*Edge, error) {
	var v Edge
	err := v.FromWire(w)
	return &v, err
}

func (v *DefaultsStruct) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TI32 {
				var x int32
				x, err = field.Value.GetI32(), error(nil)
				v.RequiredPrimitive = &x
				if err != nil {
					return err
				}
			}
		case 2:
			if field.Value.Type() == wire.TI32 {
				var x int32
				x, err = field.Value.GetI32(), error(nil)
				v.OptionalPrimitive = &x
				if err != nil {
					return err
				}
			}
		case 3:
			if field.Value.Type() == wire.TI32 {
				var x enums.EnumDefault
				x, err = _EnumDefault_Read(field.Value)
				v.RequiredEnum = &x
				if err != nil {
					return err
				}
			}
		case 4:
			if field.Value.Type() == wire.TI32 {
				var x enums.EnumDefault
				x, err = _EnumDefault_Read(field.Value)
				v.OptionalEnum = &x
				if err != nil {
					return err
				}
			}
		case 5:
			if field.Value.Type() == wire.TList {
				v.RequiredList, err = _List_String_Read(field.Value.GetList())
				if err != nil {
					return err
				}
			}
		case 6:
			if field.Value.Type() == wire.TList {
				v.OptionalList, err = _List_Double_Read(field.Value.GetList())
				if err != nil {
					return err
				}
			}
		case 7:
			if field.Value.Type() == wire.TStruct {
				v.RequiredStruct, err = _Frame_Read(field.Value)
				if err != nil {
					return err
				}
			}
		case 8:
			if field.Value.Type() == wire.TStruct {
				v.OptionalStruct, err = _Edge_Read(field.Value)
				if err != nil {
					return err
				}
			}
		}
	}
	if v.RequiredPrimitive == nil {
		v.RequiredPrimitive = ptr.Int32(100)
	}
	if v.OptionalPrimitive == nil {
		v.OptionalPrimitive = ptr.Int32(200)
	}
	if v.RequiredEnum == nil {
		v.RequiredEnum = _EnumDefault_ptr(enums.EnumDefaultBar)
	}
	if v.OptionalEnum == nil {
		v.OptionalEnum = _EnumDefault_ptr(enums.EnumDefaultBaz)
	}
	if v.RequiredList == nil {
		v.RequiredList = []string{"hello", "world"}
	}
	if v.OptionalList == nil {
		v.OptionalList = []float64{1, 2, 3}
	}
	if v.RequiredStruct == nil {
		v.RequiredStruct = &Frame{Size: &Size{Height: 200, Width: 100}, TopLeft: &Point{X: 1, Y: 2}}
	}
	if v.OptionalStruct == nil {
		v.OptionalStruct = &Edge{EndPoint: &Point{X: 3, Y: 4}, StartPoint: &Point{X: 1, Y: 2}}
	}
	return nil
}

func (v *DefaultsStruct) String() string {
	var fields [8]string
	i := 0
	if v.RequiredPrimitive != nil {
		fields[i] = fmt.Sprintf("RequiredPrimitive: %v", *(v.RequiredPrimitive))
		i++
	}
	if v.OptionalPrimitive != nil {
		fields[i] = fmt.Sprintf("OptionalPrimitive: %v", *(v.OptionalPrimitive))
		i++
	}
	if v.RequiredEnum != nil {
		fields[i] = fmt.Sprintf("RequiredEnum: %v", *(v.RequiredEnum))
		i++
	}
	if v.OptionalEnum != nil {
		fields[i] = fmt.Sprintf("OptionalEnum: %v", *(v.OptionalEnum))
		i++
	}
	if v.RequiredList != nil {
		fields[i] = fmt.Sprintf("RequiredList: %v", v.RequiredList)
		i++
	}
	if v.OptionalList != nil {
		fields[i] = fmt.Sprintf("OptionalList: %v", v.OptionalList)
		i++
	}
	if v.RequiredStruct != nil {
		fields[i] = fmt.Sprintf("RequiredStruct: %v", v.RequiredStruct)
		i++
	}
	if v.OptionalStruct != nil {
		fields[i] = fmt.Sprintf("OptionalStruct: %v", v.OptionalStruct)
		i++
	}
	return fmt.Sprintf("DefaultsStruct{%v}", strings.Join(fields[:i], ", "))
}

func _i32_EqualsPtr(lhs, rhs *int32) bool {
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

func _EnumDefault_EqualsPtr(lhs, rhs *enums.EnumDefault) bool {
	if lhs != nil && rhs != nil {
		x := *lhs
		y := *rhs
		return x.Equals(y)
	} else if lhs == nil && rhs == nil {
		return true
	} else {
		return false
	}
}

func _List_String_Equals(lhs, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, lv := range lhs {
		rv := rhs[i]
		if !(lv == rv) {
			return false
		}
	}
	return true
}

func _List_Double_Equals(lhs, rhs []float64) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, lv := range lhs {
		rv := rhs[i]
		if !(lv == rv) {
			return false
		}
	}
	return true
}

func (lhs *DefaultsStruct) Equals(rhs *DefaultsStruct) bool {
	{
		if !(_i32_EqualsPtr(lhs.RequiredPrimitive, rhs.RequiredPrimitive)) {
			return false
		}
	}
	{
		if !(_i32_EqualsPtr(lhs.OptionalPrimitive, rhs.OptionalPrimitive)) {
			return false
		}
	}
	{
		if !(_EnumDefault_EqualsPtr(lhs.RequiredEnum, rhs.RequiredEnum)) {
			return false
		}
	}
	{
		if !(_EnumDefault_EqualsPtr(lhs.OptionalEnum, rhs.OptionalEnum)) {
			return false
		}
	}
	if (lhs.RequiredList == nil && rhs.RequiredList != nil) || (lhs.RequiredList != nil && rhs.RequiredList == nil) {
		return false
	} else if lhs.RequiredList != nil && rhs.RequiredList != nil {
		if !(_List_String_Equals(lhs.RequiredList, rhs.RequiredList)) {
			return false
		}
	}
	if (lhs.OptionalList == nil && rhs.OptionalList != nil) || (lhs.OptionalList != nil && rhs.OptionalList == nil) {
		return false
	} else if lhs.OptionalList != nil && rhs.OptionalList != nil {
		if !(_List_Double_Equals(lhs.OptionalList, rhs.OptionalList)) {
			return false
		}
	}
	if (lhs.RequiredStruct == nil && rhs.RequiredStruct != nil) || (lhs.RequiredStruct != nil && rhs.RequiredStruct == nil) {
		return false
	} else if lhs.RequiredStruct != nil && rhs.RequiredStruct != nil {
		if !(lhs.RequiredStruct.Equals(rhs.RequiredStruct)) {
			return false
		}
	}
	if (lhs.OptionalStruct == nil && rhs.OptionalStruct != nil) || (lhs.OptionalStruct != nil && rhs.OptionalStruct == nil) {
		return false
	} else if lhs.OptionalStruct != nil && rhs.OptionalStruct != nil {
		if !(lhs.OptionalStruct.Equals(rhs.OptionalStruct)) {
			return false
		}
	}
	return true
}

type Edge struct {
	StartPoint *Point `json:"startPoint"`
	EndPoint   *Point `json:"endPoint"`
}

func (v *Edge) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.StartPoint == nil {
		return w, errors.New("field StartPoint of Edge is required")
	}
	w, err = v.StartPoint.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	if v.EndPoint == nil {
		return w, errors.New("field EndPoint of Edge is required")
	}
	w, err = v.EndPoint.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _Point_Read(w wire.Value) (*Point, error) {
	var v Point
	err := v.FromWire(w)
	return &v, err
}

func (v *Edge) FromWire(w wire.Value) error {
	var err error
	startPointIsSet := false
	endPointIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.StartPoint, err = _Point_Read(field.Value)
				if err != nil {
					return err
				}
				startPointIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TStruct {
				v.EndPoint, err = _Point_Read(field.Value)
				if err != nil {
					return err
				}
				endPointIsSet = true
			}
		}
	}
	if !startPointIsSet {
		return errors.New("field StartPoint of Edge is required")
	}
	if !endPointIsSet {
		return errors.New("field EndPoint of Edge is required")
	}
	return nil
}

func (v *Edge) String() string {
	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("StartPoint: %v", v.StartPoint)
	i++
	fields[i] = fmt.Sprintf("EndPoint: %v", v.EndPoint)
	i++
	return fmt.Sprintf("Edge{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *Edge) Equals(rhs *Edge) bool {
	if !(lhs.StartPoint.Equals(rhs.StartPoint)) {
		return false
	}
	if !(lhs.EndPoint.Equals(rhs.EndPoint)) {
		return false
	}
	return true
}

type EmptyStruct struct{}

func (v *EmptyStruct) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *EmptyStruct) FromWire(w wire.Value) error {
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}
	return nil
}

func (v *EmptyStruct) String() string {
	var fields [0]string
	i := 0
	return fmt.Sprintf("EmptyStruct{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *EmptyStruct) Equals(rhs *EmptyStruct) bool {
	return true
}

type Frame struct {
	TopLeft *Point `json:"topLeft"`
	Size    *Size  `json:"size"`
}

func (v *Frame) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.TopLeft == nil {
		return w, errors.New("field TopLeft of Frame is required")
	}
	w, err = v.TopLeft.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	if v.Size == nil {
		return w, errors.New("field Size of Frame is required")
	}
	w, err = v.Size.ToWire()
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _Size_Read(w wire.Value) (*Size, error) {
	var v Size
	err := v.FromWire(w)
	return &v, err
}

func (v *Frame) FromWire(w wire.Value) error {
	var err error
	topLeftIsSet := false
	sizeIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.TopLeft, err = _Point_Read(field.Value)
				if err != nil {
					return err
				}
				topLeftIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TStruct {
				v.Size, err = _Size_Read(field.Value)
				if err != nil {
					return err
				}
				sizeIsSet = true
			}
		}
	}
	if !topLeftIsSet {
		return errors.New("field TopLeft of Frame is required")
	}
	if !sizeIsSet {
		return errors.New("field Size of Frame is required")
	}
	return nil
}

func (v *Frame) String() string {
	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("TopLeft: %v", v.TopLeft)
	i++
	fields[i] = fmt.Sprintf("Size: %v", v.Size)
	i++
	return fmt.Sprintf("Frame{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *Frame) Equals(rhs *Frame) bool {
	if !(lhs.TopLeft.Equals(rhs.TopLeft)) {
		return false
	}
	if !(lhs.Size.Equals(rhs.Size)) {
		return false
	}
	return true
}

type Graph struct {
	Edges []*Edge `json:"edges"`
}

type _List_Edge_ValueList []*Edge

func (v _List_Edge_ValueList) ForEach(f func(wire.Value) error) error {
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

func (v _List_Edge_ValueList) Size() int {
	return len(v)
}

func (_List_Edge_ValueList) ValueType() wire.Type {
	return wire.TStruct
}

func (_List_Edge_ValueList) Close() {
}

func (v *Graph) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.Edges == nil {
		return w, errors.New("field Edges of Graph is required")
	}
	w, err = wire.NewValueList(_List_Edge_ValueList(v.Edges)), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _List_Edge_Read(l wire.ValueList) ([]*Edge, error) {
	if l.ValueType() != wire.TStruct {
		return nil, nil
	}
	o := make([]*Edge, 0, l.Size())
	err := l.ForEach(func(x wire.Value) error {
		i, err := _Edge_Read(x)
		if err != nil {
			return err
		}
		o = append(o, i)
		return nil
	})
	l.Close()
	return o, err
}

func (v *Graph) FromWire(w wire.Value) error {
	var err error
	edgesIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TList {
				v.Edges, err = _List_Edge_Read(field.Value.GetList())
				if err != nil {
					return err
				}
				edgesIsSet = true
			}
		}
	}
	if !edgesIsSet {
		return errors.New("field Edges of Graph is required")
	}
	return nil
}

func (v *Graph) String() string {
	var fields [1]string
	i := 0
	fields[i] = fmt.Sprintf("Edges: %v", v.Edges)
	i++
	return fmt.Sprintf("Graph{%v}", strings.Join(fields[:i], ", "))
}

func _List_Edge_Equals(lhs, rhs []*Edge) bool {
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

func (lhs *Graph) Equals(rhs *Graph) bool {
	if !(_List_Edge_Equals(lhs.Edges, rhs.Edges)) {
		return false
	}
	return true
}

type List Node

func (v *List) ToWire() (wire.Value, error) {
	x := (*Node)(v)
	return x.ToWire()
}

func (v *List) String() string {
	x := (*Node)(v)
	return fmt.Sprint(x)
}

func (v *List) FromWire(w wire.Value) error {
	return (*Node)(v).FromWire(w)
}

func (lhs *List) Equals(rhs *List) bool {
	return (*Node)(lhs).Equals((*Node)(rhs))
}

type Node struct {
	Value int32 `json:"value"`
	Tail  *List `json:"tail,omitempty"`
}

func (v *Node) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	w, err = wire.NewValueI32(v.Value), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	if v.Tail != nil {
		w, err = v.Tail.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _List_Read(w wire.Value) (*List, error) {
	var x List
	err := x.FromWire(w)
	return &x, err
}

func (v *Node) FromWire(w wire.Value) error {
	var err error
	valueIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TI32 {
				v.Value, err = field.Value.GetI32(), error(nil)
				if err != nil {
					return err
				}
				valueIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TStruct {
				v.Tail, err = _List_Read(field.Value)
				if err != nil {
					return err
				}
			}
		}
	}
	if !valueIsSet {
		return errors.New("field Value of Node is required")
	}
	return nil
}

func (v *Node) String() string {
	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("Value: %v", v.Value)
	i++
	if v.Tail != nil {
		fields[i] = fmt.Sprintf("Tail: %v", v.Tail)
		i++
	}
	return fmt.Sprintf("Node{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *Node) Equals(rhs *Node) bool {
	if !(lhs.Value == rhs.Value) {
		return false
	}
	if (lhs.Tail == nil && rhs.Tail != nil) || (lhs.Tail != nil && rhs.Tail == nil) {
		return false
	} else if lhs.Tail != nil && rhs.Tail != nil {
		if !(lhs.Tail.Equals(rhs.Tail)) {
			return false
		}
	}
	return true
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v *Point) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	w, err = wire.NewValueDouble(v.X), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	w, err = wire.NewValueDouble(v.Y), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *Point) FromWire(w wire.Value) error {
	var err error
	xIsSet := false
	yIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TDouble {
				v.X, err = field.Value.GetDouble(), error(nil)
				if err != nil {
					return err
				}
				xIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TDouble {
				v.Y, err = field.Value.GetDouble(), error(nil)
				if err != nil {
					return err
				}
				yIsSet = true
			}
		}
	}
	if !xIsSet {
		return errors.New("field X of Point is required")
	}
	if !yIsSet {
		return errors.New("field Y of Point is required")
	}
	return nil
}

func (v *Point) String() string {
	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("X: %v", v.X)
	i++
	fields[i] = fmt.Sprintf("Y: %v", v.Y)
	i++
	return fmt.Sprintf("Point{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *Point) Equals(rhs *Point) bool {
	if !(lhs.X == rhs.X) {
		return false
	}
	if !(lhs.Y == rhs.Y) {
		return false
	}
	return true
}

type PrimitiveOptionalStruct struct {
	BoolField   *bool    `json:"boolField,omitempty"`
	ByteField   *int8    `json:"byteField,omitempty"`
	Int16Field  *int16   `json:"int16Field,omitempty"`
	Int32Field  *int32   `json:"int32Field,omitempty"`
	Int64Field  *int64   `json:"int64Field,omitempty"`
	DoubleField *float64 `json:"doubleField,omitempty"`
	StringField *string  `json:"stringField,omitempty"`
	BinaryField []byte   `json:"binaryField"`
}

func (v *PrimitiveOptionalStruct) ToWire() (wire.Value, error) {
	var (
		fields [8]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.BoolField != nil {
		w, err = wire.NewValueBool(*(v.BoolField)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	if v.ByteField != nil {
		w, err = wire.NewValueI8(*(v.ByteField)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}
	if v.Int16Field != nil {
		w, err = wire.NewValueI16(*(v.Int16Field)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 3, Value: w}
		i++
	}
	if v.Int32Field != nil {
		w, err = wire.NewValueI32(*(v.Int32Field)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 4, Value: w}
		i++
	}
	if v.Int64Field != nil {
		w, err = wire.NewValueI64(*(v.Int64Field)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 5, Value: w}
		i++
	}
	if v.DoubleField != nil {
		w, err = wire.NewValueDouble(*(v.DoubleField)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 6, Value: w}
		i++
	}
	if v.StringField != nil {
		w, err = wire.NewValueString(*(v.StringField)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 7, Value: w}
		i++
	}
	if v.BinaryField != nil {
		w, err = wire.NewValueBinary(v.BinaryField), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 8, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *PrimitiveOptionalStruct) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBool {
				var x bool
				x, err = field.Value.GetBool(), error(nil)
				v.BoolField = &x
				if err != nil {
					return err
				}
			}
		case 2:
			if field.Value.Type() == wire.TI8 {
				var x int8
				x, err = field.Value.GetI8(), error(nil)
				v.ByteField = &x
				if err != nil {
					return err
				}
			}
		case 3:
			if field.Value.Type() == wire.TI16 {
				var x int16
				x, err = field.Value.GetI16(), error(nil)
				v.Int16Field = &x
				if err != nil {
					return err
				}
			}
		case 4:
			if field.Value.Type() == wire.TI32 {
				var x int32
				x, err = field.Value.GetI32(), error(nil)
				v.Int32Field = &x
				if err != nil {
					return err
				}
			}
		case 5:
			if field.Value.Type() == wire.TI64 {
				var x int64
				x, err = field.Value.GetI64(), error(nil)
				v.Int64Field = &x
				if err != nil {
					return err
				}
			}
		case 6:
			if field.Value.Type() == wire.TDouble {
				var x float64
				x, err = field.Value.GetDouble(), error(nil)
				v.DoubleField = &x
				if err != nil {
					return err
				}
			}
		case 7:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.StringField = &x
				if err != nil {
					return err
				}
			}
		case 8:
			if field.Value.Type() == wire.TBinary {
				v.BinaryField, err = field.Value.GetBinary(), error(nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *PrimitiveOptionalStruct) String() string {
	var fields [8]string
	i := 0
	if v.BoolField != nil {
		fields[i] = fmt.Sprintf("BoolField: %v", *(v.BoolField))
		i++
	}
	if v.ByteField != nil {
		fields[i] = fmt.Sprintf("ByteField: %v", *(v.ByteField))
		i++
	}
	if v.Int16Field != nil {
		fields[i] = fmt.Sprintf("Int16Field: %v", *(v.Int16Field))
		i++
	}
	if v.Int32Field != nil {
		fields[i] = fmt.Sprintf("Int32Field: %v", *(v.Int32Field))
		i++
	}
	if v.Int64Field != nil {
		fields[i] = fmt.Sprintf("Int64Field: %v", *(v.Int64Field))
		i++
	}
	if v.DoubleField != nil {
		fields[i] = fmt.Sprintf("DoubleField: %v", *(v.DoubleField))
		i++
	}
	if v.StringField != nil {
		fields[i] = fmt.Sprintf("StringField: %v", *(v.StringField))
		i++
	}
	if v.BinaryField != nil {
		fields[i] = fmt.Sprintf("BinaryField: %v", v.BinaryField)
		i++
	}
	return fmt.Sprintf("PrimitiveOptionalStruct{%v}", strings.Join(fields[:i], ", "))
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

func _byte_EqualsPtr(lhs, rhs *int8) bool {
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

func _i16_EqualsPtr(lhs, rhs *int16) bool {
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

func _double_EqualsPtr(lhs, rhs *float64) bool {
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

func (lhs *PrimitiveOptionalStruct) Equals(rhs *PrimitiveOptionalStruct) bool {
	{
		if !(_bool_EqualsPtr(lhs.BoolField, rhs.BoolField)) {
			return false
		}
	}
	{
		if !(_byte_EqualsPtr(lhs.ByteField, rhs.ByteField)) {
			return false
		}
	}
	{
		if !(_i16_EqualsPtr(lhs.Int16Field, rhs.Int16Field)) {
			return false
		}
	}
	{
		if !(_i32_EqualsPtr(lhs.Int32Field, rhs.Int32Field)) {
			return false
		}
	}
	{
		if !(_i64_EqualsPtr(lhs.Int64Field, rhs.Int64Field)) {
			return false
		}
	}
	{
		if !(_double_EqualsPtr(lhs.DoubleField, rhs.DoubleField)) {
			return false
		}
	}
	{
		if !(_string_EqualsPtr(lhs.StringField, rhs.StringField)) {
			return false
		}
	}
	if (lhs.BinaryField == nil && rhs.BinaryField != nil) || (lhs.BinaryField != nil && rhs.BinaryField == nil) {
		return false
	} else if lhs.BinaryField != nil && rhs.BinaryField != nil {
		if !(bytes.Equal(lhs.BinaryField, rhs.BinaryField)) {
			return false
		}
	}
	return true
}

type PrimitiveRequiredStruct struct {
	BoolField   bool    `json:"boolField"`
	ByteField   int8    `json:"byteField"`
	Int16Field  int16   `json:"int16Field"`
	Int32Field  int32   `json:"int32Field"`
	Int64Field  int64   `json:"int64Field"`
	DoubleField float64 `json:"doubleField"`
	StringField string  `json:"stringField"`
	BinaryField []byte  `json:"binaryField"`
}

func (v *PrimitiveRequiredStruct) ToWire() (wire.Value, error) {
	var (
		fields [8]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	w, err = wire.NewValueBool(v.BoolField), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	w, err = wire.NewValueI8(v.ByteField), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++
	w, err = wire.NewValueI16(v.Int16Field), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 3, Value: w}
	i++
	w, err = wire.NewValueI32(v.Int32Field), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 4, Value: w}
	i++
	w, err = wire.NewValueI64(v.Int64Field), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 5, Value: w}
	i++
	w, err = wire.NewValueDouble(v.DoubleField), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 6, Value: w}
	i++
	w, err = wire.NewValueString(v.StringField), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 7, Value: w}
	i++
	if v.BinaryField == nil {
		return w, errors.New("field BinaryField of PrimitiveRequiredStruct is required")
	}
	w, err = wire.NewValueBinary(v.BinaryField), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 8, Value: w}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *PrimitiveRequiredStruct) FromWire(w wire.Value) error {
	var err error
	boolFieldIsSet := false
	byteFieldIsSet := false
	int16FieldIsSet := false
	int32FieldIsSet := false
	int64FieldIsSet := false
	doubleFieldIsSet := false
	stringFieldIsSet := false
	binaryFieldIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBool {
				v.BoolField, err = field.Value.GetBool(), error(nil)
				if err != nil {
					return err
				}
				boolFieldIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TI8 {
				v.ByteField, err = field.Value.GetI8(), error(nil)
				if err != nil {
					return err
				}
				byteFieldIsSet = true
			}
		case 3:
			if field.Value.Type() == wire.TI16 {
				v.Int16Field, err = field.Value.GetI16(), error(nil)
				if err != nil {
					return err
				}
				int16FieldIsSet = true
			}
		case 4:
			if field.Value.Type() == wire.TI32 {
				v.Int32Field, err = field.Value.GetI32(), error(nil)
				if err != nil {
					return err
				}
				int32FieldIsSet = true
			}
		case 5:
			if field.Value.Type() == wire.TI64 {
				v.Int64Field, err = field.Value.GetI64(), error(nil)
				if err != nil {
					return err
				}
				int64FieldIsSet = true
			}
		case 6:
			if field.Value.Type() == wire.TDouble {
				v.DoubleField, err = field.Value.GetDouble(), error(nil)
				if err != nil {
					return err
				}
				doubleFieldIsSet = true
			}
		case 7:
			if field.Value.Type() == wire.TBinary {
				v.StringField, err = field.Value.GetString(), error(nil)
				if err != nil {
					return err
				}
				stringFieldIsSet = true
			}
		case 8:
			if field.Value.Type() == wire.TBinary {
				v.BinaryField, err = field.Value.GetBinary(), error(nil)
				if err != nil {
					return err
				}
				binaryFieldIsSet = true
			}
		}
	}
	if !boolFieldIsSet {
		return errors.New("field BoolField of PrimitiveRequiredStruct is required")
	}
	if !byteFieldIsSet {
		return errors.New("field ByteField of PrimitiveRequiredStruct is required")
	}
	if !int16FieldIsSet {
		return errors.New("field Int16Field of PrimitiveRequiredStruct is required")
	}
	if !int32FieldIsSet {
		return errors.New("field Int32Field of PrimitiveRequiredStruct is required")
	}
	if !int64FieldIsSet {
		return errors.New("field Int64Field of PrimitiveRequiredStruct is required")
	}
	if !doubleFieldIsSet {
		return errors.New("field DoubleField of PrimitiveRequiredStruct is required")
	}
	if !stringFieldIsSet {
		return errors.New("field StringField of PrimitiveRequiredStruct is required")
	}
	if !binaryFieldIsSet {
		return errors.New("field BinaryField of PrimitiveRequiredStruct is required")
	}
	return nil
}

func (v *PrimitiveRequiredStruct) String() string {
	var fields [8]string
	i := 0
	fields[i] = fmt.Sprintf("BoolField: %v", v.BoolField)
	i++
	fields[i] = fmt.Sprintf("ByteField: %v", v.ByteField)
	i++
	fields[i] = fmt.Sprintf("Int16Field: %v", v.Int16Field)
	i++
	fields[i] = fmt.Sprintf("Int32Field: %v", v.Int32Field)
	i++
	fields[i] = fmt.Sprintf("Int64Field: %v", v.Int64Field)
	i++
	fields[i] = fmt.Sprintf("DoubleField: %v", v.DoubleField)
	i++
	fields[i] = fmt.Sprintf("StringField: %v", v.StringField)
	i++
	fields[i] = fmt.Sprintf("BinaryField: %v", v.BinaryField)
	i++
	return fmt.Sprintf("PrimitiveRequiredStruct{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *PrimitiveRequiredStruct) Equals(rhs *PrimitiveRequiredStruct) bool {
	if !(lhs.BoolField == rhs.BoolField) {
		return false
	}
	if !(lhs.ByteField == rhs.ByteField) {
		return false
	}
	if !(lhs.Int16Field == rhs.Int16Field) {
		return false
	}
	if !(lhs.Int32Field == rhs.Int32Field) {
		return false
	}
	if !(lhs.Int64Field == rhs.Int64Field) {
		return false
	}
	if !(lhs.DoubleField == rhs.DoubleField) {
		return false
	}
	if !(lhs.StringField == rhs.StringField) {
		return false
	}
	if !(bytes.Equal(lhs.BinaryField, rhs.BinaryField)) {
		return false
	}
	return true
}

type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (v *Size) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	w, err = wire.NewValueDouble(v.Width), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	w, err = wire.NewValueDouble(v.Height), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 2, Value: w}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *Size) FromWire(w wire.Value) error {
	var err error
	widthIsSet := false
	heightIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TDouble {
				v.Width, err = field.Value.GetDouble(), error(nil)
				if err != nil {
					return err
				}
				widthIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TDouble {
				v.Height, err = field.Value.GetDouble(), error(nil)
				if err != nil {
					return err
				}
				heightIsSet = true
			}
		}
	}
	if !widthIsSet {
		return errors.New("field Width of Size is required")
	}
	if !heightIsSet {
		return errors.New("field Height of Size is required")
	}
	return nil
}

func (v *Size) String() string {
	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("Width: %v", v.Width)
	i++
	fields[i] = fmt.Sprintf("Height: %v", v.Height)
	i++
	return fmt.Sprintf("Size{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *Size) Equals(rhs *Size) bool {
	if !(lhs.Width == rhs.Width) {
		return false
	}
	if !(lhs.Height == rhs.Height) {
		return false
	}
	return true
}

type User struct {
	Name    string       `json:"name"`
	Contact *ContactInfo `json:"contact,omitempty"`
}

func (v *User) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	w, err = wire.NewValueString(v.Name), error(nil)
	if err != nil {
		return w, err
	}
	fields[i] = wire.Field{ID: 1, Value: w}
	i++
	if v.Contact != nil {
		w, err = v.Contact.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _ContactInfo_Read(w wire.Value) (*ContactInfo, error) {
	var v ContactInfo
	err := v.FromWire(w)
	return &v, err
}

func (v *User) FromWire(w wire.Value) error {
	var err error
	nameIsSet := false
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.Name, err = field.Value.GetString(), error(nil)
				if err != nil {
					return err
				}
				nameIsSet = true
			}
		case 2:
			if field.Value.Type() == wire.TStruct {
				v.Contact, err = _ContactInfo_Read(field.Value)
				if err != nil {
					return err
				}
			}
		}
	}
	if !nameIsSet {
		return errors.New("field Name of User is required")
	}
	return nil
}

func (v *User) String() string {
	var fields [2]string
	i := 0
	fields[i] = fmt.Sprintf("Name: %v", v.Name)
	i++
	if v.Contact != nil {
		fields[i] = fmt.Sprintf("Contact: %v", v.Contact)
		i++
	}
	return fmt.Sprintf("User{%v}", strings.Join(fields[:i], ", "))
}

func (lhs *User) Equals(rhs *User) bool {
	if !(lhs.Name == rhs.Name) {
		return false
	}
	if (lhs.Contact == nil && rhs.Contact != nil) || (lhs.Contact != nil && rhs.Contact == nil) {
		return false
	} else if lhs.Contact != nil && rhs.Contact != nil {
		if !(lhs.Contact.Equals(rhs.Contact)) {
			return false
		}
	}
	return true
}
