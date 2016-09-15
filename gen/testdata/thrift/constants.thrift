include "./other_constants.thrift"
include "./containers.thrift"
include "./enums.thrift"
include "./exceptions.thrift"
include "./structs.thrift"
include "./unions.thrift"
include "./typedefs.thrift"

const containers.PrimitiveContainers primitiveContainers = {
    "listOfInts": other_constants.listOfInts, // imported constant
    "setOfStrings": ["foo", "bar"],
    "setOfBytes": other_constants.listOfInts, // imported constant with type casting
    "mapOfIntToString": {
        1: "1",
        2: "2",
        3: "3",
    },
    "mapOfStringToBool": {
        "1": 0,
        "2": 1,
        "3": 1,
    }
}

const containers.EnumContainers enumContainers = {
    "listOfEnums": [1, enums.EnumDefault.Foo],
    "setOfEnums": [123, enums.EnumWithValues.Y],
    "mapOfEnums": {
        0: 1,
        enums.EnumWithDuplicateValues.Q: 2,
    },
}

const containers.ContainersOfContainers containersOfContainers = {
    "listOfLists": [[1, 2, 3], [4, 5, 6]],
    "listOfSets": [[1, 2, 3], [4, 5, 6]],
    "listOfMaps": [{1: 2, 3: 4, 5: 6}, {7: 8, 9: 10, 11: 12}],
    "setOfSets": [["1", "2", "3"], ["4", "5", "6"]],
    "setOfLists": [["1", "2", "3"], ["4", "5", "6"]],
    "setOfMaps": [
        {"1": "2", "3": "4", "5": "6"},
        {"7": "8", "9": "10", "11": "12"},
    ],
    "mapOfMapToInt": {
        {"1": 1, "2": 2, "3": 3}: 100,
        {"4": 4, "5": 5, "6": 6}: 200,
    },
    "mapOfListToSet": {
        // more type casting
        other_constants.listOfInts: other_constants.listOfInts,
        [4, 5, 6]: [4, 5, 6],
    },
    "mapOfSetToListOfDouble": {
        [1, 2, 3]: [1.2, 3.4],
        [4, 5, 6]: [5.6, 7.8],
    },
}

const enums.StructWithOptionalEnum structWithOptionalEnum = {
    "e": enums.EnumDefault.Baz
}

const exceptions.EmptyException emptyException = {}

const structs.Graph graph = {
    "edges": [
        {"startPoint": other_constants.some_point, "endPoint": {"x": 3, "y": 4}},
        {"startPoint": {"x": 5, "y": 6}, "endPoint": {"x": 7, "y": 8}},
    ]
}

const structs.Node lastNode = {"value": 3}
const structs.Node node = {
    "value": 1,
    "tail": {"value": 2, "tail": lastNode},
}

const unions.ArbitraryValue arbitraryValue = {
    "listValue": [
        {"boolValue": 1},
        {"int64Value": 2},
        {"stringValue": "hello"},
        {"mapValue": {"foo": {"stringValue": "bar"}}},
    ],
}
// TODO: union validation for constants?

const typedefs.i128 i128 = uuid
const typedefs.UUID uuid = {"high": 1234, "low": 5678}

const typedefs.Timestamp beginningOfTime = 0
const typedefs.FrameGroup frameGroup = [
    {
        "topLeft": {"x": 1, "y": 2},
        "size": {"width": 100, "height": 200},
    }
    {
        "topLeft": {"x": 3, "y": 4},
        "size": {"width": 300, "height": 400},
    },
]

const typedefs.MyEnum myEnum = enums.EnumWithValues.Y

const enums.RecordType NAME = enums.RecordType.NAME
const enums.RecordType HOME = enums.RecordType.HOME_ADDRESS
const enums.RecordType WORK_ADDRESS = enums.RecordType.WORK_ADDRESS

const enums.lowerCaseEnum lower = enums.lowerCaseEnum.items
