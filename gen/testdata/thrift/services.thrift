include "./unions.thrift"
include "./exceptions.thrift"

typedef string Key

exception InternalError {
    1: optional string message
}

service KeyValue {
    // void and no exceptions
    void setValue(1: Key key, 2: unions.ArbitraryValue value)

    void setValueV2(
        1: required Key key,
        2: required unions.ArbitraryValue value,
    )

    // Return with exceptions
    unions.ArbitraryValue getValue(1: Key key)
        throws (1: exceptions.DoesNotExistException doesNotExist)

    // void with exceptions
    void deleteValue(1: Key key)
        throws (
            1: exceptions.DoesNotExistException doesNotExist,
            2: InternalError internalError
        )

    list<unions.ArbitraryValue> getManyValues(
        1: list<Key> range  // < reserved keyword as an argument
    ) throws (
        1: exceptions.DoesNotExistException doesNotExist,
    )

    i64 size()  // < primitve return value
}

service Cache {
    oneway void clear()
    oneway void clearAfter(1: i64 durationMS)
}

struct ConflictingNames_SetValue_Args {
    1: required string key
    2: required binary value
}

service ConflictingNames {
    void setValue(1: ConflictingNames_SetValue_Args request)
}

service non_standard_service_name {
    void non_standard_function_name()
}
