union Value {
    1: string s
    2: i64 i
    3: double f
    4: bool b
    5: Value val
}

union NewValue {
    1: string s
    2: i64 i
    3: double f
    4: bool b
    5: NewValue val
    6: byte c
}
