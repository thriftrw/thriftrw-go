// Code generated by thriftrw v1.6.0. DO NOT EDIT.
// @generated

package services

import (
	"fmt"
	"go.uber.org/thriftrw/wire"
	"strings"
)

type Cache_ClearAfter_Args struct {
	DurationMS *int64 `json:"durationMS,omitempty"`
}

func (v *Cache_ClearAfter_Args) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.DurationMS != nil {
		w, err = wire.NewValueI64(*(v.DurationMS)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *Cache_ClearAfter_Args) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TI64 {
				var x int64
				x, err = field.Value.GetI64(), error(nil)
				v.DurationMS = &x
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *Cache_ClearAfter_Args) String() string {
	if v == nil {
		return "<nil>"
	}
	var fields [1]string
	i := 0
	if v.DurationMS != nil {
		fields[i] = fmt.Sprintf("DurationMS: %v", *(v.DurationMS))
		i++
	}
	return fmt.Sprintf("Cache_ClearAfter_Args{%v}", strings.Join(fields[:i], ", "))
}

func _I64_EqualsPtr(lhs, rhs *int64) bool {
	if lhs != nil && rhs != nil {
		x := *lhs
		y := *rhs
		return (x == y)
	}
	return lhs == nil && rhs == nil
}

func (v *Cache_ClearAfter_Args) Equals(rhs *Cache_ClearAfter_Args) bool {
	if !_I64_EqualsPtr(v.DurationMS, rhs.DurationMS) {
		return false
	}
	return true
}

func (v *Cache_ClearAfter_Args) MethodName() string {
	return "clearAfter"
}

func (v *Cache_ClearAfter_Args) EnvelopeType() wire.EnvelopeType {
	return wire.OneWay
}

var Cache_ClearAfter_Helper = struct {
	Args func(durationMS *int64) *Cache_ClearAfter_Args
}{}

func init() {
	Cache_ClearAfter_Helper.Args = func(durationMS *int64) *Cache_ClearAfter_Args {
		return &Cache_ClearAfter_Args{DurationMS: durationMS}
	}
}
