include "./enums.thrift"
include "./enum_conflict.thrift"
include "./typedefs.thrift"
include "./uuid_conflict.thrift"

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

struct ContainersOfContainers {
    1: optional list<list<i32>> listOfLists;
    2: optional list<set<i32>> listOfSets;
    3: optional list<map<i32, i32>> listOfMaps;

    4: optional set<set<string>> setOfSets;
    5: optional set<list<string>> setOfLists;
    6: optional set<map<string, string>> setOfMaps;

    7: optional map<map<string, i32>, i64> mapOfMapToInt;
    8: optional map<list<i32>, set<i64>> mapOfListToSet;
    9: optional map<set<i32>, list<double>> mapOfSetToListOfDouble;
}

struct MapOfBinaryAndString {
    1: optional map<binary, string> binaryToString;
    2: optional map<string, binary> stringToBinary;
}

struct ListOfConflictingEnums {
    1: required list<enum_conflict.RecordType> records
    2: required list<enums.RecordType> otherRecords
}

struct ListOfConflictingUUIDs {
    1: required list<typedefs.UUID> uuids
    2: required list<uuid_conflict.UUID> otherUUIDs
}
