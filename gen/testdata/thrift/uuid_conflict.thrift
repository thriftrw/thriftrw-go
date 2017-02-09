include "./typedefs.thrift"

typedef string UUID

struct UUIDConflict {
    1: required UUID localUUID
    2: required typedefs.UUID importedUUID
}
