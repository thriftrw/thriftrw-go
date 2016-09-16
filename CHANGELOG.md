Releases
========

v0.3.0 (unreleased)
------------------

-   **Breaking**: Most keywords reserved by Apache Thrift are now disallowed as
    identifiers in Thrift files.


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
