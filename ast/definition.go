// Copyright (c) 2019 Uber Technologies, Inc.
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

package ast

// DefinitionInfo provides a common way to access name and line information
// for definitions.
type DefinitionInfo struct {
	Name string
	Line int
}

// Definition unifies the different types representing items defined in the
// Thrift file.
type Definition interface {
	Node

	Info() DefinitionInfo
	definition()
}

// Constant is a constant declared in the Thrift file using a const statement.
//
// 	const i32 foo = 42
type Constant struct {
	Name  string
	Type  Type
	Value ConstantValue
	Line  int
	Doc   string
}

func (*Constant) node()       {}
func (*Constant) definition() {}

func (c *Constant) lineNumber() int { return c.Line }

func (c *Constant) visitChildren(ss nodeStack, v visitor) {
	v.visit(ss, c.Type)
	v.visit(ss, c.Value)
}

// Info for Constant
func (c *Constant) Info() DefinitionInfo {
	return DefinitionInfo{Name: c.Name, Line: c.Line}
}

// Typedef is used to define an alias for another type.
//
// 	typedef string UUID
// 	typedef i64 Timestamp (unit = "milliseconds")
type Typedef struct {
	Name        string
	Type        Type
	Annotations []*Annotation
	Line        int
	Doc         string
}

// Definition implementation for Typedef.
func (*Typedef) node()       {}
func (*Typedef) definition() {}

func (t *Typedef) lineNumber() int { return t.Line }

func (t *Typedef) visitChildren(ss nodeStack, v visitor) {
	v.visit(ss, t.Type)
	for _, ann := range t.Annotations {
		v.visit(ss, ann)
	}
}

// Info for Typedef.
func (t *Typedef) Info() DefinitionInfo {
	return DefinitionInfo{Name: t.Name, Line: t.Line}
}

// Enum is a set of named integer values.
//
// 	enum Status { Enabled, Disabled }
//
// 	enum Role {
// 		User = 1,
// 		Moderator = 2 (py.name = "Mod"),
// 		Admin = 3
// 	} (go.name = "UserRole")
type Enum struct {
	Name        string
	Items       []*EnumItem
	Annotations []*Annotation
	Line        int
	Doc         string
}

func (*Enum) node()       {}
func (*Enum) definition() {}

func (e *Enum) lineNumber() int { return e.Line }

func (e *Enum) visitChildren(ss nodeStack, v visitor) {
	for _, item := range e.Items {
		v.visit(ss, item)
	}

	for _, ann := range e.Annotations {
		v.visit(ss, ann)
	}
}

// Info for Enum.
func (e *Enum) Info() DefinitionInfo {
	return DefinitionInfo{Name: e.Name, Line: e.Line}
}

// EnumItem is a single item in an Enum definition.
type EnumItem struct {
	Name string
	// Value of the item. This is nil if the user did not specify anything.
	Value       *int
	Annotations []*Annotation
	Line        int
	Doc         string
}

func (*EnumItem) node() {}

func (i *EnumItem) lineNumber() int { return i.Line }

func (i *EnumItem) visitChildren(ss nodeStack, v visitor) {
	for _, ann := range i.Annotations {
		v.visit(ss, ann)
	}
}

// StructureType specifies whether a struct-like type is a struct, union, or
// exception.
type StructureType int

// Different kinds of struct-like objects supported by us.
const (
	StructType    StructureType = iota + 1 // struct
	UnionType                              // union
	ExceptionType                          // exception
)

// Struct is a collection of named fields with different types.
//
// This type encompasses structs, unions, and exceptions.
//
// 	struct User {
// 		1: required string name (min_length = "3")
// 		2: optional Status status = Enabled;
// 	}
//
// 	struct i128 {
// 		1: required i64 high
// 		2: required i64 low
// 	} (py.serializer = "foo.Int128Serializer")
//
// 	union Contents {
// 		1: string plainText
// 		2: binary pdf
// 	}
//
// 	exception ServiceError { 1: required string message }
type Struct struct {
	Name        string
	Type        StructureType
	Fields      []*Field
	Annotations []*Annotation
	Line        int
	Doc         string
}

func (*Struct) node()       {}
func (*Struct) definition() {}

func (s *Struct) lineNumber() int { return s.Line }

func (s *Struct) visitChildren(ss nodeStack, v visitor) {
	for _, field := range s.Fields {
		v.visit(ss, field)
	}
	for _, ann := range s.Annotations {
		v.visit(ss, ann)
	}
}

// Info for Struct.
func (s *Struct) Info() DefinitionInfo {
	return DefinitionInfo{Name: s.Name, Line: s.Line}
}

// Service is a collection of functions.
//
// 	service KeyValue {
// 		void setValue(1: string key, 2: binary value)
// 		binary getValue(1: string key)
// 	} (router.serviceName = "key_value")
type Service struct {
	Name      string
	Functions []*Function
	// Reference to the parent service if this service inherits another
	// service, nil otherwise.
	Parent      *ServiceReference
	Annotations []*Annotation
	Line        int
	Doc         string
}

func (*Service) node()       {}
func (*Service) definition() {}

func (s *Service) lineNumber() int { return s.Line }

func (s *Service) visitChildren(ss nodeStack, v visitor) {
	for _, function := range s.Functions {
		v.visit(ss, function)
	}
	for _, ann := range s.Annotations {
		v.visit(ss, ann)
	}
}

// Info for Service.
func (s *Service) Info() DefinitionInfo {
	return DefinitionInfo{Name: s.Name, Line: s.Line}
}

// Function is a single function inside a service.
//
// 	binary getValue(1: string key)
// 		throws (1: KeyNotFoundError notFound) (
// 			ttl.milliseconds = "250"
// 		)
type Function struct {
	Name        string
	Parameters  []*Field
	ReturnType  Type
	Exceptions  []*Field
	OneWay      bool
	Annotations []*Annotation
	Line        int
	Doc         string
}

func (*Function) node() {}

func (n *Function) lineNumber() int { return n.Line }

func (n *Function) visitChildren(ss nodeStack, v visitor) {
	v.visit(ss, n.ReturnType)
	for _, field := range n.Parameters {
		v.visit(ss, field)
	}
	for _, exc := range n.Exceptions {
		v.visit(ss, exc)
	}
	for _, ann := range n.Annotations {
		v.visit(ss, ann)
	}
}

// Requiredness represents whether a field was marked as required or optional,
// or if the user did not specify either.
type Requiredness int

// Different requiredness levels that are supported.
const (
	Unspecified Requiredness = iota // unspecified (default)
	Required                        // required
	Optional                        // optional
)

// Field is a single field inside a struct, union, exception, or a single item
// in the parameter or exception list of a function.
//
// 	1: required i32 foo = 0
// 	2: optional binary (max_length = "4096") bar
// 	3: i64 baz (go.name = "qux")
//
type Field struct {
	ID           int
	Name         string
	Type         Type
	Requiredness Requiredness
	Default      ConstantValue
	Annotations  []*Annotation
	Line         int
	Doc          string
}

func (*Field) node() {}

func (n *Field) lineNumber() int { return n.Line }

func (n *Field) visitChildren(ss nodeStack, v visitor) {
	v.visit(ss, n.Type)
	v.visit(ss, n.Default)
	for _, ann := range n.Annotations {
		v.visit(ss, ann)
	}
}

// ServiceReference is a reference to another service.
type ServiceReference struct {
	Name string
	Line int
}
