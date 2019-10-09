include "./typedefs.thrift"

union EmptyUnion {}

union Document {
    1: typedefs.PDF pdf
    2: string plainText
}

/**
 * ArbitraryValue allows constructing complex values without a schema.
 *
 * A value is one of,
 *
 * * Boolean
 * * Integer
 * * String
 * * A list of other values
 * * A dictionary of other values
 */
union ArbitraryValue {
    1: bool boolValue
    2: i64 int64Value
    3: string stringValue
    4: list<ArbitraryValue> listValue
    5: map<string, ArbitraryValue> mapValue
}
