namespace rb v1

struct AddedRequiredField {
    1: optional string A
    2: optional string B
}

service Foo {
    void methodA()
}

service Bar {}