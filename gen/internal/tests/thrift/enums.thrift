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

/**
 * Kinds of records stored in the database.
 */
enum RecordType {
  /** Name of the user. */
  NAME,

  /**
   * Home address of the user.
   *
   * This record is always present.
   */
  HOME_ADDRESS,

  /**
   * Home address of the user.
   *
   * This record may not be present.
   */
  WORK_ADDRESS
}

enum lowerCaseEnum {
    containing, lower_case, items
}

// collision with RecordType_Values() function.
enum RecordType_Values { FOO, BAR }
