include "./enums.thrift"

enum RecordType {
    Name, Email
}

const RecordType defaultRecordType = RecordType.Name

const enums.RecordType defaultOtherRecordType = enums.RecordType.NAME

struct Records {
    1: optional RecordType recordType = defaultRecordType
    2: optional enums.RecordType otherRecordType = defaultOtherRecordType
}
