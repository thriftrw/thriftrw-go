exception EmptyException {}

exception DoesNotExistException {
    1: required string key
    2: optional string Error (go.name="Error2")
}
