exception EmptyException {}

/**
 * Raised when something doesn't exist.
 */
exception DoesNotExistException {
    /** Key that was missing. */
    1: required string key
    2: optional string Error (go.name="Error2")
}
