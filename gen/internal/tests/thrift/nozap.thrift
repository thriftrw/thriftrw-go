enum EnumDefault {
    Foo, Bar, Baz
}

struct PrimitiveRequiredStruct {
    1: required bool boolField
    2: required byte byteField
    3: required i16 int16Field
    4: required i32 int32Field
    5: required i64 int64Field
    6: required double doubleField
    7: required string stringField
    8: required binary binaryField
    9: required list<string> listOfStrings
    10: required set<i32> setOfInts
    11: required map<i64, double> mapOfIntsToDoubles
}

typedef map<string, string> StringMap
typedef PrimitiveRequiredStruct Primitives
typedef list<string> StringList
