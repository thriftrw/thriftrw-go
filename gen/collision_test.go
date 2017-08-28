// Copyright (c) 2015 Uber Technologies, Inc.
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

package gen

import (
	"bytes"
	"reflect"
	"testing"

	tc "go.uber.org/thriftrw/gen/testdata/collision"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStruct(t *testing.T) {
	tests := []struct {
		desc string
		x    thriftType
	}{
		{
			"StructCollision",
			&tc.StructCollision{
				CollisionField:  true,
				CollisionField2: "true",
			},
		},
		{
			"StructCollision2 (struct_collision)",
			&tc.StructCollision2{
				CollisionField:  true,
				CollisionField2: "such trueness",
			},
		},
		{
			"PrimitiveContainers",
			&tc.PrimitiveContainers{
				A: []string{"arbre", "fleur"},
				B: map[string]struct{}{
					"croissant": {},
					"baguette":  {},
				},
				C: map[string]string{
					"voiture": "bleue",
					"cammion": "rouge",
				},
			},
		},
		{
			"StructConstant (struct_constant)",
			tc.StructConstant,
		},
		{
			"UnionCollision .CollisionField",
			&tc.UnionCollision{
				CollisionField: ptr.Bool(true),
			},
		},
		{
			"UnionCollision .collision_field",
			&tc.UnionCollision{
				CollisionField2: ptr.String("so true"),
			},
		},
		{
			"UnionCollision2 (union_collision) .CollisionField",
			&tc.UnionCollision2{
				CollisionField: &[]bool{true}[0],
			},
		},
		{
			"UnionCollision (union_collision) .collision_field",
			&tc.UnionCollision2{
				CollisionField2: ptr.String("true indeed"),
			},
		},
	}
	for _, tt := range tests {
		roundTrip(t, tt.x, tt.desc)
	}
}

func TestStructEquals(t *testing.T) {
	tests := []struct {
		desc     string
		lhs, rhs interface{}
		want     bool
	}{
		{
			"StructCollision equal",
			&tc.StructCollision{
				CollisionField: true,
			},
			&tc.StructCollision{
				CollisionField: true,
			},
			true,
		},
		{
			"StructCollision unequal",
			&tc.StructCollision{
				CollisionField:  true,
				CollisionField2: "true",
			},
			&tc.StructCollision{
				CollisionField:  true,
				CollisionField2: "false",
			},
			false,
		},
		{
			"Empty StructCollision equal",
			&tc.StructCollision{},
			&tc.StructCollision{},
			true,
		},
		{
			"Empty StructCollision unequal",
			&tc.StructCollision{},
			&tc.StructCollision{
				CollisionField: true,
			},
			false,
		},
		{
			"StructCollision2 (struct_collision) equal",
			&tc.StructCollision2{
				CollisionField:  true,
				CollisionField2: "such trueness",
			},
			&tc.StructCollision2{
				CollisionField:  true,
				CollisionField2: "such trueness",
			},
			true,
		},
		{
			"StructCollision2 (struct_collision) unequal",
			&tc.StructCollision2{
				CollisionField:  false,
				CollisionField2: "such trueness",
			},
			&tc.StructCollision2{
				CollisionField:  true,
				CollisionField2: "such trueness",
			},
			false,
		},
		{
			"Empty StructCollision2 (struct_collision) equal",
			&tc.StructCollision2{},
			&tc.StructCollision2{},
			true,
		},
		{
			"Empty StructCollision2 (struct_collision) unequal",
			&tc.StructCollision2{},
			&tc.StructCollision2{
				CollisionField: true,
			},
			false,
		},
		{
			"PrimitiveContainers equal",
			&tc.PrimitiveContainers{
				A: []string{"arbre", "fleur"},
				B: map[string]struct{}{
					"croissant": {},
					"baguette":  {},
				},
				C: map[string]string{
					"voiture": "bleue",
					"cammion": "rouge",
				},
			},
			&tc.PrimitiveContainers{
				A: []string{"arbre", "fleur"},
				B: map[string]struct{}{
					"croissant": {},
					"baguette":  {},
				},
				C: map[string]string{
					"voiture": "bleue",
					"cammion": "rouge",
				},
			},
			true,
		},
		{
			"PrimitiveContainers unequal",
			&tc.PrimitiveContainers{
				A: []string{"arbre", "fleur"},
				B: map[string]struct{}{
					"croissant": {},
					"baguette":  {},
				},
				C: map[string]string{
					"voiture": "bleue",
					"cammion": "rouge",
				},
			},
			&tc.PrimitiveContainers{
				A: []string{"arbre", "fleur"},
				B: map[string]struct{}{
					"croissant": {}, // No baguette
				},
				C: map[string]string{
					"voiture": "bleue",
					"cammion": "rouge",
				},
			},
			false,
		},
		{
			"Empty PrimitiveContainers equal",
			&tc.PrimitiveContainers{},
			&tc.PrimitiveContainers{},
			true,
		},
		{
			"Empty PrimitiveContainers unequal",
			&tc.PrimitiveContainers{},
			&tc.PrimitiveContainers{
				A: []string{"arbre", "fleur"},
			},
			false,
		},
		{
			"StructConstant (struct_constant) equal",
			tc.StructConstant,
			&tc.StructCollision2{
				CollisionField2: "false indeed",
			},
			true,
		},
		{
			"StructConstant (struct_constant) unequal",
			tc.StructConstant,
			&tc.StructCollision2{
				CollisionField2: "true indeed",
			},
			false,
		},
		{
			"UnionCollision .CollisionField equal",
			&tc.UnionCollision{
				CollisionField: ptr.Bool(true),
			},
			&tc.UnionCollision{
				CollisionField: ptr.Bool(true),
			},
			true,
		},
		{
			"UnionCollision .collision_field2 equal",
			&tc.UnionCollision{
				CollisionField2: ptr.String("so true"),
			},
			&tc.UnionCollision{
				CollisionField2: ptr.String("so true"),
			},
			true,
		},
		{
			"UnionCollision .CollisionField unequal",
			&tc.UnionCollision{
				CollisionField: ptr.Bool(true),
			},
			&tc.UnionCollision{
				CollisionField: ptr.Bool(false),
			},
			false,
		},
		{
			"UnionCollision .collision_field2 unequal",
			&tc.UnionCollision{
				CollisionField2: ptr.String("so true"),
			},
			&tc.UnionCollision{
				CollisionField2: ptr.String("not so true"),
			},
			false,
		},
		{
			"UnionCollision equal",
			&tc.UnionCollision{
				CollisionField:  ptr.Bool(true),
				CollisionField2: ptr.String("so true"),
			},
			&tc.UnionCollision{
				CollisionField:  ptr.Bool(true),
				CollisionField2: ptr.String("so true"),
			},
			true,
		},
		{
			"UnionCollision unequal",
			&tc.UnionCollision{
				CollisionField:  ptr.Bool(true),
				CollisionField2: ptr.String("so true"),
			},
			&tc.UnionCollision{
				CollisionField:  ptr.Bool(true),
				CollisionField2: ptr.String("not so true"),
			},
			false,
		},
		{
			"Empty UnionCollision equal",
			&tc.UnionCollision{},
			&tc.UnionCollision{},
			true,
		},
		{
			"Empty UnionCollision unequal",
			&tc.UnionCollision{},
			&tc.UnionCollision{
				CollisionField: ptr.Bool(false),
			},
			false,
		},
		{
			"UnionCollision2 (union_collision) .CollisionField equal",
			&tc.UnionCollision2{
				CollisionField: &[]bool{true}[0],
			},
			&tc.UnionCollision2{
				CollisionField: &[]bool{true}[0],
			},
			true,
		},
		{
			"UnionCollision (union_collision) .collision_field2 equal",
			&tc.UnionCollision2{
				CollisionField2: ptr.String("true indeed"),
			},
			&tc.UnionCollision2{
				CollisionField2: ptr.String("true indeed"),
			},
			true,
		},
		{
			"UnionCollision2 (union_collision) .CollisionField unequal",
			&tc.UnionCollision2{
				CollisionField: &[]bool{false}[0],
			},
			&tc.UnionCollision2{
				CollisionField: &[]bool{true}[0],
			},
			false,
		},
		{
			"UnionCollision (union_collision) .collision_field2 unequal",
			&tc.UnionCollision2{
				CollisionField2: ptr.String("true indeed"),
			},
			&tc.UnionCollision2{
				CollisionField2: ptr.String("false indeed"),
			},
			false,
		},
		{
			"UnionCollision (union_collision) equal",
			&tc.UnionCollision2{
				CollisionField:  &[]bool{false}[0],
				CollisionField2: ptr.String("true indeed"),
			},
			&tc.UnionCollision2{
				CollisionField:  &[]bool{false}[0],
				CollisionField2: ptr.String("true indeed"),
			},
			true,
		},
		{
			"UnionCollision (union_collision) unequal",
			&tc.UnionCollision2{
				CollisionField:  &[]bool{true}[0],
				CollisionField2: ptr.String("true indeed"),
			},
			&tc.UnionCollision2{
				CollisionField:  &[]bool{false}[0],
				CollisionField2: ptr.String("true indeed"),
			},
			false,
		},
		{
			"Empty UnionCollision (union_collision) equal",
			&tc.UnionCollision2{},
			&tc.UnionCollision2{},
			true,
		},
		{
			"Empty UnionCollision (union_collision) unequal",
			&tc.UnionCollision{},
			&tc.UnionCollision{
				CollisionField: &[]bool{true}[0],
			},
			false,
		},
	}
	for _, tt := range tests {
		args := []reflect.Value{reflect.ValueOf(tt.rhs)}
		ret := reflect.ValueOf(tt.lhs).MethodByName("Equals").Call(args)[0].Interface()
		assert.Equal(t, ret, tt.want, tt.desc)
	}
}

func TestConstant(t *testing.T) {
	require.Equal(t, tc.StructCollision2{
		CollisionField2: "false indeed",
	}, *tc.StructConstant)
}

func TestWithDefault(t *testing.T) {
	a := &tc.WithDefault{
		Pouet: &tc.StructCollision2{
			CollisionField2: "false indeed",
		},
	}
	b := &tc.WithDefault{}

	a = roundTrip(t, a, "WithDefault{filled in}").(*tc.WithDefault)
	b = roundTrip(t, b, "WithDefault{}").(*tc.WithDefault)
	if a != nil && b != nil {
		require.Equal(t, a, b)
	}
}

func TestWithDefaultEquals(t *testing.T) {
	a := tc.WithDefault{
		Pouet: &tc.StructCollision2{
			CollisionField2: "false indeed",
		},
	}
	b := tc.WithDefault{}
	c := tc.WithDefault{
		Pouet: &tc.StructCollision2{
			CollisionField: true,
		},
	}

	tests := []struct {
		lhs, rhs tc.WithDefault
		want     bool
	}{
		{a, a, true},
		{b, b, true},
		{c, c, true},
		{a, b, false},
		{a, c, false},
		{b, c, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.lhs.Equals(&tt.rhs), tt.want)
	}
}

func TestMyEnum(t *testing.T) {
	tests := []struct {
		e tc.MyEnum
		n string
		v int64
	}{
		{tc.MyEnumX, "X", 123},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.n, tt.e.String())
		assert.Equal(t, tt.v, int64(tt.e))
	}
}

func TestMyEnumEquals(t *testing.T) {
	tests := []struct {
		lhs, rhs tc.MyEnum
		want     bool
	}{
		{tc.MyEnumX, tc.MyEnumX, true},
		{tc.MyEnumFooBar, tc.MyEnumFooBar2, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.lhs.Equals(tt.rhs), tt.want)
	}
}

func TestMyEnum2(t *testing.T) {
	tests := []struct {
		e tc.MyEnum2
		n string
		v int64
	}{
		{tc.MyEnum2X, "X", 12},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.n, tt.e.String())
		assert.Equal(t, tt.v, int64(tt.e))
	}
}

func TestMyEnum2Equals(t *testing.T) {
	tests := []struct {
		lhs, rhs tc.MyEnum2
		want     bool
	}{
		{tc.MyEnum2X, tc.MyEnum2X, true},
		{tc.MyEnum2Y, tc.MyEnum2Z, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.lhs.Equals(tt.rhs), tt.want)
	}
}

func roundTrip(t *testing.T, x thriftType, msg string) thriftType {
	v, err := x.ToWire()
	if !assert.NoError(t, err, "failed to serialize: %v", x) {
		return nil
	}

	var buff bytes.Buffer
	if !assert.NoError(t, protocol.Binary.Encode(v, &buff), "%v: failed to serialize", msg) {
		return nil
	}

	newV, err := protocol.Binary.Decode(bytes.NewReader(buff.Bytes()), v.Type())
	if !assert.NoError(t, err, "%v: failed to deserialize", msg) {
		return nil
	}

	if !assert.True(
		t, wire.ValuesAreEqual(newV, v), "%v: deserialize(serialize(%v.ToWire())) != %v", msg, x, v) {
		return nil
	}

	xType := reflect.TypeOf(x)
	if xType.Kind() == reflect.Ptr {
		xType = xType.Elem()
	}

	gotX := reflect.New(xType)
	errval := gotX.MethodByName("FromWire").
		Call([]reflect.Value{reflect.ValueOf(newV)})[0].
		Interface()

	if !assert.Nil(t, errval, "FromWire: %v", msg) {
		return nil
	}

	assert.Equal(t, x, gotX.Interface(), "FromWire: %v", msg)
	return gotX.Interface().(thriftType)
}

func TestCollisionAccessors(t *testing.T) {
	t.Run("UnionCollision", func(t *testing.T) {
		t.Run("bool", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				u := tc.UnionCollision{CollisionField: ptr.Bool(true)}
				assert.True(t, u.GetCollisionField())
			})
			t.Run("unset", func(t *testing.T) {
				var u tc.UnionCollision
				assert.False(t, u.GetCollisionField())
			})
		})

		t.Run("string", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				u := tc.UnionCollision{CollisionField2: ptr.String("foo")}
				assert.Equal(t, "foo", u.GetCollisionField2())
			})
			t.Run("unset", func(t *testing.T) {
				var u tc.UnionCollision
				assert.Equal(t, "", u.GetCollisionField2())
			})
		})
	})

	t.Run("UnionCollision2", func(t *testing.T) {
		t.Run("bool", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				u := tc.UnionCollision2{CollisionField: ptr.Bool(true)}
				assert.True(t, u.GetCollisionField())
			})
			t.Run("unset", func(t *testing.T) {
				var u tc.UnionCollision2
				assert.False(t, u.GetCollisionField())
			})
		})

		t.Run("string", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				u := tc.UnionCollision2{CollisionField2: ptr.String("foo")}
				assert.Equal(t, "foo", u.GetCollisionField2())
			})
			t.Run("unset", func(t *testing.T) {
				var u tc.UnionCollision2
				assert.Equal(t, "", u.GetCollisionField2())
			})
		})
	})

	t.Run("AccessorConflict", func(t *testing.T) {
		t.Run("name", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				u := tc.AccessorConflict{Name: ptr.String("foo")}
				assert.Equal(t, "foo", u.GetName())
			})
			t.Run("unset", func(t *testing.T) {
				u := tc.AccessorConflict{GetName2: ptr.String("bar")}
				assert.Equal(t, "", u.GetName())
			})
		})

		t.Run("get_name", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				u := tc.AccessorConflict{GetName2: ptr.String("foo")}
				assert.Equal(t, "foo", u.GetGetName2())
			})
			t.Run("unset", func(t *testing.T) {
				u := tc.AccessorConflict{Name: ptr.String("bar")}
				assert.Equal(t, "", u.GetGetName2())
			})
		})
	})
}
