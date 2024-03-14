exception EmptyException {}

/**
 * Raised when something doesn't exist.
 */
exception DoesNotExistException {
    /** Key that was missing. */
    1: required string key
    2: optional string Error (go.name="Error2")
    3: optional string userName (go.redacted)
}

exception Does_Not_Exist_Exception_Collision {
 /** Key that was missing. */
    1: required string key
    2: optional string Error (go.name="Error2")
} (go.name="DoesNotExistException2")
