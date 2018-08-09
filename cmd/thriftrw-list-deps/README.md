# list-deps

This tool can be used on a given Thrift file to output the list of Thrift files that the given file depends on.

## Installation

```bash
$ go get go.uber.org/thriftrw/cmd/thriftrw-list-deps
```

## Usage

```bash
$ thriftrw-list-deps --relative=$(pwd) gen/testdata/thrift/structs.thrift
gen/testdata/thrift/enums.thrift
$
```