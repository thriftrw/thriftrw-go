include "./structs.thrift"

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

typedef set<structs.Frame> FrameGroup

typedef map<structs.Point, structs.Point> PointMap
