// Code generated by thriftrw

package keyvalue

import (
	"errors"
	"fmt"
	"github.com/thriftrw/thriftrw-go/gen/testdata/exceptions"
	"github.com/thriftrw/thriftrw-go/gen/testdata/services"
	"github.com/thriftrw/thriftrw-go/wire"
	"strings"
)

type DeleteValueArgs struct {
	Key *services.Key `json:"key,omitempty"`
}

func (v *DeleteValueArgs) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.Key != nil {
		w, err = v.Key.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _Key_Read(w wire.Value) (services.Key, error) {
	var x services.Key
	err := x.FromWire(w)
	return x, err
}

func (v *DeleteValueArgs) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				var x services.Key
				x, err = _Key_Read(field.Value)
				v.Key = &x
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *DeleteValueArgs) String() string {
	var fields [1]string
	i := 0
	if v.Key != nil {
		fields[i] = fmt.Sprintf("Key: %v", *(v.Key))
		i++
	}
	return fmt.Sprintf("DeleteValueArgs{%v}", strings.Join(fields[:i], ", "))
}

func (v *DeleteValueArgs) MethodName() string {
	return "deleteValue"
}

func (v *DeleteValueArgs) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

type DeleteValueResult struct {
	DoesNotExist  *exceptions.DoesNotExistException `json:"doesNotExist,omitempty"`
	InternalError *services.InternalError           `json:"internalError,omitempty"`
}

func (v *DeleteValueResult) ToWire() (wire.Value, error) {
	var (
		fields [2]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.DoesNotExist != nil {
		w, err = v.DoesNotExist.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	if v.InternalError != nil {
		w, err = v.InternalError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}
	if i > 1 {
		return wire.Value{}, fmt.Errorf("DeleteValueResult should receive at most one field value: received %v values", i)
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _DoesNotExistException_Read(w wire.Value) (*exceptions.DoesNotExistException, error) {
	var v exceptions.DoesNotExistException
	err := v.FromWire(w)
	return &v, err
}

func _InternalError_Read(w wire.Value) (*services.InternalError, error) {
	var v services.InternalError
	err := v.FromWire(w)
	return &v, err
}

func (v *DeleteValueResult) FromWire(w wire.Value) error {
	var err error
	count := 0
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.DoesNotExist, err = _DoesNotExistException_Read(field.Value)
				if err != nil {
					return err
				}
				count++
			}
		case 2:
			if field.Value.Type() == wire.TStruct {
				v.InternalError, err = _InternalError_Read(field.Value)
				if err != nil {
					return err
				}
				count++
			}
		}
	}
	if count > 1 {
		return fmt.Errorf("DeleteValueResult should receive at most one field value: received %v values", count)
	}
	return nil
}

func (v *DeleteValueResult) String() string {
	var fields [2]string
	i := 0
	if v.DoesNotExist != nil {
		fields[i] = fmt.Sprintf("DoesNotExist: %v", v.DoesNotExist)
		i++
	}
	if v.InternalError != nil {
		fields[i] = fmt.Sprintf("InternalError: %v", v.InternalError)
		i++
	}
	return fmt.Sprintf("DeleteValueResult{%v}", strings.Join(fields[:i], ", "))
}

func (v *DeleteValueResult) MethodName() string {
	return "deleteValue"
}

func (v *DeleteValueResult) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}

var DeleteValueHelper = struct {
	IsException    func(error) bool
	Args           func(key *services.Key) *DeleteValueArgs
	WrapResponse   func(error) (*DeleteValueResult, error)
	UnwrapResponse func(*DeleteValueResult) error
}{}

func init() {
	DeleteValueHelper.IsException = func(err error) bool {
		switch err.(type) {
		case *exceptions.DoesNotExistException:
			return true
		case *services.InternalError:
			return true
		default:
			return false
		}
	}
	DeleteValueHelper.Args = func(key *services.Key) *DeleteValueArgs {
		return &DeleteValueArgs{Key: key}
	}
	DeleteValueHelper.WrapResponse = func(err error) (*DeleteValueResult, error) {
		if err == nil {
			return &DeleteValueResult{}, nil
		}
		switch e := err.(type) {
		case *exceptions.DoesNotExistException:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for DeleteValueResult.DoesNotExist")
			}
			return &DeleteValueResult{DoesNotExist: e}, nil
		case *services.InternalError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for DeleteValueResult.InternalError")
			}
			return &DeleteValueResult{InternalError: e}, nil
		}
		return nil, err
	}
	DeleteValueHelper.UnwrapResponse = func(result *DeleteValueResult) (err error) {
		if result.DoesNotExist != nil {
			err = result.DoesNotExist
			return
		}
		if result.InternalError != nil {
			err = result.InternalError
			return
		}
		return
	}
}
