# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
 No changes yet.

## [1.33.0] - 2025-07-09
### Changed
- formatType template function takes into account go.type annotation.

## [1.32.0] - 2024-04-23
## Added
- Redacted annotation provides a mechanism to redact certain struct fields from
errors messages and log objects.

## [1.31.0] - 2023-06-09
### Changed
- StreamReader ReadString() and WriteString() performance improvements.

## [1.30.0] - 2023-04-06
### Added
- AddTemplate template option.
- thriftbreak: support for changed types, new files, and optional JSON output.
### Changed
- String() performance improvements for string type definitions.

## [1.29.2] - 2021-09-09
### Fixed
- Streaming now handles unrecongized fields and non-matching field types when
  deserializing Thrift structs.

## [1.29.1] - 2021-08-31
### Fixed
- Streaming encodes now handle `nil` items properly in containers.

## [1.29.0] - 2021-08-30
This release includes support for (de)serializing Thrift structs directly
to/from IO streams without converting them to the intermediate `wire.Value`
representation. This method of (de)serialization is significantly faster and
less memory intensive.

### Added
- `protocol/stream` and `envelope/stream` packages defining the core types
  needed to implement support for streaming serialization.
- `protocol`: `BinaryStreamer` that exports the Binary protocol as a
  `stream.Protocol`.
- `protocol/binary`: The new `Default` variable is the default value of
  `*binary.Protocol` that most users should use. This variable retains the
  struct type so that it can be used for any new interfaces declared in the
  future without another `protocol.Binary*` export.
- All generated types now include `Encode` and `Decode` methods that can
  serialize or deserialize those types using `stream.Writer` and
  `stream.Reader` objects.
- `ast`: All nodes now track the column numbers they're defined on in addition
  to the line numbers.
- `ast`: Add `Annotations(Node)` function that reports the annotations for AST
  nodes that record annotations.

### Changed
- `protocol`:
    - Deprecate `Binary` and `EnvelopeAgnosticBinary` in favor of
      `protocol/binary.Default`.
    - Deprecate `NoEnvelopeResponder`, `EnvelopeV0Responder`, and
      `EnvelopeV1Responder` in favor of versions declared in the
      `protocol/binary` package.

Thanks to [@witriew](https://github.com/witriew), [@dianale31](https://github.com/dianale31), [@usmyth](https://github.com/usmyth), and [@jparise](https://github.com/jparise) for their contributions
to this release.

## [1.28.0] - 2021-07-26
### Added
- `idl.Parse` now returns structured `ParseError` on parse failures.
- `idl.Config` provides a new means of parsing Thrift IDLs. This is the
  recommended API going forward.

### Changed
- ThriftRW now tracks positional information about constant values in the
  Thrift IDL. These were previously ignored. To access this information, see
  the documentation for `idl.Config`.

### Fixed
- Support parsing `struct` fields without field identifiers. This is legacy
  syntax still supported by Apache Thrift. Note that this is a parser-level
  change only; ThriftRW will refuse to generate Go code from these.

Thanks to [@jparise](https://github.com/jparise) for their contributions to
this release.

## [1.27.0] - 2021-05-20
### Added
- ThriftRW is now able to parse Thrift files with `cpp_include` statements.

### Fixed
- `double` constants with exponents but without decimal components are now supported.
- Fix handling of escaped quotes inside string literals.

## [1.26.0] - 2021-02-18
### Changed
- Codegeneration for typedefs now uses use generated `MarshalLog...` functions
  where appropriate to avoid casting to a root go type from go packages of
  transitive thrift dependencies.
- Rewrote internal wire's unsafeStringToBytes to adhere to 1.16 vet check.

## [1.25.1] - 2021-01-04
### Fixed
- Boolean fields with default value `false` now are sent over the wire.

## [1.25.0] - 2020-09-09
### Added
- Add RootModules field to api.GenerateServiceRequest.

## [1.24.0] - 2020-06-18
### Added
- Generate `Default_*` methods that construct Thrift structs with defined
  default values pre-populated.

### Changed
- gen: Redefine Options.Plugin as a struct usable outside go.uber.org/thriftrw.

### Fixed
- Serializing generated Thrift objects no longer mutates them with default
  values for unspecified fields.

## [1.23.0] - 2020-03-31
### Added
- Support opting out of the `omitempty` JSON option by adding `!omitempty` to the
  JSON struct tag.

### Changed
- Drop library dependency in tools.go. This includes: `github.com/golang/mock/mockgen`,
  `golang.org/x/lint`, `golang.org/x/tools`, `honnef.co/go/tools/cmd/staticcheck`.

## [1.22.0] - 2020-01-22
### Added
- Arguments now include Annotations as defined in the Thrift file.

## [1.21.0] - 2020-01-02
### Added
- Generated exceptions now include an `ErrorName()` method that returns the
  name of the exception as defined in the Thrift file.
- gen: Templates now have access to `enumItemName` to determine the Go-level
  name of the enum item.

### Changed
- `nil` slices are now treated as empty lists for fields of `list` type. This
  relaxes the previous requirement of accepting only non-`nil` slices for
  required `list` fields. Note: this does not affect `map` and `set` types
  with unhashable keys slices are accepted.
- Migrated to Go modules.

## [1.20.2] - 2019-10-17
### Fixed
- Added canonical import path directive to avoid checking out ThriftRW at the
  wrong import path.
- Package names are now normalized before Go files are generated.

## [1.20.1] - 2019-07-30
### Fixed
- Fixed field compilation to allow fields with similar looking names and
  different casing.

## [1.20.0] - 2019-06-12
### Changed
- ThriftRW now generates non-plugin code into a single file.
- Module data is now provided to ThriftRW plugins when the Module does not
  contain a service.

## [1.19.1] - 2019-05-16
### Fixed
- Fixed a bug that caused invalid code to be generated if two slices of the
  same type with different `go.type` annotations were encountered in the same
  Thrift file.

## [1.19.0] - 2019-04-26
### Added
- Sets now support a `(go.type = "slice")` annotation to be generated as
  slices rather than maps.

## [1.18.0] - 2019-03-28
### Added
- `Ptr` methods for primititve typedefs.

## [1.17.0] - 2019-03-15
### Changed
- Imports in generated code are now always named imports.

## [1.16.1] - 2019-01-23
### Fixed
- Bump API Version for ThriftRW plugins because the previous release contained
  a significant change to the ThriftRW Plugin API.

## [1.16.0] - 2019-01-22
### Added
- Expose Thrift file names, package prefix, and Thrift root directory to
  plugins.

### Fixed
- plugin: Library version matching was dropped.

## [1.15.0] - 2019-01-14
### Changed
-  Generated`Get*` and `IsSet*` methods on structs are now nil-safe.

## [1.14.0] - 2018-10-18
### Added
- Structs now include `IsSet*` methods for fields that can be nil.

## [1.13.1] - 2018-10-04
### Fixed
- gen/plugin: Fixed a bug where typedefs of structs were mishandled; while they
  should have been pointers, they were generated without `*` and failed to
  compile.
- gen/zap: Fixed a bug where logging nil structs would panic.

## [1.13.0] - 2018-09-10
### Added
- gen: Added support for a `go.label` annotation that allows overriding the
  user-readable string names of enum items and struct fields. This has no
  effect on the names of the generated Go entities.
- Generated types now implement `zapcore.ObjectMarshaler` or
  `zapcore.ArrayMarshaler` where appropriate. This should lead to much faster
  logging of these objects.
- Added `go.nolog` annotation for struct fields: Those with
  this annotation will not be included in Zap logging.
- gen/enum: `MarshalText` and `UnmarshalText` now round-trips, even if
  the enum value is unrecognized.

### Fixed
- ThriftRW now does a bounds-check on field identifiers rather than silently
  truncating them.
- gen: Equals methods on generated structs no longer panic if either value is
  nil.
- gen: Fixed a bug where `*_Values` functions for empty enums would not be
  generated.
- gen: Fixed infinite loop in generated `Equals` methods of specific typedefs.

## [1.12.0] - 2018-06-25
### Added
- gen: Added `ThriftPackageImporter` to control import path
  resolution Thrift files.
- Structs now include getter functions for all fields. This
  improves Apache Thrift compatibility.
- Enums now implement encoding.TextMarshaler.

### Changed
- gen: `NewGenerator` is now usable from other packages.

## [1.11.0] - 2018-03-27
### Added
- Plugins now have access to service and function annotations.

## [1.10.0] - 2018-01-11
### Removed
- Removed version check. Version checks would force code regeneration after
  installing backward-compatible versions of ThriftRW. This change relaxes that
  requirement.

## [1.9.0] - 2017-12-12
### Added
- Adds support for EnvelopeAgnosticProtocol. This upcast for protocol.Binary
  can decode both enveloped or not-enveloped requests and respond in kind, by
  exploiting the non-overlapping grammars of these message types.
- Generated enum types now include a `Ptr()` method.

## [1.8.0] - 2017-09-29
### Added
- Plugins: Annotations declared on user-defined types are now exposed on
  `TypeReference`.

### Changed
- Optional fields of generated structs now enable the `omitempty` JSON option
  if the field holds a list, set, or map.

## [1.7.0] - 2017-09-12
### Added
- AST: Parts of the AST now support parsing docstrings in the `/** ... */`
  style.
- Docstrings on types in the Thrift file are now included in the generated
  code.

## [1.6.0] - 2017-08-28
### Added
- Structs now include getter functions for primitive types.
- Fields of generated struct types are now tagged with easyjson-compatible
  `required` tags.

### Fixed
- Fixed code generation bug for default values of fields with a `typedef` type.

## [1.5.0] - 2017-08-03
### Added
- Code generated by ThriftRW now conforms to
  <https://golang.org/s/generatedcode>.

### Fixed
- Removed gomock and testify from Glide imports to test imports.

## [1.4.0] - 2017-07-21
### Added
- Added support for `go.tag` annotations on struct fields. Corresponding fields
  of the generated Go structs will be tagged with these values.
- AST: Added a `LineNumber` function to get the line on which an AST Node was
  defined.

## [1.3.0] - 2017-07-05
### Added
- Plugins: Added support for overriding the communication channels used by
  plugins.
- Enums now implement the `encoding.TextUnmarshaler` interface. This allows
  retrieving enum values by name and integrates better with other encoding
  formats.
- Enums now have a `<EnumName>_Values()` function which returns all known
  values for that enum.

### Fixed
- Plugins: Template output is now compliant with gofmt.

## [1.2.0] - 2017-04-17
### Added
- The Thrift IDL is now embedded withing the generated package. It is
  accessible via the package global ThriftModule.

## [1.1.0] - 2017-03-23
### Added
- `Equals()` methods are now generated for all custom types.
- Added flags `--no-types`, `--no-constants`, `--no-service-helpers`.
- AST: All type references and complex constant values now record the line
  numbers on which they were specified.
- AST: Added a Node interface to unify different AST object types.
- AST: Added Walk to traverse the AST using a Visitor.

### Fixed
- Handle `nil` values in generated `String()` methods.
- Plugins: Fail code generation if communication with a plugin fails to
  disconnect properly.
- Fixed conflicts in helper functions when imported types had names similar to
  locally defined types.

## [1.0.0] - 2016-11-14
### Changed
- Field names `ToWire`, `FromWire`, `String`, and for exceptions `Error` are
  now reserved. Override the names of these fields in the generated code with
  the `go.name` annotation.
- Plugins: The version of ThriftRW used to compile the plugin is now matched
  against the version actually generating code.

## [0.5.0] - 2016-11-10
### Added
- A `go.name` annotation may now be specified to override the names of entities
  in the generated Go code. The annotation is supported for struct, union, and
  exception types, and their fields, enum types and enum items, and parameters
  of functions.
- Plugins: Added a new `Service.Name` field which contains the name of field
  normalized per Go naming conventions. This, along with `Function.Name` may be
  used to build the names of the `Args`, `Result`, and `Helper` types for a
  function.

### Changed
- **Breaking**: Generated enums now have first-class JSON support. Enums are
  (un)marshalled from/to strings if possible with fallback to integer for
  unrecognized values.
- **Breaking**: `Args`, `Result`, and `Helper` types for service functions are
  now generated in the same package as the user-defined types. These types are
  now named similarly to `$service_$function_Args` where `$service` and
  `$function` are the names of the Thrift service and function normalized based
  on Go naming conventions.
- Code generation will abort if struct fields, after conversion to Go style
  names, are not unique in the structure.
- Plugins: Renamed `Service.Name` to `Service.ThriftName` since it contains the
  name of the service as it appeared in the Thrift file.
- Plugins: Constructors for `Plugin` and `ServiceGenerator` clients and
  handlers are now exposed in the same package as the interfaces.
- Non-primitive types constants are now inlined in the generated Go code
  instead of being referenced in an effort to reduce the impact of user errors
  on the generated code. This is because non-primitive constants were
  previously implemented as global `var`s which might be modified by user code.

### Removed
- Plugins: Removed `Service.Directory` and `Service.ImportPath` because these
  are now same as the corresponding module.

## [0.4.0] - 2016-11-01
### Added
- Expose whether a function is oneway to plugins.
- Expose the version of the library under
  `go.uber.org/thriftrw/version.Version`.
- Generated code will test for version compatibility with the current version
  of ThriftRW during initialization.

### Changed
- **Breaking**: The `compile` API now exposes annotations made while
  referencing native Thrift types. This changes the `TypeSpec`s for primitive
  types from values to types.
- The `compile` API now also exposes annotations for `typedef` declarations.
- Generate args structs and helpers for oneway functions.

### Removed
- **Breaking**: Remove the `--yarpc` flag. Install the ThriftRW YARPC plugin
  from `go.uber.org/yarpc/encoding/thrift/thriftrw-plugin-yarpc` and use
  `--plugin=yarpc` instead.

## [0.3.2] - 2016-10-05
### Fixed
- Fix import paths for code generated using `--yarpc`. Note that this flag will
  be removed in a future version.

## [0.3.1] - 2016-09-30
### Fixed
- Fix missing canonical import path to `go.uber.org/thriftrw`.

## [0.3.0] - 2016-09-29
### Changed
- **Breaking**: Renamed project to `go.uber.org/thriftrw`.
- **Breaking**: Keywords reserved by Apache Thrift are now disallowed as
  identifiers in Thrift files.
- **Breaking**: The `Package` field of the `plugin.TypeReference`,
  `plugin.Service`, and `plugin.Module` structs was renamed to `ImportPath`.

## [0.2.1] - 2016-09-27
### Added
- Plugin templates: Imports in generated code are now always qualified if the
  package name doesn't match the base name.

### Fixed
- Plugin templates: Fixed a bug where imports in templates would use the base
  name of the package even if it had a hyphen in it if it wasn't available on
  the `GOPATH`.

## [0.2.0] - 2016-09-09
### Added
- Added a `-v`/`--version` flag.
- Added a plugin system.

  ThriftRW now provides a plugin system to allow customizing code generation.
  Initially, only the generated code for `service` declarations is
  customizable. Check the documentation for more details.

### Changed
- **Breaking**: The `String()` method for `enum` types now returns the name of
  the item as specified in the Thrift file.

### Fixed
- **Breaking**: Fixed a bug where all-caps attributes that are not known
  abbreviations were changed to PascalCase.

## 0.1.0 - 2016-08-31
### Added
- Initial release.

[Unreleased]: https://github.com/thriftrw/thriftrw-go/compare/v1.33.0...HEAD
[1.33.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.32.0...v1.33.0
[1.32.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.31.0...v1.32.0
[1.31.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.30.0...v1.31.0
[1.30.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.29.2...v1.30.0
[1.29.2]: https://github.com/thriftrw/thriftrw-go/compare/v1.29.1...v1.29.2
[1.29.1]: https://github.com/thriftrw/thriftrw-go/compare/v1.29.0...v1.29.1
[1.29.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.28.0...v1.29.0
[1.28.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.27.0...v1.28.0
[1.27.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.26.0...v1.27.0
[1.26.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.25.1...v1.26.0
[1.25.1]: https://github.com/thriftrw/thriftrw-go/compare/v1.25.0...v1.25.1
[1.25.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.24.0...v1.25.0
[1.24.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.23.0...v1.24.0
[1.23.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.22.0...v1.23.0
[1.22.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.21.0...v1.22.0
[1.21.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.20.2...v1.21.0
[1.20.2]: https://github.com/thriftrw/thriftrw-go/compare/v1.20.1...v1.20.2
[1.20.1]: https://github.com/thriftrw/thriftrw-go/compare/v1.20.0...v1.20.1
[1.20.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.19.1...v1.20.0
[1.19.1]: https://github.com/thriftrw/thriftrw-go/compare/v1.19.0...v1.19.1
[1.19.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.18.0...v1.19.0
[1.18.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.17.0...v1.18.0
[1.17.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.16.1...v1.17.0
[1.16.1]: https://github.com/thriftrw/thriftrw-go/compare/v1.16.0...v1.16.1
[1.16.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.13.1...v1.14.0
[1.13.1]: https://github.com/thriftrw/thriftrw-go/compare/v1.13.0...v.13.1
[1.13.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.10.0...v1.11.0
[1.10.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/thriftrw/thriftrw-go/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/thriftrw/thriftrw-go/compare/v0.5.0...v1.0.0
[0.5.0]: https://github.com/thriftrw/thriftrw-go/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/thriftrw/thriftrw-go/compare/v0.3.2...v0.4.0
[0.3.2]: https://github.com/thriftrw/thriftrw-go/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/thriftrw/thriftrw-go/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/thriftrw/thriftrw-go/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/thriftrw/thriftrw-go/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/thriftrw/thriftrw-go/compare/v0.1.0...v0.2.0
