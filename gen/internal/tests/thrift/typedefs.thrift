include "./structs.thrift"
include "./enums.thrift"

/**
 * Number of seconds since epoch.
 *
 * Deprecated: Use ISOTime instead.
 */
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

struct DefaultPrimitiveTypedef {
    1: optional State state = "hello"
}

struct Transition {
    1: required State fromState
    2: required State toState
    3: optional EventGroup events
}

typedef binary PDF  // alias of []byte

typedef set<structs.Frame> FrameGroup

typedef map<structs.Point, structs.Point> PointMap

typedef set<binary> BinarySet

typedef map<structs.Edge, structs.Edge> EdgeMap

typedef enums.EnumWithValues MyEnum
