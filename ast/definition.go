package ast

// Definition unifies the different types representing items defined in the
// Thrift file.
type Definition interface {
	definition()
}

func (*Constant) definition()  {}
func (*Enum) definition()      {}
func (*Exception) definition() {}
func (*Service) definition()   {}
func (*Struct) definition()    {}
func (*Typedef) definition()   {}
func (*Union) definition()     {}

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

// Struct is a collection of named fields with different types.
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
type Struct struct {
	Name       string
	Fields     []*Field
	Annotation []*Annotation
	Line       int
}

// Union is similar to a Struct except only a single field of the union can be
// populated at any time.
//
// 	union Contents {
// 		1: string plainText
// 		2: binary pdf
// 	}
type Union struct {
	Struct
}

// Exception is similar to a Struct, but it represents error types that may be
// returned by methods in case of failure.
//
// 	exception ServiceError { 1: required string message }
type Exception struct {
	Struct
}

// Service is a collection of functions.
//
// 	service KeyValue {
// 		void setValue(1: string key, 2: binary value)
// 		binary getValue(1: string key)
// 	} (router.serviceName = "key_value")
type Service struct {
	Name     string
	Function []*Function
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
}

// Field is a single field inside a struct, union, exception, or a single item
// in the parameter or exception list of a function.
//
// 	1: required i32 foo = 0
// 	2: optional binary (max_length = "4096") bar
// 	3: i64 baz (go.name = "qux")
//
type Field struct {
	ID   int
	Name string
	Type Type
	// Requiredness may be true or false, or nil if the optional/required
	// wasn't specified.
	Requiredness *bool
	DefaultValue *ConstantValue
	Annotations  []*Annotation
	Line         int
}

// ServiceReference is a reference to another service.
type ServiceReference struct {
	Name string
	Line int
}
