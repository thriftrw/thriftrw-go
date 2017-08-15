include "./enums.thrift"

struct EmptyStruct {}

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
    1: required Point startPoint
    2: required Point endPoint
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

//////////////////////////////////////////////////////////////////////////////
// self-referential struct

typedef Node List

struct Node {
    1: required i32 value
    2: optional List tail
}

//////////////////////////////////////////////////////////////////////////////
// JSON tagged structs

struct Rename {
    1: required string Default (go.tag = 'json:"default"')
    2: required string camelCase (go.tag = 'json:"snake_case"')
}

struct Omit {
    1: required string serialized
    2: required string hidden (go.tag = 'json:"-"')
}

struct GoTags {
        1: required string Foo (go.tag = 'json:"-" foo:"bar"')
        2: optional string Bar (go.tag = 'bar:"foo"')
        3: required string FooBar (go.tag = 'json:"foobar,option1,option2" bar:"foo,option1" foo:"foobar"')
        4: required string FooBarWithSpace (go.tag = 'json:"foobarWithSpace" foo:"foo bar foobar barfoo"')
        5: optional string FooBarWithOmitEmpty (go.tag = 'json:"foobarWithOmitEmpty,omitempty"')
        6: required string FooBarWithRequired (go.tag = 'json:"foobarWithRequired,required"')
}

//////////////////////////////////////////////////////////////////////////////
// Default values

struct DefaultsStruct {
    1: required i32 requiredPrimitive = 100
    2: optional i32 optionalPrimitive = 200

    3: required enums.EnumDefault requiredEnum = enums.EnumDefault.Bar
    4: optional enums.EnumDefault optionalEnum = 2

    5: required list<string> requiredList = ["hello", "world"]
    6: optional list<double> optionalList = [1, 2.0, 3]

    7: required Frame requiredStruct = {
        "topLeft": {"x": 1, "y": 2},
        "size": {"width": 100, "height": 200},
    }
    8: optional Edge optionalStruct = {
        "startPoint": {"x": 1, "y": 2},
        "endPoint":   {"x": 3, "y": 4},
    }
}
