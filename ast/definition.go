package ast

// Definition unifies the different types representing items defined in the
// Thrift file.
type Definition interface {
	definition()
}

func (*Constant) definition() {}
func (*Enum) definition()     {}
func (*Service) definition()  {}
func (*Struct) definition()   {}
func (*Typedef) definition()  {}

// Constant is a constant declared in the Thrift file using a const statement.
//
// 	const i32 foo = 42
type Constant struct {
	Name  string
	Type  Type
	Value ConstantValue
	Line  int
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
}

// EnumItem is a single item in an Enum definition.
type EnumItem struct {
	Name string
	// Value of the item. This is nil if the user did not specify anything.
	Value       *int
	Annotations []*Annotation
	Line        int
}

// StructType specifies whether a struct-like type is a struct, union, or
// exception.
type StructureType int

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
}

// Requiredness represents whether a field was marked as required or optional,
// or if the user did not specify either.
type Requiredness int

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
