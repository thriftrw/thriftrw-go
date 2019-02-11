exception NoErrorException {}

typedef string Key

service NoErrorService {
    // void and no exceptions
    void setValue(1: Key key, 2: string value)

    // Return with exceptions
    Key getValue(1: Key key)
        throws (1: NoErrorException doesNotExist)

    // void with exceptions
    void deleteValue(1: Key key)
        throws (1: NoErrorException doesNotExist)

    i64 size()  // < primitve return value
}
