enum EmptyEnum {}

enum EnumDefault {
    Foo, Bar, Baz
}

enum EnumWithValues {
    X = 123,
    Y = 456,
    Z = 789,
}

enum EnumWithDuplicateValues {
    P, // 0
    Q = -1,
    R, // 0
}

// enum with item names conflicting with those of another enum
enum EnumWithDuplicateName {
    A, B, C, P, Q, R, X, Y, Z
}

// Enum treated as optional inside a struct
struct StructWithOptionalEnum {
    1: optional EnumDefault e
}

enum RecordType {
  NAME,
  HOME_ADDRESS,
  WORK_ADDRESS
}

enum lowerCaseEnum {
    containing, lower_case, items
}
