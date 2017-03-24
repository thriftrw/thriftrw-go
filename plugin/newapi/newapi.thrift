const i32 API_VERSION = 4 // maybe change to i64?

typedef i64 ID

/**
 * Definition represents all definitions for a plugin request.
 */
struct Definition {
  1: optional map<ID, Module> modules
  2: optional map<ID, Const> consts
  3: optional map<ID, Typedef> typedefs
  4: optional map<ID, Enum> enums
  5: optional map<ID, Struct> structs
  6: optional map<ID, Service> services
}

/**
 * Module represents an entire file.
 *
 * The name will include the file path from the common ancestor
 * of the file and all it's includes, for example:
 *
 * a/b/c/foo.thrift: include "../bar.thrift", incluide "../d/baz.thrift"
 *   name: "c/foo.thrift"
 *   includes: "bar.thrift", "d/baz.thrift"
 */
struct Module {
  1: optional ID id
  2: optional string name
  3: optional map<string, string> namespace_language_to_name;
  4: optional set<ID> included_module_ids
  5: optional set<ID> const_ids
  6: optional set<ID> typedef_ids
  7: optional set<ID> enum_ids
  8: optional set<ID> struct_ids
  9: optional set<ID> service_ids
}

enum Primitive {
  BOOL = 1,
  BYTE,
  I16,
  I32,
  I64,
  DOUBLE,
  BINARY,
  STRING,
}

struct TypePair {
  1: optional Type left
  2: optional Type right
}

union Type {
  1: Primitive primitive
  2: Type listType
  3: Type setType
  4: TypePair mapType
  5: ID typedefID
  6: ID enumID
  7: ID structID
}

struct Const {
  1: optional ID id
  2: optional ID module_id
  3: optional Type type
  4: optional string name
  5: optional string value // should we abstract value?
  6: optional map<string, string> annotations
}

struct Typedef {
  1: optional ID id
  2: optional ID module_id
  3: optional Type type
  4: optional string name
  5: optional map<string, string> annotations
}

struct Enum {
  1: optional ID id
  2: optional ID module_id
  3: optional string name
  4: optional map<i32, EnumItem> value_to_enum_item
  5: optional map<string, string> annotations
}

struct EnumItem {
  1: optional i32 value // are we sure this is an i32 per thrift spec?
  2: optional string name
  3: optional map<string, string> annotations
}

struct Struct {
  1: optional ID id
  2: optional ID module_id
  3: optional string name
  4: optional map<i16, Field> tag_to_field
  5: optional map<string, string> annotations
}

struct Field {
  1: optional i16 tag // are we sure this is an i16 per thrift spec?
  2: optional string name
  3: optional Type type
  4: optional bool isRequired
  5: optional map<string, string> annotations
}

struct Service {
  1: optional ID id
  2: optional ID module_id
  3: optional ID extended_service_id
  4: optional string name
  5: optional map<string, Function> name_to_function
  6: optional map<string, string> annotations
}

struct Function {
  1: optional string name
  2: optional map<i16, Argument> tag_to_argument
  3: optional Type response_type
  4: optional map<i16, Argument> tag_to_exception // should this really be an argument?
  5: optional bool isOneway
  6: optional map<string, string> annotations
}

struct Argument {
  1: optional i16 tag // are we sure this is an i16 per thrift spec?
  2: optional string name
  3: optional Type type
  // does this have annotations?
}

enum PluginFeature {
  GENERATE = 1,
  GO_STRUCT_TAGGER
}

struct HandshakeRequest {
}

struct HandshakeResponse {
  1: optional string name
  2: optional i32 apiVersion
  3: optional set<PluginFeature> features
  4: optional string libraryVersion
}

service Plugin {
  HandshakeResponse handshake(1: HandshakeRequest request)
}

struct GenerateRequest {
  1: optional ID module_id // the module id requested to be generated
  2: optional string out_dir_path
  3: optional string thrift_root_dir_path
  4: optional string go_pkg_prefix
  5: optional bool no_recurse // do not generate code for included modules
  6: optional Definition definition
}

struct GenerateResponse {
  1: optional map<string, binary> file_path_to_file // relative to out_dir_path
}

service GenerateService {
  GenerateResponse generate(1: GenerateRequest request)
}

struct GoStructTaggerRequest {
  1: optional ID module_id // the module id requested to be generated
  2: optional Definition definition
}

/**
 * `json:"hello"`
 *    key: "json"
 *    value: "hello"
 */
struct GoStructTaggerResponse {
  // remove the keys from the struct that are usually added by default
  // if there is a key in both add and remove, this is equivalent to modify
  1: optional map<ID, map<i16, set<string>>> struct_id_to_field_tag_to_remove_keys
  // add the key/values to the struct
  2: optional map<ID, map<i16, map<string, string>>> struct_id_to_field_tag_to_add_key_to_value
}

service GoStructTaggerService {
  GoStructTaggerResponse Tag(1: GoStructTaggerRequest request)
}
