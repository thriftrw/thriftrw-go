include "./enums.thrift"

struct PrimitiveContainers {
    1: optional list<binary> listOfBinary
    2: optional list<i64> listOfInts
    3: optional set<string> setOfStrings
    4: optional set<byte> setOfBytes
    5: optional map<i32, string> mapOfIntToString
    6: optional map<string, bool> mapOfStringToBool
}

struct PrimitiveContainersRequired {
    1: required list<string> listOfStrings
    2: required set<i32> setOfInts
    3: required map<i64, double> mapOfIntsToDoubles
}

struct EnumContainers {
    1: optional list<enums.EnumDefault> listOfEnums
    2: optional set<enums.EnumWithValues> setOfEnums
    3: optional map<enums.EnumWithDuplicateValues, i32> mapOfEnums
}
