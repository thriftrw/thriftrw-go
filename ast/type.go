package ast

// Type unifies the different types representing Thrift field types.
type Type interface {
	fieldType()
}

func (BaseType) fieldType()      {}
func (MapType) fieldType()       {}
func (ListType) fieldType()      {}
func (SetType) fieldType()       {}
func (TypeReference) fieldType() {}

// BaseTypeID is an identifier for primitive types supported by Thrift.
type BaseTypeID int

//go:generate stringer -type=BaseTypeID

// IDs of the base types supported by Thrift.
const (
	BoolBaseTypeID   BaseTypeID = iota + 1 // bool
	ByteBaseTypeID                         // byte
	I16BaseTypeID                          // i16
	I32BaseTypeID                          // i32
	I64BaseTypeID                          // i64
	DoubleBaseTypeID                       // double
	StringBaseTypeID                       // string
	BinaryBaseTypeID                       // binary
)

// BaseType is a reference to a Thrift base type.
//
// 	bool, byte, i16, i32, i64, double, string, binary
//
// All references to base types in the document may be followed by type
// annotations.
//
// 	bool (go.type = "int")
type BaseType struct {
	// ID of the base type.
	ID BaseTypeID

	// Type annotations associated with this reference.
	Annotations []*Annotation
}

// MapType is a reference to a the Thrift map type.
//
// 	map<k, v>
//
// All references to map types may be followed by type annotations.
//
// 	map<string, list<i32>> (java.type = "MultiMap")
type MapType struct {
	KeyType, ValueType Type
	Annotations        []*Annotation
}

// ListType is a reference to the Thrift list type.
//
// 	list<a>
//
// All references to list types may be followed by type annotations.
//
// 	list<i64> (cpp.type = "vector")
type ListType struct {
	ValueType   Type
	Annotations []*Annotation
}

// SetType is a reference to the Thrift set type.
//
// 	set<a>
//
// All references to set types may be followed by type annotations.
//
// 	set<string> (js.type = "list")
type SetType struct {
	ValueType   Type
	Annotations []*Annotation
}

// TypeReference references a user-defined type.
type TypeReference struct {
	Name string
	Line int
}
