package ast_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"go.uber.org/thriftrw/ast"

	"github.com/golang/mock/gomock"
)

func TestWalk(t *testing.T) {
	type visit struct {
		// expected node and that node's ancestors for each visit
		node      ast.Node
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
			node: &ast.BaseType{ID: ast.BoolTypeID},
			visits: []visit{
				{node: &ast.BaseType{ID: ast.BoolTypeID}},
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
				{node: foo, ancestors: []ast.Node{baseType}},
				{node: bar, ancestors: []ast.Node{baseType}},
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
				{node: typ, ancestors: []ast.Node{constant}},
				{node: val, ancestors: []ast.Node{constant}},
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
				{node: one, ancestors: []ast.Node{clist}},
				{node: two, ancestors: []ast.Node{clist}},
				{node: three, ancestors: []ast.Node{clist}},
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

				{node: item1, ancestors: []ast.Node{cmap}},
				{node: ast.ConstantString("foo"), ancestors: []ast.Node{item1, cmap}},
				{node: ast.ConstantInteger(1), ancestors: []ast.Node{item1, cmap}},

				{node: item2, ancestors: []ast.Node{cmap}},
				{node: ast.ConstantString("bar"), ancestors: []ast.Node{item2, cmap}},
				{node: ast.ConstantInteger(2), ancestors: []ast.Node{item2, cmap}},

				{node: item3, ancestors: []ast.Node{cmap}},
				{node: ast.ConstantString("baz"), ancestors: []ast.Node{item3, cmap}},
				{node: ast.ConstantInteger(3), ancestors: []ast.Node{item3, cmap}},
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
				{node: key, ancestors: []ast.Node{item}},
				{node: value, ancestors: []ast.Node{item}},
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
				{node: item1, ancestors: []ast.Node{enum}},
				{node: item2, ancestors: []ast.Node{enum}},
				{
					node:      &ast.Annotation{Name: "k1", Value: "v1"},
					ancestors: []ast.Node{item2, enum},
				},
				{
					node:      &ast.Annotation{Name: "k2", Value: "v2"},
					ancestors: []ast.Node{item2, enum},
				},
				{
					node:      &ast.Annotation{Name: "k3", Value: "v3"},
					ancestors: []ast.Node{item2, enum},
				},
				{node: item3, ancestors: []ast.Node{enum}},
				{
					node:      &ast.Annotation{Name: "k4", Value: "v4"},
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
				{node: ann, ancestors: []ast.Node{item}},
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
					Visit(matchAncestors(visit.ancestors), visit.node).
					Return(v)
				calls = append(calls, call)
			}
			gomock.InOrder(calls...)

			ast.Walk(v, tt.node)
		})
	}
}

type matchAncestors []ast.Node

func (m matchAncestors) Matches(x interface{}) bool {
	w, ok := x.(ast.Walker)
	if !ok {
		return false
	}
	return reflect.DeepEqual([]ast.Node(m), w.Ancestors())
}

func (m matchAncestors) String() string {
	ancestors := make([]string, len(m))
	for i, n := range m {
		ancestors[i] = fmt.Sprint(n)
	}
	return "[" + strings.Join(ancestors, ", ") + "]"
}
