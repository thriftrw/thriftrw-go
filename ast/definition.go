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

package ast

// Definition unifies the different types representing items defined in the
// Thrift file.
type Definition interface {
	Name() string
	Line() int
	definition()
}

// Constant is a constant declared in the Thrift file using a const statement.
//
// 	const i32 foo = 42
type Constant struct {
	CName string
	Type  Type
	Value ConstantValue
	CLine int
}

func (c *Constant) definition()  {}
func (c *Constant) Name() string { return c.CName }
func (c *Constant) Line() int    { return c.CLine }

// Typedef is used to define an alias for another type.
//
// 	typedef string UUID
// 	typedef i64 Timestamp (unit = "milliseconds")
type Typedef struct {
	TName       string
	Type        Type
	Annotations []*Annotation
	TLine       int
}

// Definition implementation for Typedef.
func (t *Typedef) definition()  {}
func (t *Typedef) Name() string { return t.TName }
func (t *Typedef) Line() int    { return t.TLine }

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
	EName       string
	Items       []*EnumItem
	Annotations []*Annotation
	ELine       int
}

// DefinitionName for Enum.
func (e *Enum) definition()  {}
func (e *Enum) Name() string { return e.EName }
func (e *Enum) Line() int    { return e.ELine }

// EnumItem is a single item in an Enum definition.
type EnumItem struct {
	Name string
	// Value of the item. This is nil if the user did not specify anything.
	Value       *int
	Annotations []*Annotation
	Line        int
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
	SName       string
	Type        StructureType
	Fields      []*Field
	Annotations []*Annotation
	SLine       int
}

// DefinitionName implementation for Struct.
func (s *Struct) definition()  {}
func (s *Struct) Name() string { return s.SName }
func (s *Struct) Line() int    { return s.SLine }

// Service is a collection of functions.
//
// 	service KeyValue {
// 		void setValue(1: string key, 2: binary value)
// 		binary getValue(1: string key)
// 	} (router.serviceName = "key_value")
type Service struct {
	SName     string
	Functions []*Function
	// Reference to the parent service if this service inherits another
	// service, nil otherwise.
	Parent      *ServiceReference
	Annotations []*Annotation
	SLine       int
}

// DefinitionName implementation for Service.
func (s *Service) definition()  {}
func (s *Service) Name() string { return s.SName }
func (s *Service) Line() int    { return s.SLine }

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
}

// ServiceReference is a reference to another service.
type ServiceReference struct {
	Name string
	Line int
}
