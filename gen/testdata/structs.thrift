// TODO: Default values

//////////////////////////////////////////////////////////////////////////////
// Primitives

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
// Containers

// Optional containers have already been tested by containers.thrift.
struct ContainersRequiredStruct {
    1: required list<string> listOfStrings
    2: required set<i32> setOfInts
    3: required map<i64, double> mapOfIntsToDoubles
}

//////////////////////////////////////////////////////////////////////////////
// Required nested structs

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

//////////////////////////////////////////////////////////////////////////////
// Optional nested structS:

struct ContactInfo {
    1: required string emailAddress
}

struct User {
    1: required string name
    2: optional ContactInfo contact
}
