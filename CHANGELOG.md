Releases
========

v0.2.0 (unreleased)
-------------------

-   Add a plugin system.

    ThriftRW now provides a plugin system to allow customizing code generation.
    Initially, only the generated code for `service` declarations is
    customizable. Check the documentation for more details.

-   Breaking: Fixed a bug where all-caps attributes that are not known
    abbreviations were changed to PascalCase.


v0.1.0 (2016-08-31)
-------------------

-   Initial release.
