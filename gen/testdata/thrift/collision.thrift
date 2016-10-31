
struct StructCollision {
	1: required bool collisionField
	2: required string collision_field (go.name = "CollisionField2")
}

struct struct_collision {
	1: required bool collisionField
	2: required string collision_field (go.name = "CollisionField2")
} (go.name="StructCollision2")

struct PrimitiveContainers {
    1: optional list<string> ListOrSetOrMap (go.name = "A")
    3: optional set<string>  List_Or_SetOrMap (go.name = "B")
    5: optional map<string, string> ListOrSet_Or_Map (go.name = "C")
}

enum MyEnum {
    X = 123,
    Y = 456,
    Z = 789,
    FooBar,
    foo_bar (go.name="FooBar2"),
}

enum my_enum {
    X = 12,
    Y = 34,
    Z = 56,
} (go.name="MyEnum2")

typedef i64 LittlePotatoe
typedef double little_potatoe (go.name="LittlePotatoe2")

const struct_collision struct_constant = {
	"collisionField": false,
	"collision_field": "false indeed",
}

union UnionCollision {
	1: bool collisionField
	2: string collision_field (go.name = "CollisionField2")
}

union union_collision {
	1: bool collisionField
	2: string collision_field (go.name = "CollisionField2")
} (go.name="UnionCollision2")

struct WithDefault {
	1: required struct_collision pouet = struct_constant
}
