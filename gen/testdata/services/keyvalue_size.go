// Code generated by thriftrw v1.1.0
// @generated

package services

import (
	"errors"
	"fmt"
	"go.uber.org/thriftrw/wire"
	"strings"
)

type KeyValue_Size_Args struct{}

func (v *KeyValue_Size_Args) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *KeyValue_Size_Args) FromWire(w wire.Value) error {
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}
	return nil
}

func (v *KeyValue_Size_Args) String() string {
	if v == nil {
		return "<nil>"
	}
	var fields [0]string
	i := 0
	return fmt.Sprintf("KeyValue_Size_Args{%v}", strings.Join(fields[:i], ", "))
}

func (v *KeyValue_Size_Args) Equals(rhs *KeyValue_Size_Args) bool {
	return true
}

func (v *KeyValue_Size_Args) MethodName() string {
	return "size"
}

func (v *KeyValue_Size_Args) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

var KeyValue_Size_Helper = struct {
	Args           func() *KeyValue_Size_Args
	IsException    func(error) bool
	WrapResponse   func(int64, error) (*KeyValue_Size_Result, error)
	UnwrapResponse func(*KeyValue_Size_Result) (int64, error)
}{}

func init() {
	KeyValue_Size_Helper.Args = func() *KeyValue_Size_Args {
		return &KeyValue_Size_Args{}
	}
	KeyValue_Size_Helper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}
	KeyValue_Size_Helper.WrapResponse = func(success int64, err error) (*KeyValue_Size_Result, error) {
		if err == nil {
			return &KeyValue_Size_Result{Success: &success}, nil
		}
		return nil, err
	}
	KeyValue_Size_Helper.UnwrapResponse = func(result *KeyValue_Size_Result) (success int64, err error) {
		if result.Success != nil {
			success = *result.Success
			return
		}
		err = errors.New("expected a non-void result")
		return
	}
}

type KeyValue_Size_Result struct {
	Success *int64 `json:"success,omitempty"`
}

func (v *KeyValue_Size_Result) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.Success != nil {
		w, err = wire.NewValueI64(*(v.Success)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 0, Value: w}
		i++
	}
	if i != 1 {
		return wire.Value{}, fmt.Errorf("KeyValue_Size_Result should have exactly one field: got %v fields", i)
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *KeyValue_Size_Result) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 0:
			if field.Value.Type() == wire.TI64 {
				var x int64
				x, err = field.Value.GetI64(), error(nil)
				v.Success = &x
				if err != nil {
					return err
				}
			}
		}
	}
	count := 0
	if v.Success != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("KeyValue_Size_Result should have exactly one field: got %v fields", count)
	}
	return nil
}

func (v *KeyValue_Size_Result) String() string {
	if v == nil {
		return "<nil>"
	}
	var fields [1]string
	i := 0
	if v.Success != nil {
		fields[i] = fmt.Sprintf("Success: %v", *(v.Success))
		i++
	}
	return fmt.Sprintf("KeyValue_Size_Result{%v}", strings.Join(fields[:i], ", "))
}

func (v *KeyValue_Size_Result) Equals(rhs *KeyValue_Size_Result) bool {
	if !_I64_EqualsPtr(v.Success, rhs.Success) {
		return false
	}
	return true
}

func (v *KeyValue_Size_Result) MethodName() string {
	return "size"
}

func (v *KeyValue_Size_Result) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}
