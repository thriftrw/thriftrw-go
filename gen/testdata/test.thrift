//////////////////////////////////////////////////////////////////////////////
// Containers

struct PrimitiveContainers {
    1: optional list<binary> listOfBinary
    2: optional list<i64> listOfInts
    3: optional set<string> setOfStrings
    4: optional set<byte> setOfBytes
    5: optional map<i32, string> mapOfIntToString
    6: optional map<string, bool> mapOfStringToBool
}

struct PrimitiveContainersRequired {
    1: required list<string> listOfStrings
    2: required set<i32> setOfInts
    3: required map<i64, double> mapOfIntsToDoubles
}

struct EnumContainers {
    1: optional list<EnumDefault> listOfEnums
    2: optional set<EnumWithValues> setOfEnums
    3: optional map<EnumWithDuplicateValues, i32> mapOfEnums
}

//////////////////////////////////////////////////////////////////////////////
// Enums

enum EnumDefault {
    Foo, Bar, Baz
}

enum EnumWithValues {
    X = 123,
    Y = 456,
    Z = 789,
}

enum EnumWithDuplicateValues {
    P, // 0
    Q = -1,
    R, // 0
}

// enum with item names conflicting with those of another enum
enum EnumWithDuplicateName {
    A, B, C, P, Q, R, X, Y, Z
}

//////////////////////////////////////////////////////////////////////////////
// Structs with primitives

struct PrimitiveRequiredStruct {
    1: required bool boolField
    2: required byte byteField
    3: required i16 int16Field
    4: required i32 int32Field
    5: required i64 int64Field
    6: required double doubleField
    7: required string stringField
    8: required binary binaryField
}

struct PrimitiveOptionalStruct {
    1: optional bool boolField
    2: optional byte byteField
    3: optional i16 int16Field
    4: optional i32 int32Field
    5: optional i64 int64Field
    6: optional double doubleField
    7: optional string stringField
    8: optional binary binaryField
}

//////////////////////////////////////////////////////////////////////////////
// Nested structs (Required)

struct Point {
    1: required double x
    2: required double y
}

struct Size {
    1: required double width
    2: required double height
}

struct Frame {
    1: required Point topLeft
    2: required Size size
}

struct Edge {
    1: required Point start
    2: required Point end
}

struct Graph {
    1: required list<Edge> edges
}

//////////////////////////////////////////////////////////////////////////////
// Nested structs (Optional)

struct ContactInfo {
    1: required string emailAddress
}

struct User {
    1: required string name
    2: optional ContactInfo contact
}

// TODO: Default values

// TODO: unions
// TODO: exceptions

//////////////////////////////////////////////////////////////////////////////
// Typedefs

typedef i64 Timestamp  // alias of primitive
typedef string State

typedef i128 UUID  // alias of struct

typedef list<Event> EventGroup  // alias fo collection

struct i128 {
    1: required i64 high
    2: required i64 low
}

struct Event {
    1: required UUID uuid  // required typedef
    2: optional Timestamp time  // optional typedef
}

struct Transition {
    1: required State from
    2: required State to
    3: optional EventGroup events
}

typedef binary PDF  // alias of []byte
