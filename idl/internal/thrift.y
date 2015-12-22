%{
package internal

import "github.com/uber/thriftrw-go/ast"
%}

%union {
    // Used to record line numbers when the line number at the start point is
    // required.
    line int

    // Holds the final AST for the file.
    prog *ast.Program

    // Other intermediate variables:

    bul bool
    str string
    i64 int64
    dub float64

    fieldType ast.Type
    structType ast.StructureType
    baseTypeID ast.BaseTypeID
    fieldRequired ast.Requiredness

    field *ast.Field
    fields []*ast.Field

    header ast.Header
    headers []ast.Header

    function *ast.Function
    functions []*ast.Function

    enumItem *ast.EnumItem
    enumItems []*ast.EnumItem

    definition ast.Definition
    definitions []ast.Definition

    typeAnnotations []*ast.Annotation

    constantValue ast.ConstantValue
    constantValues []ast.ConstantValue
    constantMapItems []ast.ConstantMapItem
}

%token <str> IDENTIFIER
%token <str> LITERAL
%token <i64> INTCONSTANT
%token <dub> DUBCONSTANT

// Reserved keywords
%token NAMESPACE INCLUDE
%token VOID BOOL BYTE I8 I16 I32 I64 DOUBLE STRING BINARY MAP LIST SET
%token ONEWAY TYPEDEF STRUCT UNION EXCEPTION EXTENDS THROWS SERVICE ENUM CONST
%token REQUIRED OPTIONAL TRUE FALSE

%type <line> lineno
%type <prog> program
%type <fieldType> type
%type <baseTypeID> base_type_name
%type <fieldRequired> field_required
%type <structType> struct_type

%type <field> field
%type <fields> fields

%type <header> header
%type <headers> headers

%type <function> function
%type <functions> functions

%type <enumItem> enum_item
%type <enumItems> enum_items

%type <definition> definition
%type <definitions> definitions

%type <constantValue> const_value
%type <constantValues> const_list_items
%type <constantMapItems> const_map_items

%type <typeAnnotations> type_annotation_list type_annotations

%%

program
    : headers definitions
        {
            $$ = &ast.Program{Headers: $1, Definitions: $2}
            yylex.(*lexer).program = $$
            return 0
        }
    ;

/***************************************************************************
 Headers
 ***************************************************************************/

headers
    : /* no headers */     { $$ = nil }
    | headers header     { $$ = append($1, $2) }
    ;

header
    : lineno INCLUDE LITERAL
        {
            $$ = &ast.Include{
                Path: $3,
                Line: $1,
            }
        }
    | lineno INCLUDE IDENTIFIER LITERAL
        {
            $$ = &ast.Include{
                Name: $3,
                Path: $4,
                Line: $1,
            }
        }
    | lineno NAMESPACE '*' IDENTIFIER
        {
            $$ = &ast.Namespace{
                Scope: "*",
                Name: $4,
                Line: $1,
            }
        }
    | lineno NAMESPACE IDENTIFIER IDENTIFIER
        {
            $$ = &ast.Namespace{
                Scope: $3,
                Name: $4,
                Line: $1,
            }
        }
    ;

/***************************************************************************
 Definitions
 ***************************************************************************/

definitions
    : /* nothing */ { $$ = nil }
    | definitions definition optional_sep { $$ = append($1, $2) }
    ;


definition
    /* constants */
    : lineno CONST type IDENTIFIER '=' const_value
        {
            $$ = &ast.Constant{
                Name: $4,
                Type: $3,
                Value: $6,
                Line: $1,
            }
        }
    /* types */
    | lineno TYPEDEF type IDENTIFIER type_annotations
        {
            $$ = &ast.Typedef{
                Name: $4,
                Type: $3,
                Annotations: $5,
                Line: $1,
            }
        }
    | lineno ENUM IDENTIFIER '{' enum_items '}' type_annotations
        {
            $$ = &ast.Enum{
                Name: $3,
                Items: $5,
                Annotations: $7,
                Line: $1,
            }
        }
    | lineno struct_type IDENTIFIER '{' fields '}' type_annotations
        {
            $$ = &ast.Struct{
                Name: $3,
                Type: $2,
                Fields: $5,
                Annotations: $7,
                Line: $1,
            }
        }
    /* services */
    | lineno SERVICE IDENTIFIER '{' functions '}' type_annotations
        {
            $$ = &ast.Service{
                Name: $3,
                Functions: $5,
                Annotations: $7,
                Line: $1,
            }
        }
    | lineno SERVICE IDENTIFIER EXTENDS lineno IDENTIFIER '{' functions '}'
      type_annotations
        {
            parent := &ast.ServiceReference{
                Name: $6,
                Line: $5,
            }

            $$ = &ast.Service{
                Name: $3,
                Functions: $8,
                Parent: parent,
                Annotations: $10,
                Line: $1,
            }
        }
    ;

struct_type
    : STRUCT    { $$ =    ast.StructType }
    | UNION     { $$ =     ast.UnionType }
    | EXCEPTION { $$ = ast.ExceptionType }
    ;

enum_items
    : /* nothing */ { $$ = nil }
    | enum_items enum_item optional_sep { $$ = append($1, $2) }
    ;

enum_item
    : lineno IDENTIFIER type_annotations
        { $$ = &ast.EnumItem{Name: $2, Annotations: $3, Line: $1} }
    | lineno IDENTIFIER '=' INTCONSTANT type_annotations
        {
            value := int($4)
            $$ = &ast.EnumItem{
                Name: $2,
                Value: &value,
                Annotations: $5,
                Line: $1,
            }
        }
    ;

fields
    : /* nothing */ { $$ = nil }
    | fields field optional_sep { $$ = append($1, $2) }
    ;


field
    : lineno INTCONSTANT ':' field_required type IDENTIFIER type_annotations
        {
            $$ = &ast.Field{
                ID: int($2),
                Name: $6,
                Type: $5,
                Requiredness: $4,
                Annotations: $7,
                Line: $1,
            }
        }
    | lineno INTCONSTANT ':' field_required type IDENTIFIER '=' const_value
      type_annotations
        {
            $$ = &ast.Field{
                ID: int($2),
                Name: $6,
                Type: $5,
                Requiredness: $4,
                Default: $8,
                Annotations: $9,
                Line: $1,
            }
        }
    ;

field_required
    : REQUIRED { $$ =    ast.Required }
    | OPTIONAL { $$ =    ast.Optional }
    | /* na */ { $$ = ast.Unspecified }
    ;

functions
    : /* nothing */ { $$ = nil }
    | functions function optional_sep { $$ = append($1, $2) }
    ;

function
    : oneway function_type lineno IDENTIFIER '(' fields ')' throws
      type_annotations
        {
            $$ = &ast.Function{
                Name: $4,
                Parameters: $6,
                ReturnType: $<fieldType>2,
                Exceptions: $<fields>8,
                OneWay: $<bul>1,
                Annotations: $9,
                Line: $3,
            }
        }
    ;

oneway
    : ONEWAY        { $<bul>$ = true }
    | /* nothing */ { $<bul>$ = false }
    ;

function_type
    : VOID { $<fieldType>$ = nil }
    | type { $<fieldType>$ = $1  }
    ;

throws
    : /* nothing */  { $<fields>$ = nil }
    | THROWS '(' fields ')' { $<fields>$ = $3 }
    ;

/***************************************************************************
 Types
 ***************************************************************************/

type
    : base_type_name type_annotations
        { $$ = ast.BaseType{ID: $1, Annotations: $2} }

    /* container types */
    | MAP '<' type ',' type '>' type_annotations
        { $$ = ast.MapType{KeyType: $3, ValueType: $5, Annotations: $7} }
    | LIST '<' type '>' type_annotations
        { $$ = ast.ListType{ValueType: $3, Annotations: $5} }
    | SET '<' type '>' type_annotations
        { $$ = ast.SetType{ValueType: $3, Annotations: $5} }
    | lineno IDENTIFIER
        { $$ = ast.TypeReference{Name: $2, Line: $1 } }
    ;

base_type_name
    : BOOL    { $$ =   ast.BoolTypeID }
    | BYTE    { $$ =     ast.I8TypeID }
    | I8      { $$ =     ast.I8TypeID }
    | I16     { $$ =    ast.I16TypeID }
    | I32     { $$ =    ast.I32TypeID }
    | I64     { $$ =    ast.I64TypeID }
    | DOUBLE  { $$ = ast.DoubleTypeID }
    | STRING  { $$ = ast.StringTypeID }
    | BINARY  { $$ = ast.BinaryTypeID }
    ;

/***************************************************************************
 Constant values
 ***************************************************************************/

const_value
    : INTCONSTANT { $$ = ast.ConstantInteger($1) }
    | DUBCONSTANT { $$ = ast.ConstantDouble($1) }
    | TRUE        { $$ = ast.ConstantBoolean(true) }
    | FALSE       { $$ = ast.ConstantBoolean(false) }
    | LITERAL     { $$ = ast.ConstantString($1) }
    | lineno IDENTIFIER
        { $$ = ast.ConstantReference{Name: $2, Line: $1} }

    | '[' const_list_items ']' { $$ = ast.ConstantList{Items: $2} }
    | '{' const_map_items  '}' { $$ =  ast.ConstantMap{Items: $2} }
    ;

const_list_items
    : /* nothing */ { $$ = nil }
    | const_list_items const_value optional_sep
        { $$ = append($1, $2) }
    ;

const_map_items
    : /* nothing */ { $$ = nil }
    | const_map_items const_value ':' const_value optional_sep
        { $$ = append($1, ast.ConstantMapItem{Key: $2, Value: $4}) }
    ;

/***************************************************************************
 Type annotations
 ***************************************************************************/

type_annotations
    : /* nothing */         { $$ = nil }
    | '(' type_annotation_list ')' { $$ = $2 }
    ;

type_annotation_list
    : /* nothing */ { $$ = nil }
    | type_annotation_list lineno IDENTIFIER '=' LITERAL optional_sep
        { $$ = append($1, &ast.Annotation{Name: $3, Value: $5, Line: $2}) }
    | type_annotation_list lineno IDENTIFIER optional_sep
        { $$ = append($1, &ast.Annotation{Name: $3, Line: $2}) }
    ;

/***************************************************************************
 Other
 ***************************************************************************/

/* Grammar rules that need to record a line number at a specific token should
   include this somewhere. For example,

    foo : bar lineno baz { x := $2 }

  $2 in the above example contains the line number right after 'bar' but before
  'baz'. This way, if 'baz' spans mulitple lines, we still get the line number
  for where the rule started rather than where it ends.
 */
lineno
    : /* nothing */ { $$ = yylex.(*lexer).line }
    ;

optional_sep
    : ','
    | ';'
    | /* nothing */
    ;

// vim:set ft=yacc:
