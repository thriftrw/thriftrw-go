include "./enums.thrift"

struct EmptyStruct {}

//////////////////////////////////////////////////////////////////////////////
// Structs with primitives

/**
 * A struct that contains primitive fields exclusively.
 *
 * All fields are required.
 */
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

/**
 * A struct that contains primitive fields exclusively.
 *
 * All fields are optional.
 */
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

/**
 * A point in 2D space.
 */
struct Point {
    1: required double x
    2: required double y
}

/**
 * Size of something.
 */
struct Size {
    /**
     * Width in pixels.
     */
    1: required double width
    /** Height in pixels. */
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

/**
 * A graph is comprised of zero or more edges.
 */
struct Graph {
    /**
     * List of edges in the graph.
     *
     * May be empty.
     */
    1: required list<Edge> edges
}

//////////////////////////////////////////////////////////////////////////////
// Nested structs (Optional)

struct ContactInfo {
    1: required string emailAddress
}

struct PersonalInfo {
    1: optional i32 age
    2: optional string race (go.redact)
}

struct User {
    1: required string name
    2: optional ContactInfo contact
    3: optional PersonalInfo personal
}

typedef map<string, User> UserMap

//////////////////////////////////////////////////////////////////////////////
// self-referential struct

typedef Node List

/**
 * Node is linked list of values.
 * All values are 32-bit integers.
 */
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

struct NotOmitEmpty {
    1: optional string NotOmitEmptyString (go.tag = 'json:"notOmitEmptyString,!omitempty"')
    2: optional string NotOmitEmptyInt (go.tag = 'json:"notOmitEmptyInt,!omitempty"')
    3: optional string NotOmitEmptyBool (go.tag = 'json:"notOmitEmptyBool,!omitempty"')
    4: optional list<string> NotOmitEmptyList (go.tag = 'json:"notOmitEmptyList,!omitempty"')
    5: optional map<string, string> NotOmitEmptyMap (go.tag = 'json:"notOmitEmptyMap,!omitempty"')
    6: optional list<string> NotOmitEmptyListMixedWithOmitEmpty (go.tag = 'json:"notOmitEmptyListMixedWithOmitEmpty,!omitempty,omitempty"')
    7: optional list<string> NotOmitEmptyListMixedWithOmitEmptyV2 (go.tag = 'json:"notOmitEmptyListMixedWithOmitEmptyV2,omitempty,!omitempty"')
    8: optional string OmitEmptyString (go.tag = 'json:"omitEmptyString,omitempty"') // to test that there can be a mix of fields that do and don't have !omitempty
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

    9:  required bool requiredBoolDefaultTrue = true
    10: optional bool optionalBoolDefaultTrue = true

    11: required bool requiredBoolDefaultFalse = false
    12: optional bool optionalBoolDefaultFalse = false
}

//////////////////////////////////////////////////////////////////////////////
// Opt-out of Zap

struct ZapOptOutStruct {
    1: required string name
    2: required string optout (go.nolog)
}

//////////////////////////////////////////////////////////////////////////////
// Field jabels

struct StructLabels {
    // reserved keyword as label
    1: optional bool isRequired (go.label = "required")

    // go.tag's JSON tag takes precedence over go.label
    2: optional string foo (go.label = "bar", go.tag = 'json:"not_bar"')

    // Empty label
    3: optional string qux (go.label = "")

    // All-caps label
    4: optional string quux (go.label = "QUUX")
}
