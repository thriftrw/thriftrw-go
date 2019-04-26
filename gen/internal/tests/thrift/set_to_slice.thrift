typedef set<string> (go.type = "slice") StringList
typedef set<Foo> (go.type = "slice") FooList
typedef StringList MyStringList
typedef MyStringList AnotherStringList

typedef set<set<string> (go.type = "slice")> (go.type = "slice") StringListList

struct Foo {
    1: required string stringField
}

struct Bar {
    1: required set<i32> (go.type = "slice") requiredInt32ListField
    2: optional set<string> (go.type = "slice") optionalStringListField
    3: required StringList requiredTypedefStringListField
    4: optional StringList optionalTypedefStringListField
    5: required set<Foo> (go.type = "slice") requiredFooListField
    6: optional set<Foo> (go.type = "slice") optionalFooListField
    7: required FooList requiredTypedefFooListField
    8: optional FooList optionalTypedefFooListField
    9: required set<set<string> (go.type = "slice")> (go.type = "slice") requiredStringListListField
    10: required StringListList requiredTypedefStringListListField
}

const set<string> (go.type = "slice") ConstStringList = ["hello"]
const set<set<string>(go.type = "slice")> (go.type = "slice") ConstListStringList = [["hello"], ["world"]]
