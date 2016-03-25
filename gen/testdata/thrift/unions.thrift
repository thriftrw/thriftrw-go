include "./typedefs.thrift"

union Document {
    1: typedefs.PDF pdf
    2: string plainText
}

union ArbitraryValue {
    1: bool boolValue
    2: i64 int64Value
    3: string stringValue
    4: list<ArbitraryValue> listValue
    5: map<string, ArbitraryValue> mapValue
}
