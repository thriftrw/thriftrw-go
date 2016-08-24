/**
 * Package api defines the core interface which plugins use to talk to
 * ThriftRw.
 */

/**
 * Version is the plugin API version specified in this IDL. This MUST be
 * provided in HandshakeResponses to specify which version of the API the
 * plugin implements.
 */
const i32 VERSION = 1

/**
 * ServiceID is an arbitrary unique identifier to reference the diffeerent
 * services in this request.
 */
typedef i32 ServiceID

/**
 * ModuleID is an arbitrary unique identifier to reference the different
 * modules in this request.
 */
typedef i32 ModuleID

/**
 * TypeReference is a reference to a user-defined type.
 */
struct TypeReference {
    1: required string name
    /**
     * Import path for the package defining this type.
     */
    2: required string package
}

/**
 * SimpleType is a standalone native Go type.
 */
enum SimpleType {
    BOOL = 1,     // bool
    BYTE,         // byte
    INT8,         // int8
    INT16,        // int16
    INT32,        // int32
    INT64,        // int64
    FLOAT64,      // float64
    STRING,       // string
    STRUCT_EMPTY, // struct{}
}

/**
 * TypePair is a pair of two types.
 */
struct TypePair {
    1: required Type left
    2: required Type right
}

/**
 * Type is a reference to a Go type which may be native or user defined.
 */
union Type {
    1: SimpleType simpleType
    /**
     * Slice of a type
     *
     * []$sliceType
     */
    2: Type sliceType
    /**
     * Slice of key-value pairs of a pair of types.
     *
     * []struct{Key $left, Value $right}
     */
    3: TypePair keyValueSliceType
    /**
     * Map of a pair of types.
     *
     * map[$left]$right
     */
    4: TypePair mapType
    /**
     * Reference to a user-defined type.
     */
    5: TypeReference referenceType
    /**
     * Pointer to a type.
     */
    6: Type pointerType
}

/**
 * Argument is a single Argument inside a Function.
 * For,
 *
 *      void setValue(1: string key, 2: string value)
 *
 * You get the arguments,
 *
 *      Argument{Name: "Key", Type: Type{SimpleType: SimpleTypeString}}
 *
 *      Argument{Name: "Value", Type: Type{SimpleType: SimpleTypeString}}
 */
struct Argument {
    /**
     * Name of the argument. This is also the name of the argument field
     * inside the args/result struct for that function.
     */
    1: required string name
    /**
     * Argument type.
     */
    2: required Type type
}

/**
 * Function is a single function on a Thrift service.
 */
struct Function {
    /**
     * Name of the Go function.
     */
    1: required string name
    /**
     * Name of the function as defined in the Thrift file.
     */
    2: required string thriftName
    /**
     * List of arguments accepted by the function.
     *
     * This list is in the order specified by the user in the Thrift file.
     */
    3: required list<Argument> arguments
    /**
     * Return type of the function, if any. If this is not set, the function
     * is a void function.
     */
    4: optional Type returnType
    /**
     * List of exceptions raised by the function.
     *
     * This list is in the order specified by the user in the Thrift file.
     */
    5: optional list<Argument> exceptions
}

/**
 * Service is a service defined by the user in the Thrift file.
 */
struct Service {
    /**
     * Name of the service as defined in the Thrift file.
     */
    1: required string name
    /**
     * Import path for the package defining this service.
     */
    2: required string package
    /**
     * Path to the directory containing code for this service.
     *
     * The path is relative to the output directory into which ThriftRW is
     * generating code. Plugins SHOULD NOT make any assumptions about the
     * absolute location of the directory.
     */
    3: required string directory
    /**
     * ID of the parent service.
     */
    4: optional ServiceID parentID
    /**
     * List of functions defined for this service.
     */
    5: required list<Function> functions
    /**
     * ID of the module which declared this service.
     */
    6: required ModuleID moduleID
}

/**
 * Module is a module generated from a single Thrift file. Each module
 * corresponds to exactly one Thrift file and contains all the types and
 * constants defined in that Thrift file.
 */
struct Module {
    /**
     * Import path for the package defining the types for this module.
     */
    1: required string package
    /**
     * Path to the directory containing the code for this module.
     *
     * The path is relative to the output directory into which ThriftRW is
     * generating code. Plugins SHOULD NOT make any assumptions about the
     * absolute location of the directory.
     */
    2: required string directory
}

//////////////////////////////////////////////////////////////////////////////

/**
 * Feature is a functionality offered by a ThriftRW plugin.
 */
enum Feature {
    /**
     * SERVICE_GENERATOR specifies that the plugin may generate arbitrary code
     * for services defined in the Thrift file.
     *
     * If a plugin provides this, it MUST implement the ServiceGenerator
     * service.
     */
    SERVICE_GENERATOR = 1,

    // TODO: TAGGER for struct-tagging plugins
}

/**
 * HandshakeRequest is the initial request sent to the plugin as part of
 * establishing communication and feature negotiation.
 */
struct HandshakeRequest {
}

/**
 * HandshakeResponse is the response from the plugin for a HandshakeRequest.
 */
struct HandshakeResponse {
    /**
     * Name of the plugin. This MUST match the name of the plugin specified
     * over the command line or the program will fail.
     */
    1: required string name
    /**
     * Version of the plugin API.
     *
     * This is NOT the version of the plugin but the version of the plugin
     * APIw hich the plugin respects. This MUST match the version that
     * ThriftRW expects or the program will fail.
     */
    2: required i32 apiVersion
    /**
     * List of features the plugin provides.
     */
    3: required list<Feature> features
}

service Plugin {
    /**
     * handshake performs a handshake with the plugin to negotiate the
     * features provided by it and the version of the plugin API it expects.
     */
    HandshakeResponse handshake(1: HandshakeRequest request)

    /**
     * Informs the plugin process that it will not receive any more requests
     * and it is safe for it to exit.
     */
    void goodbye()
}

//////////////////////////////////////////////////////////////////////////////

/**
 * GenerateServiceRequest is a request to generate code for zero or more
 * Thrift services.
 */
struct GenerateServiceRequest {
    /**
     * IDs of services for which code should be generated.
     *
     * Note that the services map contains information about both, the
     * services being generated and their transitive dependencies. Code should
     * only be generated for service IDs listed here.
     */
    1: required list<ServiceID> rootServices
    /**
     * Map of service ID to service.
     *
     * Any service IDs present in this request will have a corresponding
     * service definition in this map, including services for which code does
     * not need to be generated.
     */
    2: required map<ServiceID, Service> services
    /**
     * Map of module ID to module.
     *
     * Any module IDs present in the request will have a corresponding module
     * definition in this map.
     */
    3: required map<ModuleID, Module> modules
}

/**
 * GenerateServiceResponse is response to a GenerateServiceRequest.
 */
struct GenerateServiceResponse {
    /**
     * Map of file path to file contents.
     *
     * All paths MUST be relative to the output directory into which ThriftRW
     * is generating code. Plugins SHOULD NOT make any assumptions about the
     * absolute location of the directory.
     *
     * The paths MUST NOT contain the string ".." or the request will fail.
     */
    1: optional map<string, binary> files
}

/**
 * ServiceGenerator generates arbitrary code for services.
 *
 * This MUST be implemented if the SERVICE_GENERATOR feature is enabled.
 */
service ServiceGenerator {
    /**
     * Generates code for requested services.
     */
    GenerateServiceResponse generate(1: GenerateServiceRequest request)
}
