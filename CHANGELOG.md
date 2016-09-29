Releases
========

v0.3.0 (unreleased)
------------------

-   **Breaking**: Renamed project to `go.uber.org/thriftrw`.
-   **Breaking**: Keywords reserved by Apache Thrift are now disallowed as
    identifiers in Thrift files.


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
