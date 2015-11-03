package ast

// ConstantValue unifies the different types representing constant values in
// Thrift files.
type ConstantValue interface {
	constantValue()
}

func (ConstantBoolean) constantValue()   {}
func (ConstantInteger) constantValue()   {}
func (ConstantString) constantValue()    {}
func (ConstantDouble) constantValue()    {}
func (ConstantReference) constantValue() {}
func (ConstantMap) constantValue()       {}
func (ConstantList) constantValue()      {}

// ConstantBoolean is a boolean value specified in the Thrift file.
//
//   true
//   false
type ConstantBoolean bool

// ConstantInteger is an integer value specified in the Thrift file.
//
//   42
type ConstantInteger int64

// ConstantString is a string literal specified in the Thrift file.
//
//   "hello world"
type ConstantString string

// ConstantDouble is a floating point value specified in the Thrift file.
//
//   1.234
type ConstantDouble float64

// ConstantMap is a map literal from the Thrift file.
//
// 	{"a": 1, "b": 2}
//
// Note that map literals can also be used to build structs.
type ConstantMap struct {
	Items []ConstantMapItem
}

// ConstantMapItem is a single item in a ConstantMap.
type ConstantMapItem struct {
	Key, Value ConstantValue
}

// ConstantList is a list literal from the Thrift file.
//
// 	[1, 2, 3]
type ConstantList struct {
	Items []ConstantValue
}

// ConstantReference is a reference to another constant value defined in the
// Thrift file.
//
// 	foo.bar
type ConstantReference struct {
	// Name of the referenced value.
	Name string

	// Line number on which this reference was made.
	Line int
}
