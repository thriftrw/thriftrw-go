namespace rb v1

struct AddedRequiredField {
    1: optional string A
    2: required string B
    3: required string C
}

service Foo {}