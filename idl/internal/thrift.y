%{
package internal

import "go.uber.org/thriftrw/ast"
%}

%union {
    // Used to record line numbers when the line number at the start point is
    // required.
    line int

    docstring string

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
%type <docstring> docstring
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
    : lineno docstring CONST type IDENTIFIER '=' const_value
        {
            $$ = &ast.Constant{
                Name: $5,
                Type: $4,
                Value: $7,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    /* types */
    | lineno docstring TYPEDEF type IDENTIFIER type_annotations
        {
            $$ = &ast.Typedef{
                Name: $5,
                Type: $4,
                Annotations: $6,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    | lineno docstring ENUM IDENTIFIER '{' enum_items '}' type_annotations
        {
            $$ = &ast.Enum{
                Name: $4,
                Items: $6,
                Annotations: $8,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    | lineno docstring struct_type IDENTIFIER '{' fields '}' type_annotations
        {
            $$ = &ast.Struct{
                Name: $4,
                Type: $3,
                Fields: $6,
                Annotations: $8,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    /* services */
    | lineno docstring SERVICE IDENTIFIER '{' functions '}' type_annotations
        {
            $$ = &ast.Service{
                Name: $4,
                Functions: $6,
                Annotations: $8,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    | lineno docstring SERVICE IDENTIFIER EXTENDS lineno IDENTIFIER '{' functions '}'
      type_annotations
        {
            parent := &ast.ServiceReference{
                Name: $7,
                Line: $6,
            }

            $$ = &ast.Service{
                Name: $4,
                Functions: $9,
                Parent: parent,
                Annotations: $11,
                Line: $1,
                Doc: ParseDocstring($2),
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
    : lineno docstring IDENTIFIER type_annotations
        {
            $$ = &ast.EnumItem{
                Name: $3,
                Annotations: $4,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    | lineno docstring IDENTIFIER '=' INTCONSTANT type_annotations
        {
            value := int($5)
            $$ = &ast.EnumItem{
                Name: $3,
                Value: &value,
                Annotations: $6,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    ;

fields
    : /* nothing */ { $$ = nil }
    | fields field optional_sep { $$ = append($1, $2) }
    ;


field
    : lineno docstring INTCONSTANT ':' field_required type IDENTIFIER type_annotations
        {
            $$ = &ast.Field{
                ID: int($3),
                Name: $7,
                Type: $6,
                Requiredness: $5,
                Annotations: $8,
                Line: $1,
                Doc: ParseDocstring($2),
            }
        }
    | lineno docstring INTCONSTANT ':' field_required type IDENTIFIER '=' const_value
      type_annotations
        {
            $$ = &ast.Field{
                ID: int($3),
                Name: $7,
                Type: $6,
                Requiredness: $5,
                Default: $9,
                Annotations: $10,
                Line: $1,
                Doc: ParseDocstring($2),
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
    : docstring oneway function_type lineno IDENTIFIER '(' fields ')' throws
      type_annotations
        {
            $$ = &ast.Function{
                Name: $5,
                Parameters: $7,
                ReturnType: $<fieldType>3,
                Exceptions: $<fields>9,
                OneWay: $<bul>2,
                Annotations: $10,
                Line: $4,
                Doc: ParseDocstring($1),
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
    : lineno base_type_name type_annotations
        { $$ = ast.BaseType{ID: $2, Annotations: $3, Line: $1} }

    /* container types */
    | lineno MAP '<' type ',' type '>' type_annotations
        { $$ = ast.MapType{KeyType: $4, ValueType: $6, Annotations: $8, Line: $1} }
    | lineno LIST '<' type '>' type_annotations
        { $$ = ast.ListType{ValueType: $4, Annotations: $6, Line: $1} }
    | lineno SET '<' type '>' type_annotations
        { $$ = ast.SetType{ValueType: $4, Annotations: $6, Line: $1} }
    | lineno IDENTIFIER
        { $$ = ast.TypeReference{Name: $2, Line: $1} }
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

    | lineno '[' const_list_items ']' { $$ = ast.ConstantList{Items: $3, Line: $1} }
    | lineno '{' const_map_items  '}' { $$ =  ast.ConstantMap{Items: $3, Line: $1} }
    ;

const_list_items
    : /* nothing */ { $$ = nil }
    | const_list_items const_value optional_sep
        { $$ = append($1, $2) }
    ;

const_map_items
    : /* nothing */ { $$ = nil }
    | const_map_items lineno const_value ':' const_value optional_sep
        { $$ = append($1, ast.ConstantMapItem{Key: $3, Value: $5, Line: $2}) }
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

docstring
    : /* nothing */ { $$ = yylex.(*lexer).LastDocstring() }
    ;

optional_sep
    : ','
    | ';'
    | /* nothing */
    ;

// vim:set ft=yacc:
