namespace rb v1

struct AddedRequiredField {
    1: optional string A
    2: optional string B
}

struct NewField {}

service Foo {
    void methodB()
}

service Bar {}