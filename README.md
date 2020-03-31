# thriftrw-go [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

A Thrift encoding code generator and library for Go.

## Installation

```
go get -u go.uber.org/thriftrw
```

If using [Glide](https://github.com/Masterminds/glide), *at least* `glide version 0.12` is required to install:

```
$ glide --version
glide version 0.12.2

$ glide get 'go.uber.org/thriftrw#^1'
```

## Development Status: Stable

Ready for most users. No breaking changes will be made within the same major
release.

[doc-img]: https://godoc.org/go.uber.org/thriftrw?status.svg
[doc]: https://godoc.org/go.uber.org/thriftrw
[ci-img]: https://travis-ci.com/thriftrw/thriftrw-go.svg?branch=master
[cov-img]: https://codecov.io/gh/thriftrw/thriftrw-go/branch/dev/graph/badge.svg
[ci]: https://travis-ci.com/thriftrw/thriftrw-go
[cov]: https://codecov.io/gh/thriftrw/thriftrw-go

## Development

To install dependencies and build the project run:

```
make build
```

### Testing

To run the entire test suite, run:

```
make test
```

#### Integration tests

A lot of the codebase has its logic exercised via integration tests residing in `gen/internal/tests`.
In other words, since ThriftRW is a code generation library, the way to test that code generation behavior is
implemented correctly is to create a real Thrift struct definition, run code generation, and assert that the output is correct.

Step by step, this process is:

1. Create or use a Thrift struct in `gen/internal/tests/thrift`.  For example, in `gen/internal/tests/thrift/structs.thrift`, you can find
`struct GoTags` that is used to exercise go.tag generation behavior, or `struct NotOmitEmpty` that is used to exercise behavior 
when a field is tagged with `!omitempty`
1. Run `make generate`.  This will generate the go struct from your definition, and place it in `gen/internal/tests/structs/structs.go`
1. Write your test in one of the `*_test.go` files that is pertinent to the feature you are adding.  Oftentimes, these tests 
take in the generated go structs as inputs, and assert on various aspects of code generation like struct tags or json marshaling/unmarshaling behavior.
1. Also remember to add your struct to `gen/quick_test.go` so that your new struct and all of its generic methods (e.g., ToWire, FromWire, String, Equals, etc.) 
can be exercised for code coverage

*Note*: Code coverage is measured across packages, rather than per package.  This is because `go test` is run with the `-coverpkg=./...` flag,
which tells the code coverage tool to measure coverage for this package and all packages in the subdirectories.


