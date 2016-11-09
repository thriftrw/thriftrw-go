Releases
========

v0.5.0 (unreleased)
-------------------

-   **Breaking**: Generated enums now have first class JSON support. Enums are
    (un)marshalled from/to strings if possible with fallback to integer for
	unrecognized values.
-   Code generation will abort if struct fields, after conversion to Go style
    names, are not unique in the structure.
-   A `go.name` annotation may now be specified to override the names of
    entities in the generated Go code. The annotation is supported for struct,
    union, and exception types, and their fields, enum types and enum items,
    and parameters of functions.


v0.4.0 (2016-11-01)
-------------------

-   **Breaking**: Remove the `--yarpc` flag. Install the ThriftRW YARPC plugin
    from `go.uber.org/yarpc/encoding/thrift/thriftrw-plugin-yarpc` and use
    `--plugin=yarpc` instead.
-   **Breaking**: The `compile` API now exposes annotations made while
    referencing native Thrift types. This changes the `TypeSpec`s for primitive
    types from values to types.
-   The `compile` API now also exposes annotations for `typedef` declarations.
-   Generate args structs and helpers for oneway functions.
-   Expose whether a function is oneway to plugins.
-   Expose the version of the library under `go.uber.org/thriftrw/version.Version`.
-   Generated code will test for version compatibility with the current version
    of ThriftRW during initialization.


v0.3.2 (2016-10-05)
-------------------

-   Fix import paths for code generated using `--yarpc`. Note that this flag
    will be removed in a future version.


v0.3.1 (2016-09-30)
-------------------

-   Fix missing canonical import path to `go.uber.org/thriftrw`.


v0.3.0 (2016-09-29)
-------------------

-   **Breaking**: Renamed project to `go.uber.org/thriftrw`.
-   **Breaking**: Keywords reserved by Apache Thrift are now disallowed as
    identifiers in Thrift files.
-   **Breaking**: The `Package` field of the `plugin.TypeReference`,
    `plugin.Service`, and `plugin.Module` structs was renamed to `ImportPath`.


v0.2.1 (2016-09-27)
-------------------

-   Plugin templates: Fixed a bug where imports in templates would use the base
    name of the package even if it had a hyphen in it if it wasn't available on
    the `GOPATH`.
-   Plugin templates: Imports in generated code are now always qualified if the
    package name doesn't match the base name.


v0.2.0 (2016-09-09)
-------------------

-   Added a `-v`/`--version` flag.
-   Added a plugin system.

    ThriftRW now provides a plugin system to allow customizing code generation.
    Initially, only the generated code for `service` declarations is
    customizable. Check the documentation for more details.
-   **Breaking**: Fixed a bug where all-caps attributes that are not known
    abbreviations were changed to PascalCase.
-   **Breaking**: The `String()` method for `enum` types now returns the name
    of the item as specified in the Thrift file.


v0.1.0 (2016-08-31)
-------------------

-   Initial release.
