package ast_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/thriftrw/ast"

	"github.com/golang/mock/gomock"
)

func TestWalk(t *testing.T) {
	type visit struct {
		// expected node and that node's ancestors for each visit
		node      ast.Node
		parent    ast.Node
		ancestors []ast.Node
	}

	type test struct {
		desc   string
		node   ast.Node
		visits []visit
	}

	tests := []test{
		{
			desc: "annotation",
			node: &ast.Annotation{Name: "foo", Value: "bar"},
			visits: []visit{
				{node: &ast.Annotation{Name: "foo", Value: "bar"}},
			},
		},
		{
			desc: "base type",
			node: ast.BaseType{ID: ast.BoolTypeID},
			visits: []visit{
				{node: ast.BaseType{ID: ast.BoolTypeID}},
			},
		},
		func() (tt test) {
			tt.desc = "base type with annotations"

			foo := &ast.Annotation{Name: "foo", Value: "bar"}
			bar := &ast.Annotation{Name: "baz", Value: "qux"}
			baseType := ast.BaseType{
				ID:          ast.I64TypeID,
				Annotations: []*ast.Annotation{foo, bar},
			}
			tt.node = baseType
			tt.visits = []visit{
				{node: baseType},
				{node: foo, parent: baseType, ancestors: []ast.Node{baseType}},
				{node: bar, parent: baseType, ancestors: []ast.Node{baseType}},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "constant"

			typ := ast.TypeReference{Name: "bar"}
			val := ast.ConstantReference{Name: "baz"}
			constant := &ast.Constant{Name: "foo", Type: typ, Value: val}

			tt.node = constant
			tt.visits = []visit{
				{node: constant},
				{node: typ, parent: constant, ancestors: []ast.Node{constant}},
				{node: val, parent: constant, ancestors: []ast.Node{constant}},
			}
			return
		}(),
		{
			desc: "constant bool",
			node: ast.ConstantBoolean(true),
			visits: []visit{
				{node: ast.ConstantBoolean(true)},
			},
		},
		{
			desc: "constant double",
			node: ast.ConstantDouble(42.0),
			visits: []visit{
				{node: ast.ConstantDouble(42.0)},
			},
		},
		{
			desc: "constant integer",
			node: ast.ConstantInteger(42),
			visits: []visit{
				{node: ast.ConstantInteger(42)},
			},
		},
		{
			desc: "constant list (empty)",
			node: ast.ConstantList{},
			visits: []visit{
				{node: ast.ConstantList{}},
			},
		},
		func() (tt test) {
			tt.desc = "constant list"

			one := ast.ConstantInteger(1)
			two := ast.ConstantInteger(2)
			three := ast.ConstantInteger(3)
			clist := ast.ConstantList{
				Items: []ast.ConstantValue{one, two, three},
			}

			tt.node = clist
			tt.visits = []visit{
				{node: clist},
				{node: one, parent: clist, ancestors: []ast.Node{clist}},
				{node: two, parent: clist, ancestors: []ast.Node{clist}},
				{node: three, parent: clist, ancestors: []ast.Node{clist}},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "constant map"

			item1 := ast.ConstantMapItem{
				Key:   ast.ConstantString("foo"),
				Value: ast.ConstantInteger(1),
			}
			item2 := ast.ConstantMapItem{
				Key:   ast.ConstantString("bar"),
				Value: ast.ConstantInteger(2),
			}
			item3 := ast.ConstantMapItem{
				Key:   ast.ConstantString("baz"),
				Value: ast.ConstantInteger(3),
			}
			cmap := ast.ConstantMap{Items: []ast.ConstantMapItem{item1, item2, item3}}

			tt.node = cmap
			tt.visits = []visit{
				{node: cmap},

				{
					node:      item1,
					parent:    cmap,
					ancestors: []ast.Node{cmap},
				},

				{
					node:      ast.ConstantString("foo"),
					parent:    item1,
					ancestors: []ast.Node{item1, cmap},
				},

				{
					node:      ast.ConstantInteger(1),
					parent:    item1,
					ancestors: []ast.Node{item1, cmap},
				},

				{
					node:      item2,
					parent:    cmap,
					ancestors: []ast.Node{cmap},
				},
				{
					node:      ast.ConstantString("bar"),
					parent:    item2,
					ancestors: []ast.Node{item2, cmap},
				},
				{
					node:      ast.ConstantInteger(2),
					parent:    item2,
					ancestors: []ast.Node{item2, cmap},
				},

				{
					node:      item3,
					parent:    cmap,
					ancestors: []ast.Node{cmap},
				},
				{
					node:      ast.ConstantString("baz"),
					parent:    item3,
					ancestors: []ast.Node{item3, cmap},
				},
				{
					node:      ast.ConstantInteger(3),
					parent:    item3,
					ancestors: []ast.Node{item3, cmap},
				},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "constant map item"

			key := ast.ConstantString("foo")
			value := ast.ConstantDouble(42.0)
			item := ast.ConstantMapItem{Key: key, Value: value}

			tt.node = item
			tt.visits = []visit{
				{node: item},
				{node: key, parent: item, ancestors: []ast.Node{item}},
				{node: value, parent: item, ancestors: []ast.Node{item}},
			}
			return
		}(),
		{
			desc: "constant reference",
			node: ast.ConstantReference{Name: "foo"},
			visits: []visit{
				{node: ast.ConstantReference{Name: "foo"}},
			},
		},
		{
			desc: "constant string",
			node: ast.ConstantString("foo"),
			visits: []visit{
				{node: ast.ConstantString("foo")},
			},
		},
		func() (tt test) {
			tt.desc = "enum"

			item1 := &ast.EnumItem{Name: "foo"}
			item2 := &ast.EnumItem{
				Name: "bar",
				Annotations: []*ast.Annotation{
					{Name: "k1", Value: "v1"},
					{Name: "k2", Value: "v2"},
					{Name: "k3", Value: "v3"},
				},
			}
			item3 := &ast.EnumItem{Name: "baz"}
			enum := &ast.Enum{
				Name:  "e",
				Items: []*ast.EnumItem{item1, item2, item3},
				Annotations: []*ast.Annotation{
					{Name: "k4", Value: "v4"},
				},
			}

			tt.node = enum
			tt.visits = []visit{
				{node: enum},
				{node: item1, parent: enum, ancestors: []ast.Node{enum}},
				{node: item2, parent: enum, ancestors: []ast.Node{enum}},
				{
					node:      &ast.Annotation{Name: "k1", Value: "v1"},
					parent:    item2,
					ancestors: []ast.Node{item2, enum},
				},
				{
					node:      &ast.Annotation{Name: "k2", Value: "v2"},
					parent:    item2,
					ancestors: []ast.Node{item2, enum},
				},
				{
					node:      &ast.Annotation{Name: "k3", Value: "v3"},
					parent:    item2,
					ancestors: []ast.Node{item2, enum},
				},
				{node: item3, parent: enum, ancestors: []ast.Node{enum}},
				{
					node:      &ast.Annotation{Name: "k4", Value: "v4"},
					parent:    enum,
					ancestors: []ast.Node{enum},
				},
			}
			return
		}(),
		{
			desc: "enum item",
			node: &ast.EnumItem{Name: "foo"},
			visits: []visit{
				{node: &ast.EnumItem{Name: "foo"}},
			},
		},
		func() (tt test) {
			tt.desc = "enum item with annotations"

			ann := &ast.Annotation{Name: "k1", Value: "v1"}
			item := &ast.EnumItem{Name: "foo", Annotations: []*ast.Annotation{ann}}

			tt.node = item
			tt.visits = []visit{
				{node: item},
				{node: ann, parent: item, ancestors: []ast.Node{item}},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "field"

			field := &ast.Field{
				ID:           42,
				Name:         "foo",
				Type:         ast.BaseType{ID: ast.BoolTypeID},
				Requiredness: ast.Required,
			}

			tt.node = field
			tt.visits = []visit{
				{node: field},
				{
					node:      ast.BaseType{ID: ast.BoolTypeID},
					parent:    field,
					ancestors: []ast.Node{field},
				},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "field with default"

			field := &ast.Field{
				ID:           2,
				Name:         "bar",
				Type:         ast.BaseType{ID: ast.StringTypeID},
				Requiredness: ast.Required,
				Default:      ast.ConstantString("hi"),
			}

			tt.node = field
			tt.visits = []visit{
				{node: field},
				{
					node:      ast.BaseType{ID: ast.StringTypeID},
					parent:    field,
					ancestors: []ast.Node{field},
				},
				{
					node:      ast.ConstantString("hi"),
					parent:    field,
					ancestors: []ast.Node{field},
				},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "field with annotations"

			typ := ast.BaseType{
				ID: ast.StringTypeID,
				Annotations: []*ast.Annotation{
					{Name: "k1", Value: "v1"},
					{Name: "k2", Value: "v2"},
				},
			}
			field := &ast.Field{
				ID:           2,
				Name:         "bar",
				Type:         typ,
				Requiredness: ast.Required,
				Default:      ast.ConstantString("hi"),
				Annotations: []*ast.Annotation{
					{Name: "a", Value: "b"},
					{Name: "c", Value: "d"},
				},
			}

			tt.node = field
			tt.visits = []visit{
				{node: field},
				{node: typ, parent: field, ancestors: []ast.Node{field}},
				{
					node:      &ast.Annotation{Name: "k1", Value: "v1"},
					parent:    typ,
					ancestors: []ast.Node{typ, field},
				},
				{
					node:      &ast.Annotation{Name: "k2", Value: "v2"},
					parent:    typ,
					ancestors: []ast.Node{typ, field},
				},
				{
					node:      ast.ConstantString("hi"),
					parent:    field,
					ancestors: []ast.Node{field},
				},
				{
					node:      &ast.Annotation{Name: "a", Value: "b"},
					parent:    field,
					ancestors: []ast.Node{field},
				},
				{
					node:      &ast.Annotation{Name: "c", Value: "d"},
					parent:    field,
					ancestors: []ast.Node{field},
				},
			}
			return
		}(),
		{
			desc: "function with nothing",
			node: &ast.Function{Name: "noop"},
			visits: []visit{
				{node: &ast.Function{Name: "noop"}},
			},
		},
		func() (tt test) {
			tt.desc = "function with everything"

			keyType := ast.BaseType{
				ID: ast.StringTypeID,
				Annotations: []*ast.Annotation{
					{Name: "validator", Value: "alphanumeric"},
				},
			}
			key := &ast.Field{
				ID:   1,
				Name: "key",
				Type: keyType,
				Annotations: []*ast.Annotation{
					{Name: "http.param", Value: "key"},
				},
			}

			value := &ast.Field{
				ID:          2,
				Name:        "value",
				Type:        ast.TypeReference{Name: "Value"},
				Annotations: []*ast.Annotation{{Name: "http.body"}},
			}

			doesNotExist := &ast.Field{
				ID:   1,
				Name: "doesNotExist",
				Type: ast.TypeReference{Name: "DoesNotExistException"},
				Annotations: []*ast.Annotation{
					{Name: "http.status", Value: "404"},
				},
			}

			function := &ast.Function{
				Name:       "getAndSet",
				Parameters: []*ast.Field{key, value},
				ReturnType: ast.TypeReference{Name: "Value"},
				Exceptions: []*ast.Field{doesNotExist},
				Annotations: []*ast.Annotation{
					{Name: "http.url", Value: "/update/:key"},
					{Name: "http.method", Value: "POST"},
				},
			}

			tt.node = function
			tt.visits = []visit{
				{node: function},

				// Return type
				{
					node:      ast.TypeReference{Name: "Value"},
					parent:    function,
					ancestors: []ast.Node{function},
				},

				// Key param
				{node: key, parent: function, ancestors: []ast.Node{function}},
				{node: keyType, parent: key, ancestors: []ast.Node{key, function}},
				{
					node:      &ast.Annotation{Name: "validator", Value: "alphanumeric"},
					parent:    keyType,
					ancestors: []ast.Node{keyType, key, function},
				},
				{
					node:      &ast.Annotation{Name: "http.param", Value: "key"},
					parent:    key,
					ancestors: []ast.Node{key, function},
				},

				// Value param
				{node: value, parent: function, ancestors: []ast.Node{function}},
				{
					node:      ast.TypeReference{Name: "Value"},
					parent:    value,
					ancestors: []ast.Node{value, function},
				},
				{
					node:      &ast.Annotation{Name: "http.body"},
					parent:    value,
					ancestors: []ast.Node{value, function},
				},

				// Exception
				{node: doesNotExist, parent: function, ancestors: []ast.Node{function}},
				{
					node:      ast.TypeReference{Name: "DoesNotExistException"},
					parent:    doesNotExist,
					ancestors: []ast.Node{doesNotExist, function},
				},
				{
					node:      &ast.Annotation{Name: "http.status", Value: "404"},
					parent:    doesNotExist,
					ancestors: []ast.Node{doesNotExist, function},
				},

				// Annotations
				{
					node:      &ast.Annotation{Name: "http.url", Value: "/update/:key"},
					parent:    function,
					ancestors: []ast.Node{function},
				},
				{
					node:      &ast.Annotation{Name: "http.method", Value: "POST"},
					parent:    function,
					ancestors: []ast.Node{function},
				},
			}
			return
		}(),
		{
			desc: "include",
			node: &ast.Include{Path: "foo.thrift"},
			visits: []visit{
				{node: &ast.Include{Path: "foo.thrift"}},
			},
		},
		func() (tt test) {
			tt.desc = "list type"

			itemType := ast.BaseType{ID: ast.I64TypeID}
			listType := ast.ListType{ValueType: itemType}

			tt.node = listType
			tt.visits = []visit{
				{node: listType},
				{node: itemType, parent: listType, ancestors: []ast.Node{listType}},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "list type annotations"

			itemType := ast.BaseType{
				ID: ast.I64TypeID,
				Annotations: []*ast.Annotation{
					{Name: "foo", Value: "bar"},
				},
			}
			listType := ast.ListType{
				ValueType: itemType,
				Annotations: []*ast.Annotation{
					{Name: "baz", Value: "qux"},
				},
			}

			tt.node = listType
			tt.visits = []visit{
				{node: listType},
				{node: itemType, parent: listType, ancestors: []ast.Node{listType}},
				{
					node:      &ast.Annotation{Name: "foo", Value: "bar"},
					parent:    itemType,
					ancestors: []ast.Node{itemType, listType},
				},
				{
					node:      &ast.Annotation{Name: "baz", Value: "qux"},
					parent:    listType,
					ancestors: []ast.Node{listType},
				},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "map type"

			keyType := ast.BaseType{ID: ast.StringTypeID}
			valueType := ast.BaseType{ID: ast.BinaryTypeID}
			mapType := ast.MapType{KeyType: keyType, ValueType: valueType}

			tt.node = mapType
			tt.visits = []visit{
				{node: mapType},
				{node: keyType, parent: mapType, ancestors: []ast.Node{mapType}},
				{node: valueType, parent: mapType, ancestors: []ast.Node{mapType}},
			}
			return
		}(),
		func() (tt test) {
			tt.desc = "map type with annotations"

			keyType := ast.BaseType{ID: ast.StringTypeID}
			valueType := ast.BaseType{ID: ast.BinaryTypeID}
			mapType := ast.MapType{
				KeyType:   keyType,
				ValueType: valueType,
				Annotations: []*ast.Annotation{
					{Name: "foo", Value: "bar"},
					{Name: "baz", Value: "qux"},
				},
			}

			tt.node = mapType
			tt.visits = []visit{
				{node: mapType},
				{node: keyType, parent: mapType, ancestors: []ast.Node{mapType}},
				{node: valueType, parent: mapType, ancestors: []ast.Node{mapType}},
				{
					node:      &ast.Annotation{Name: "foo", Value: "bar"},
					parent:    mapType,
					ancestors: []ast.Node{mapType},
				},
				{
					node:      &ast.Annotation{Name: "baz", Value: "qux"},
					parent:    mapType,
					ancestors: []ast.Node{mapType},
				},
			}
			return
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			v := ast.NewMockVisitor(mockCtrl)

			var calls []*gomock.Call
			for _, visit := range tt.visits {
				call := v.EXPECT().
					Visit(
						walkerMatcher{
							Parent:    visit.parent,
							Ancestors: visit.ancestors,
						},
						visit.node,
					).Return(v)
				calls = append(calls, call)
			}
			gomock.InOrder(calls...)

			ast.Walk(v, tt.node)
		})
	}
}

type walkerMatcher struct {
	Ancestors []ast.Node
	Parent    ast.Node
}

var _ gomock.Matcher = walkerMatcher{}

func (m walkerMatcher) Matches(x interface{}) bool {
	w, ok := x.(ast.Walker)
	if !ok {
		return false
	}

	return reflect.DeepEqual(m.Parent, w.Parent()) &&
		reflect.DeepEqual(m.Ancestors, w.Ancestors())
}

func (m walkerMatcher) String() string {
	buff := bytes.NewBufferString("Walker{")
	if m.Parent != nil {
		fmt.Fprintf(buff, "Parent: %#v", m.Parent)
		if len(m.Ancestors) > 0 {
			buff.WriteString(", ")
		}
	}

	if len(m.Ancestors) > 0 {
		buff.WriteString("Ancestors: [")

		first := true
		for _, n := range m.Ancestors {
			if first {
				first = false
			} else {
				buff.WriteString(", ")
			}
			fmt.Fprintf(buff, "%#v", n)
		}

		buff.WriteString("]")
	}
	buff.WriteString("}")

	return buff.String()
}
