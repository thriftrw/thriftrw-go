// Copyright (c) 2024 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/wire"
)

const generateTestIncludedThrift = `
namespace go shared

struct Shared {
  1: required string id
}
`

const generateTestMainThrift = `
namespace go main
namespace java com.example.main

include "included.thrift"

const i32 MaxRetries = 3
const string Greeting = "hello\n\"world\""
const list<i32> Primes = [2, 3, 5]
const map<string, i32> Scores = {"a": 1, "b": 2}

typedef i64 Timestamp
typedef list<string> StringList

enum Color {
  RED = 0,
  GREEN = 1,
  BLUE = 2
}

struct Point {
  1: required i32 x
  2: optional i32 y = 10
  3: optional Color color = Color.RED
  4: optional included.Shared shared
}

union Value {
  1: i32 intValue
  2: string strValue
}

exception NotFound {
  1: required string message
}

service Base {
  void ping()
}

service API extends Base {
  oneway void notify(1: string msg)
  Point getPoint(1: required i32 id) throws (1: NotFound notFound)
}
`

func writeGenerateTestFiles(t *testing.T) (dir, mainPath string) {
	t.Helper()
	dir = t.TempDir()

	mainPath = filepath.Join(dir, "main.thrift")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "included.thrift"), []byte(generateTestIncludedThrift), 0o644))
	require.NoError(t, os.WriteFile(mainPath, []byte(generateTestMainThrift), 0o644))
	return dir, mainPath
}

func typeNames(m *Module) []string {
	names := make([]string, 0, len(m.Types))
	for name := range m.Types {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func constantNames(m *Module) []string {
	names := make([]string, 0, len(m.Constants))
	for name := range m.Constants {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func serviceNames(m *Module) []string {
	names := make([]string, 0, len(m.Services))
	for name := range m.Services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func TestGenerateThriftFile_RoundTrip(t *testing.T) {
	dir, mainPath := writeGenerateTestFiles(t)

	original, err := Compile(mainPath)
	require.NoError(t, err, "compiling source thrift")

	// Generate into the same directory so the relative include still resolves.
	outPath := filepath.Join(dir, "out_main.thrift")
	require.NoError(t, original.GenerateThriftFile(outPath), "generating thrift file")

	generated, err := os.ReadFile(outPath)
	require.NoError(t, err)
	content := string(generated)

	t.Run("content_contains_expected_constructs", func(t *testing.T) {
		for _, want := range []string{
			"namespace go main",
			"namespace java com.example.main",
			`include "included.thrift"`,
			"const i32 MaxRetries = 3",
			`\n\"world\"`, // escaped newline + escaped quote in the string constant
			"typedef i64 Timestamp",
			"typedef list<string> StringList",
			"enum Color {",
			"RED = 0",
			"struct Point {",
			"1: required i32 x",
			"2: optional i32 y = 10",
			"3: optional Color color = Color.RED",
			"included.Shared shared",
			"union Value {",
			"exception NotFound {",
			"service Base {",
			"service API extends Base {",
			"oneway void notify(",
			"throws (",
			"void ping()",
		} {
			assert.Contains(t, content, want)
		}
	})

	t.Run("recompiles_to_equivalent_module", func(t *testing.T) {
		recompiled, err := Compile(outPath)
		require.NoError(t, err, "recompiling generated thrift")

		assert.Equal(t, typeNames(original), typeNames(recompiled), "type names")
		assert.Equal(t, constantNames(original), constantNames(recompiled), "constant names")
		assert.Equal(t, serviceNames(original), serviceNames(recompiled), "service names")
		assert.Equal(t, original.Namespaces, recompiled.Namespaces, "namespaces")
		assert.Equal(t, len(original.Includes), len(recompiled.Includes), "include count")

		// Spot-check a struct's field-level details survive the round trip.
		pt, ok := recompiled.Types["Point"].(*StructSpec)
		require.True(t, ok, "Point should be a struct")
		require.Len(t, pt.Fields, 4)
		assert.True(t, pt.Fields[0].Required, "field x should be required")
		require.NotNil(t, pt.Fields[1].Default, "field y should have a default")

		// The union's fields should not be required.
		union, ok := recompiled.Types["Value"].(*StructSpec)
		require.True(t, ok)
		for _, f := range union.Fields {
			assert.False(t, f.Required, "union field %s should not be required", f.Name)
		}
	})
}

func TestGenerateThriftFile_Annotations(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	thrift := `
namespace go main

typedef set<string> (go.type = "slice") StringList
typedef i64 (js.type = "Long") Timestamp (unit = "ms")

service API {
  oneway void notify(1: string msg = "hello" (go.name = "Message"))
}
`
	require.NoError(t, os.WriteFile(mainPath, []byte(thrift), 0o644))

	original, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, original.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	generated := string(content)

	assert.Contains(t, generated, `typedef set<string> (go.type = "slice") StringList`)
	assert.Contains(t, generated, `typedef i64 (js.type = "Long") Timestamp (unit = "ms")`)
	assert.Contains(t, generated, `oneway void notify(1: string msg = "hello" (go.name = "Message"))`)

	recompiled, err := Compile(outPath)
	require.NoError(t, err, "recompiling generated thrift with annotations")

	stringList, ok := recompiled.Types["StringList"].(*TypedefSpec)
	require.True(t, ok)
	assert.Equal(t, Annotations{"go.type": "slice"}, stringList.Target.ThriftAnnotations())

	timestamp, ok := recompiled.Types["Timestamp"].(*TypedefSpec)
	require.True(t, ok)
	assert.Equal(t, Annotations{"js.type": "Long"}, timestamp.Target.ThriftAnnotations())
	assert.Equal(t, Annotations{"unit": "ms"}, timestamp.Annotations)

	notifyArg := recompiled.Services["API"].Functions["notify"].ArgsSpec[0]
	assert.Equal(t, Annotations{"go.name": "Message"}, notifyArg.Annotations)
	require.NotNil(t, notifyArg.Default)
}

func TestGenerateThriftFile_AnnotatedTypes(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	thrift := `
namespace go main

const list<string (go.name = "x")> ListConst = ["a"]
const list<i32> (deprecated) DeprecatedList = [1]

struct ContainerField {
  1: optional list<i32 (obfuscate)> (c = "d") items
  2: optional map<string, string (a = "b")> (m = "n") mapping
}

typedef list<string (inner = "v")> (outer = "w") NestedTypedef

exception E {
  1: required string message
}

service S {
  list<string (foo = "bar")> search(
    /** query id */
    1: i32 q
  )
  void get() throws (1: E err (http.status = "404"))
}
`
	require.NoError(t, os.WriteFile(mainPath, []byte(thrift), 0o644))

	original, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, original.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	generated := string(content)

	assert.Contains(t, generated, `const list<string (go.name = "x")> ListConst`)
	assert.Contains(t, generated, `const list<i32> (deprecated) DeprecatedList`)
	assert.Contains(t, generated, `1: optional list<i32 (obfuscate)> (c = "d") items`)
	assert.Contains(t, generated, `2: optional map<string, string (a = "b")> (m = "n") mapping`)
	assert.Contains(t, generated, `typedef list<string (inner = "v")> (outer = "w") NestedTypedef`)
	assert.Contains(t, generated, `list<string (foo = "bar")> search(`)
	assert.Contains(t, generated, `query id`)
	assert.Contains(t, generated, `throws (1: E err (http.status = "404"))`)

	recompiled, err := Compile(outPath)
	require.NoError(t, err, "recompiling generated thrift with annotated types")

	listConst := recompiled.Constants["ListConst"].Type.(*ListSpec)
	assert.Equal(t, Annotations{"go.name": "x"}, listConst.ValueSpec.ThriftAnnotations())

	deprecatedList := recompiled.Constants["DeprecatedList"].Type.(*ListSpec)
	assert.Equal(t, Annotations{"deprecated": ""}, deprecatedList.ThriftAnnotations())

	itemsField := recompiled.Types["ContainerField"].(*StructSpec).Fields[0]
	listField := itemsField.Type.(*ListSpec)
	assert.Equal(t, Annotations{"obfuscate": ""}, listField.ValueSpec.ThriftAnnotations())
	assert.Equal(t, Annotations{"c": "d"}, listField.ThriftAnnotations())

	searchArg := recompiled.Services["S"].Functions["search"].ArgsSpec[0]
	assert.Equal(t, "query id", searchArg.Doc)

	throwsExc := recompiled.Services["S"].Functions["get"].ResultSpec.Exceptions[0]
	assert.Equal(t, Annotations{"http.status": "404"}, throwsExc.Annotations)
}

func TestGenerateThriftFile_PrimitiveTypeNames(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	thrift := `
namespace go main

typedef byte ByteAlias
typedef i64 Timestamp

struct Primitives {
  1: required bool flag
  2: required byte raw
  3: required i16 small
  4: required i32 count
  5: required i64 big
  6: required double ratio
  7: required string name
  8: required binary data
  9: required byte (go.name = "Tagged") tagged
}

service S {
  byte echo(1: byte input)
}
`
	require.NoError(t, os.WriteFile(mainPath, []byte(thrift), 0o644))

	original, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, original.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	generated := string(content)

	// I8Spec.ThriftName() is "byte"; generation must not hardcode "i8".
	assert.Contains(t, generated, "typedef byte ByteAlias")
	assert.Contains(t, generated, "2: required byte raw")
	assert.Contains(t, generated, "8: required binary data")
	assert.Contains(t, generated, "9: required byte (go.name = \"Tagged\") tagged")
	assert.Contains(t, generated, "byte echo(1: byte input)")
	assert.NotContains(t, generated, " i8 ")

	recompiled, err := Compile(outPath)
	require.NoError(t, err, "recompiling generated thrift with primitive type names")

	byteField := recompiled.Types["Primitives"].(*StructSpec).Fields[1]
	assert.IsType(t, &I8Spec{}, byteField.Type)

	echoArg := recompiled.Services["S"].Functions["echo"].ArgsSpec[0]
	assert.IsType(t, &I8Spec{}, echoArg.Type)
}

func TestGenerateThriftFile_NilModule(t *testing.T) {
	var m *Module
	assert.Error(t, m.GenerateThriftFile(filepath.Join(t.TempDir(), "out.thrift")))
}

// TestGenerateThriftFile_IncludeRelativeToOutputDir verifies that when the
// generated file is written to a directory different from the source module's
// directory, the emitted include statement is computed relative to the output
// file's location so the generated file recompiles correctly.
func TestGenerateThriftFile_IncludeRelativeToOutputDir(t *testing.T) {
	srcDir, mainPath := writeGenerateTestFiles(t)

	original, err := Compile(mainPath)
	require.NoError(t, err, "compiling source thrift")

	// Write the generated file to a sibling directory. The included.thrift
	// file remains only in srcDir, so the emitted include must traverse
	// upwards (e.g. "../src/included.thrift") for recompilation to succeed.
	outDir := filepath.Join(filepath.Dir(srcDir), "out")
	require.NoError(t, os.MkdirAll(outDir, 0o755))
	outPath := filepath.Join(outDir, "out_main.thrift")

	require.NoError(t, original.GenerateThriftFile(outPath), "generating thrift file")

	generated, err := os.ReadFile(outPath)
	require.NoError(t, err)

	expectedRel, err := filepath.Rel(outDir, filepath.Join(srcDir, "included.thrift"))
	require.NoError(t, err)
	assert.Contains(t, string(generated), `include "`+filepath.ToSlash(expectedRel)+`"`,
		"include should be relative to the output file's directory")

	// The generated file must recompile: this proves the include actually
	// resolves from the output directory, not just that the string looks right.
	recompiled, err := Compile(outPath)
	require.NoError(t, err, "recompiling generated thrift from a different directory")
	assert.Equal(t, len(original.Includes), len(recompiled.Includes), "include count preserved")
}

func TestModuleThriftIDL_Deterministic(t *testing.T) {
	_, mainPath := writeGenerateTestFiles(t)

	m, err := Compile(mainPath)
	require.NoError(t, err)

	first, err := m.thriftIDL(mainPath)
	require.NoError(t, err)
	second, err := m.thriftIDL(mainPath)
	require.NoError(t, err)

	assert.Equal(t, first, second, "generated output should be deterministic across runs")
}

// unknownGenerateType is a TypeSpec outside the set handled by getAnnotatedType.
type unknownGenerateType struct{ nativeThriftType }

func (unknownGenerateType) ThriftName() string                              { return "unknown" }
func (unknownGenerateType) ThriftAnnotations() Annotations                  { return nil }
func (unknownGenerateType) Link(Scope) (TypeSpec, error)                    { return nil, nil }
func (unknownGenerateType) TypeCode() wire.Type                             { return wire.TI32 }
func (unknownGenerateType) ForEachTypeReference(func(TypeSpec) error) error { return nil }

type fakeConstantValue struct{}

func (fakeConstantValue) Link(Scope, TypeSpec) (ConstantValue, error) { return nil, nil }

func TestGenerateThriftFile_ConstantValues(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	thrift := `
namespace go main

const i32 Seed = 7
const i32 Ref = Seed
const bool Flag = false
const double Ratio = 1.5
const set<i32> Values = [1, 2]
const list<i32> EmptyList = []
const map<string, i32> EmptyMap = {}

struct Payload {
  1: required bool ok
  2: required string msg
}
const Payload DefaultPayload = {"ok": false, "msg": "hi"}
`
	require.NoError(t, os.WriteFile(mainPath, []byte(thrift), 0o644))

	m, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, m.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	generated := string(content)

	assert.Contains(t, generated, "const bool Flag = false")
	assert.Contains(t, generated, "const double Ratio = 1.5")
	assert.Contains(t, generated, "const i32 Ref = 7")
	assert.Contains(t, generated, "const set<i32> Values")
	assert.Contains(t, generated, "const list<i32> EmptyList = []")
	assert.Contains(t, generated, "const map<string, i32> EmptyMap = {}")
	assert.Contains(t, generated, "const Payload DefaultPayload")
	assert.Contains(t, generated, `"ok": false`)
	assert.Contains(t, generated, `"msg": "hi"`)
}

func TestGenerateThriftFile_IncludedServiceParent(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "shared.thrift"), []byte(`
namespace go shared

service KeyValue {
  void get()
}
`), 0o644))
	mainPath := filepath.Join(dir, "main.thrift")
	require.NoError(t, os.WriteFile(mainPath, []byte(`
namespace go main

include "shared.thrift"

service BulkKeyValue extends shared.KeyValue {
  void setValues(1: required map<string, binary> items)
}
`), 0o644))

	m, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, m.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	generated := string(content)

	assert.Contains(t, generated, "service BulkKeyValue extends shared.KeyValue")
	assert.Contains(t, generated, "1: required map<string, binary> items")
}

func TestGenerateThriftFile_EnumDetails(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	thrift := `
namespace go main

/**
 * Status codes.
 */
enum Status {
  /** ok */
  OK = 0 (stable),
  FAIL = 1,
} (go.name = "StatusCode")
`
	require.NoError(t, os.WriteFile(mainPath, []byte(thrift), 0o644))

	m, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, m.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	generated := string(content)

	assert.Contains(t, generated, "Status codes")
	assert.Contains(t, generated, "OK = 0 (stable)")
	assert.Contains(t, generated, "} (go.name = \"StatusCode\")")
}

func TestGenerateThriftFile_NoNamespaces(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	require.NoError(t, os.WriteFile(mainPath, []byte(`
struct Empty {}
`), 0o644))

	m, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, m.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "struct Empty")
	assert.NotContains(t, string(content), "namespace")
}

func TestGenerateThriftFile_CreatesNestedOutputDirectory(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	require.NoError(t, os.WriteFile(mainPath, []byte("struct S { 1: required i32 x }"), 0o644))

	m, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "nested", "deep", "out.thrift")
	require.NoError(t, m.GenerateThriftFile(outPath))
	assert.FileExists(t, outPath)
}

func TestGenerateThriftFile_WriteBlockedByFile(t *testing.T) {
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")
	require.NoError(t, os.WriteFile(blocker, []byte("x"), 0o644))

	m := &Module{ThriftPath: filepath.Join(dir, "main.thrift")}
	err := m.GenerateThriftFile(filepath.Join(blocker, "out.thrift"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create directory")
}

func TestGenerateThriftFile_ThriftIDLError(t *testing.T) {
	m := &Module{
		ThriftPath: "/tmp/main.thrift",
		Includes: map[string]*IncludedModule{
			"bad": {Name: "bad", Module: nil},
		},
	}
	err := m.GenerateThriftFile(filepath.Join(t.TempDir(), "out.thrift"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "generate thrift content")
}

func TestWriteIncludes_NilIncludedModule(t *testing.T) {
	var b strings.Builder
	m := &Module{
		ThriftPath: "/tmp/main.thrift",
		Includes: map[string]*IncludedModule{
			"bad": {Name: "bad", Module: nil},
		},
	}
	err := writeIncludes(&b, m, "/tmp/out.thrift")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "included module")
}

func TestWriteConstants_SkipsNilConstant(t *testing.T) {
	var b strings.Builder
	m := &Module{
		ThriftPath: "/tmp/main.thrift",
		Constants: map[string]*Constant{
			"skip": nil,
			"ok": {
				Name:  "ok",
				Type:  &I32Spec{},
				Value: ConstantInt(1),
			},
		},
	}
	require.NoError(t, writeConstants(&b, m))
	assert.Contains(t, b.String(), "const i32 ok = 1")
}

func TestConstantValueToString_AllBranches(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}

	boolStr, err := constantValueToString(ConstantBool(false), module, 0)
	require.NoError(t, err)
	assert.Equal(t, "false", boolStr)

	doubleStr, err := constantValueToString(ConstantDouble(2.5), module, 0)
	require.NoError(t, err)
	assert.Equal(t, "2.5", doubleStr)

	setStr, err := constantValueToString(ConstantSet{ConstantInt(1)}, module, 0)
	require.NoError(t, err)
	assert.Contains(t, setStr, "1")

	enum := &EnumSpec{
		Name: "Color",
		File: "/tmp/main.thrift",
		Items: []EnumItem{
			{Name: "RED", Value: 0},
		},
	}
	enumRef, err := constantValueToString(EnumItemReference{
		Enum: enum,
		Item: &enum.Items[0],
	}, module, 0)
	require.NoError(t, err)
	assert.Equal(t, "Color.RED", enumRef)

	_, err = constantValueToString(ConstReference{}, module, 0)
	require.Error(t, err)

	_, err = constantValueToString(EnumItemReference{}, module, 0)
	require.Error(t, err)

	_, err = constantValueToString(fakeConstantValue{}, module, 0)
	require.Error(t, err)
}

func TestConstantStructToString(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}

	empty, err := constantStructToString(&ConstantStruct{}, module, 0)
	require.NoError(t, err)
	assert.Equal(t, "{}", empty)

	nilStruct, err := constantStructToString(nil, module, 0)
	require.NoError(t, err)
	assert.Equal(t, "{}", nilStruct)

	populated, err := constantStructToString(&ConstantStruct{
		Fields: map[string]ConstantValue{
			"b": ConstantInt(2),
			"a": ConstantInt(1),
		},
	}, module, 0)
	require.NoError(t, err)
	assert.Contains(t, populated, `"a": 1`)
	assert.Contains(t, populated, `"b": 2`)
}

func TestGetAnnotatedType_Errors(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}

	_, err := getAnnotatedType(nil, module)
	require.Error(t, err)

	_, err = getAnnotatedType(unknownGenerateType{}, module)
	require.Error(t, err)

	_, err = getAnnotatedType(typeSpecReference{Name: "Missing"}, module)
	require.Error(t, err)
}

func TestGetQualifiedTypeName_Error(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	_, err := getQualifiedTypeName("Foo", "/other/other.thrift", module)
	require.Error(t, err)
}

func TestGetOriginalServiceParentName(t *testing.T) {
	module := &Module{
		ThriftPath: "/tmp/main.thrift",
		Includes: map[string]*IncludedModule{
			"shared": {
				Name: "shared",
				Module: &Module{
					ThriftPath: "/tmp/shared.thrift",
				},
			},
		},
	}

	localParent := &ServiceSpec{Name: "Base", File: "/tmp/main.thrift"}
	name, err := getOriginalServiceParentName(module, "API", localParent)
	require.NoError(t, err)
	assert.Equal(t, "Base", name)

	includedParent := &ServiceSpec{Name: "KeyValue", File: "/tmp/shared.thrift"}
	name, err = getOriginalServiceParentName(module, "Bulk", includedParent)
	require.NoError(t, err)
	assert.Equal(t, "shared.KeyValue", name)

	_, err = getOriginalServiceParentName(module, "API", nil)
	require.Error(t, err)

	missingParent := &ServiceSpec{Name: "Missing", File: "/tmp/missing.thrift"}
	_, err = getOriginalServiceParentName(module, "API", missingParent)
	require.Error(t, err)
}

func TestWriteTypedefWriteEnumWriteStruct_Errors(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	var b strings.Builder

	require.Error(t, writeTypedef(&b, nil, module))
	require.Error(t, writeEnum(&b, nil))
	require.Error(t, writeStruct(&b, nil, module))

	b.Reset()
	require.Error(t, writeTypedef(&b, &TypedefSpec{
		Name:   "Bad",
		Target: unknownGenerateType{},
	}, module))

	b.Reset()
	require.Error(t, writeStruct(&b, &StructSpec{
		Name: "Bad",
		Type: ast.StructureType(99),
	}, module))
}

func TestWriteFieldNameDefaultAndAnnotations_Error(t *testing.T) {
	var b strings.Builder
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	field := &FieldSpec{
		Name:    "bad",
		Default: ConstReference{},
	}
	err := writeFieldNameDefaultAndAnnotations(&b, field, module)
	require.Error(t, err)
}

func TestWriteFunctionFields_MultilineAndErrors(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	var b strings.Builder

	require.NoError(t, writeFunctionFields(&b, nil, module, "fn", "argument"))

	b.Reset()
	fields := FieldGroup{
		{ID: 1, Name: "a", Type: &I32Spec{}, Doc: "first", Required: true},
		{ID: 2, Name: "b", Type: &StringSpec{}},
	}
	require.NoError(t, writeFunctionFields(&b, fields, module, "fn", "argument"))
	assert.Contains(t, b.String(), "first")
	assert.Contains(t, b.String(), "required")

	b.Reset()
	require.Error(t, writeFunctionFields(&b, FieldGroup{
		{ID: 1, Name: "bad", Type: unknownGenerateType{}},
	}, module, "fn", "argument"))

	b.Reset()
	require.Error(t, writeFunctionFields(&b, FieldGroup{
		{ID: 1, Name: "bad", Default: ConstReference{}},
	}, module, "fn", "exception"))
}

func TestWriteServices_ParentResolutionError(t *testing.T) {
	var b strings.Builder
	module := &Module{
		ThriftPath: "/tmp/main.thrift",
		Services: map[string]*ServiceSpec{
			"Broken": {
				Name: "Broken",
				File: "/tmp/main.thrift",
				Parent: &ServiceSpec{
					Name: "Missing",
					File: "/tmp/missing.thrift",
				},
			},
		},
	}
	err := writeServices(&b, module)
	require.Error(t, err)
}

func TestSortTypesByKind_UnknownTypePriority(t *testing.T) {
	module := &Module{
		Types: map[string]TypeSpec{
			"orphan": &I32Spec{},
		},
	}
	names := sortTypesByKind(module)
	require.Equal(t, []string{"orphan"}, names)
}

func TestConstantMapToString_EmptyAndPopulated(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}

	empty, err := constantMapToString(ConstantMap{}, module, 0)
	require.NoError(t, err)
	assert.Equal(t, "{}", empty)

	populated, err := constantMapToString(ConstantMap{
		{Key: ConstantString("k"), Value: ConstantInt(1)},
	}, module, 0)
	require.NoError(t, err)
	assert.Contains(t, populated, `"k": 1`)
}

func TestConstantSliceToString_Empty(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	out, err := constantSliceToString(nil, module, 0)
	require.NoError(t, err)
	assert.Equal(t, "[]", out)
}

func TestConstReference_QualifiedName(t *testing.T) {
	module := &Module{
		ThriftPath: "/tmp/main.thrift",
		Includes: map[string]*IncludedModule{
			"shared": {
				Name: "shared",
				Module: &Module{
					ThriftPath: "/tmp/shared.thrift",
					Constants: map[string]*Constant{
						"Seed": {Name: "Seed", File: "/tmp/shared.thrift"},
					},
				},
			},
		},
	}
	target := module.Includes["shared"].Module.Constants["Seed"]
	out, err := constantValueToString(ConstReference{Target: target}, module, 0)
	require.NoError(t, err)
	assert.Equal(t, "shared.Seed", out)
}

func TestGenerateThriftFile_WriteFileError(t *testing.T) {
	dir := t.TempDir()
	m := &Module{ThriftPath: filepath.Join(dir, "main.thrift")}
	err := m.GenerateThriftFile(dir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "write thrift file")
}

func TestThriftIDL_ErrorPropagation(t *testing.T) {
	badConstants := &Module{
		ThriftPath: "/tmp/main.thrift",
		Constants: map[string]*Constant{
			"bad": {Name: "bad", Type: unknownGenerateType{}, Value: ConstantInt(1)},
		},
	}
	_, err := badConstants.thriftIDL("/tmp/out.thrift")
	require.Error(t, err)

	badTypes := &Module{
		ThriftPath: "/tmp/main.thrift",
		Types: map[string]TypeSpec{
			"bad": &TypedefSpec{Name: "bad", Target: unknownGenerateType{}},
		},
	}
	_, err = badTypes.thriftIDL("/tmp/out.thrift")
	require.Error(t, err)

	badServices := &Module{
		ThriftPath: "/tmp/main.thrift",
		Services: map[string]*ServiceSpec{
			"Broken": {
				Name:   "Broken",
				File:   "/tmp/main.thrift",
				Parent: &ServiceSpec{Name: "Missing", File: "/tmp/missing.thrift"},
			},
		},
	}
	_, err = badServices.thriftIDL("/tmp/out.thrift")
	require.Error(t, err)
}

func TestWriteConstants_Errors(t *testing.T) {
	var b strings.Builder
	require.Error(t, writeConstants(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Constants: map[string]*Constant{
			"badType": {Name: "badType", Type: unknownGenerateType{}, Value: ConstantInt(1)},
		},
	}))

	b.Reset()
	require.Error(t, writeConstants(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Constants: map[string]*Constant{
			"badValue": {Name: "badValue", Type: &I32Spec{}, Value: fakeConstantValue{}},
		},
	}))

	b.Reset()
	require.NoError(t, writeConstants(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Constants: map[string]*Constant{
			"ok": {Name: "ok", Type: &BoolSpec{}, Value: ConstantBool(true), Doc: "doc"},
		},
	}))
	assert.Contains(t, b.String(), "const bool ok = true")
	assert.Contains(t, b.String(), "doc")
}

func TestWriteTypes_Errors(t *testing.T) {
	var b strings.Builder
	require.Error(t, writeTypes(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Types: map[string]TypeSpec{
			"badTypedef": &TypedefSpec{Name: "badTypedef", Target: unknownGenerateType{}},
		},
	}))

	b.Reset()
	require.Error(t, writeTypes(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Types: map[string]TypeSpec{
			"badStruct": &StructSpec{Name: "badStruct", Type: ast.StructureType(99)},
		},
	}))
}

func TestGetAnnotatedType_ContainerErrors(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}

	_, err := getAnnotatedType(&ListSpec{ValueSpec: unknownGenerateType{}}, module)
	require.Error(t, err)

	_, err = getAnnotatedType(&SetSpec{ValueSpec: unknownGenerateType{}}, module)
	require.Error(t, err)

	_, err = getAnnotatedType(&MapSpec{KeySpec: unknownGenerateType{}, ValueSpec: &I32Spec{}}, module)
	require.Error(t, err)

	_, err = getAnnotatedType(&MapSpec{KeySpec: &StringSpec{}, ValueSpec: unknownGenerateType{}}, module)
	require.Error(t, err)
}

func TestWriteStruct_SuccessAndFieldErrors(t *testing.T) {
	var b strings.Builder
	module := &Module{ThriftPath: "/tmp/main.thrift"}

	require.NoError(t, writeStruct(&b, &StructSpec{
		Name:        "Annotated",
		Type:        ast.StructType,
		Annotations: Annotations{"s": "v"},
		Fields: FieldGroup{
			{ID: 1, Name: "x", Type: &I32Spec{}, Required: true},
			{ID: 2, Name: "y", Type: &StringSpec{}, Default: ConstantString("d")},
		},
	}, module))
	assert.Contains(t, b.String(), "(s = \"v\")")
	assert.Contains(t, b.String(), `= "d"`)

	b.Reset()
	require.Error(t, writeStruct(&b, &StructSpec{
		Name:   "BadFieldType",
		Type:   ast.StructType,
		Fields: FieldGroup{{ID: 1, Name: "x", Type: unknownGenerateType{}}},
	}, module))

	b.Reset()
	require.Error(t, writeStruct(&b, &StructSpec{
		Name:   "BadDefault",
		Type:   ast.StructType,
		Fields: FieldGroup{{ID: 1, Name: "x", Type: &I32Spec{}, Default: ConstReference{}}},
	}, module))
}

func TestWriteFunctionFields_InlineComma(t *testing.T) {
	var b strings.Builder
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	require.NoError(t, writeFunctionFields(&b, FieldGroup{
		{ID: 1, Name: "a", Type: &I32Spec{}},
		{ID: 2, Name: "b", Type: &StringSpec{}},
	}, module, "fn", "argument"))
	assert.Contains(t, b.String(), ", ")
}

func TestWriteServices_SuccessPaths(t *testing.T) {
	var b strings.Builder
	module := &Module{
		ThriftPath: "/tmp/main.thrift",
		Services: map[string]*ServiceSpec{
			"S": {
				Name:        "S",
				File:        "/tmp/main.thrift",
				Annotations: Annotations{"deprecated": ""},
				Functions: map[string]*FunctionSpec{
					"notify": {
						Name:   "notify",
						OneWay: true,
						ArgsSpec: ArgsSpec{
							{ID: 1, Name: "msg", Type: &StringSpec{}},
						},
					},
					"fetch": {
						Name: "fetch",
						ArgsSpec: ArgsSpec{
							{ID: 1, Name: "id", Type: &I32Spec{}, Required: true},
						},
						ResultSpec: &ResultSpec{
							ReturnType: &StringSpec{},
							Exceptions: FieldGroup{
								{ID: 1, Name: "err", Type: &StructSpec{Name: "E", Type: ast.ExceptionType, File: "/tmp/main.thrift"}},
							},
						},
						Annotations: Annotations{"http.method": "GET"},
					},
					"ping": {
						Name:       "ping",
						ResultSpec: &ResultSpec{},
					},
				},
			},
		},
		Types: map[string]TypeSpec{
			"E": &StructSpec{Name: "E", Type: ast.ExceptionType, File: "/tmp/main.thrift"},
		},
	}
	require.NoError(t, writeServices(&b, module))
	out := b.String()
	assert.Contains(t, out, "oneway")
	assert.Contains(t, out, "void ping()")
	assert.Contains(t, out, "string fetch")
	assert.Contains(t, out, "throws")
	assert.Contains(t, out, "(http.method = \"GET\")")
	assert.Contains(t, out, "(deprecated)")

	b.Reset()
	require.Error(t, writeServices(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Services: map[string]*ServiceSpec{
			"BadReturn": {
				Name: "BadReturn",
				File: "/tmp/main.thrift",
				Functions: map[string]*FunctionSpec{
					"bad": {
						Name: "bad",
						ResultSpec: &ResultSpec{
							ReturnType: unknownGenerateType{},
						},
					},
				},
			},
		},
	}))
}

func TestConstantValueToString_ErrorPaths(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}

	_, err := constantValueToString(EnumItemReference{
		Enum: &EnumSpec{Name: "E", File: "/tmp/other.thrift", Items: []EnumItem{{Name: "A", Value: 0}}},
		Item: &EnumItem{Name: "A", Value: 0},
	}, module, 0)
	require.Error(t, err)

	_, err = constantSliceToString([]ConstantValue{ConstReference{}}, module, 0)
	require.Error(t, err)

	_, err = constantMapToString(ConstantMap{
		{Key: ConstantString("k"), Value: ConstReference{}},
	}, module, 0)
	require.Error(t, err)

	_, err = constantStructToString(&ConstantStruct{
		Fields: map[string]ConstantValue{"bad": ConstReference{}},
	}, module, 0)
	require.Error(t, err)
}

func TestGetTypePriority_UnknownStructKind(t *testing.T) {
	priority := getTypePriority(&StructSpec{Name: "X", Type: ast.StructureType(99)})
	assert.Equal(t, 5, priority)
}

func TestWriteServices_WriteFunctionFieldsError(t *testing.T) {
	var b strings.Builder
	err := writeServices(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Services: map[string]*ServiceSpec{
			"S": {
				Name: "S",
				File: "/tmp/main.thrift",
				Functions: map[string]*FunctionSpec{
					"bad": {
						Name: "bad",
						ArgsSpec: ArgsSpec{
							{ID: 1, Name: "x", Type: unknownGenerateType{}},
						},
						ResultSpec: &ResultSpec{},
					},
				},
			},
		},
	})
	require.Error(t, err)

	b.Reset()
	err = writeServices(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Services: map[string]*ServiceSpec{
			"S": {
				Name: "S",
				File: "/tmp/main.thrift",
				Functions: map[string]*FunctionSpec{
					"badArgDefault": {
						Name: "badArgDefault",
						ArgsSpec: ArgsSpec{
							{ID: 1, Name: "x", Type: &I32Spec{}, Default: ConstReference{}},
						},
						ResultSpec: &ResultSpec{},
					},
				},
			},
		},
	})
	require.Error(t, err)
}

func TestWriteServices_ThrowsFieldError(t *testing.T) {
	var b strings.Builder
	err := writeServices(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Services: map[string]*ServiceSpec{
			"S": {
				Name: "S",
				File: "/tmp/main.thrift",
				Functions: map[string]*FunctionSpec{
					"badThrows": {
						Name: "badThrows",
						ResultSpec: &ResultSpec{
							Exceptions: FieldGroup{
								{ID: 1, Name: "e", Type: &StructSpec{Name: "E", Type: ast.ExceptionType, File: "/tmp/main.thrift"}, Default: ConstReference{}},
							},
						},
					},
				},
			},
		},
	})
	require.Error(t, err)
}

func TestWriteTypes_NilEnumError(t *testing.T) {
	var b strings.Builder
	err := writeTypes(&b, &Module{
		ThriftPath: "/tmp/main.thrift",
		Types: map[string]TypeSpec{
			"bad": (*EnumSpec)(nil),
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "enum")
}

func TestGetAnnotatedType_TypedefReference(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	out, err := getAnnotatedType(&TypedefSpec{
		Name:   "Timestamp",
		File:   "/tmp/main.thrift",
		Target: &I64Spec{},
	}, module)
	require.NoError(t, err)
	assert.Equal(t, "Timestamp", out)

	_, err = getAnnotatedType(&TypedefSpec{
		Name:   "Foreign",
		File:   "/tmp/other/other.thrift",
		Target: &I64Spec{},
	}, module)
	require.Error(t, err)
}

func TestConstantMapToString_KeyError(t *testing.T) {
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	_, err := constantMapToString(ConstantMap{
		{Key: ConstReference{}, Value: ConstantInt(1)},
	}, module, 0)
	require.Error(t, err)
}

func TestWriteFunctionFields_DefaultError(t *testing.T) {
	var b strings.Builder
	module := &Module{ThriftPath: "/tmp/main.thrift"}
	err := writeFunctionFields(&b, FieldGroup{
		{ID: 1, Name: "x", Type: &I32Spec{}, Default: ConstReference{}},
	}, module, "fn", "argument")
	require.Error(t, err)
}

func TestWriteIncludes_InvalidPaths(t *testing.T) {
	var b strings.Builder
	m := &Module{
		ThriftPath: "/tmp/main.thrift",
		Includes: map[string]*IncludedModule{
			"shared": {
				Name: "shared",
				Module: &Module{
					ThriftPath: "/tmp/shared.thrift",
				},
			},
		},
	}
	require.NoError(t, writeIncludes(&b, m, "/tmp/out/out.thrift"))

	oldAbs, oldRel := filepathAbs, filepathRel
	t.Cleanup(func() {
		filepathAbs = oldAbs
		filepathRel = oldRel
	})

	b.Reset()
	filepathAbs = func(string) (string, error) {
		return "", fmt.Errorf("abs output dir failed")
	}
	err := writeIncludes(&b, m, "/tmp/out/out.thrift")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "output directory")

	b.Reset()
	filepathAbs = func(path string) (string, error) {
		if path == "/tmp/shared.thrift" {
			return "", fmt.Errorf("abs included failed")
		}
		return filepath.Abs(path)
	}
	err = writeIncludes(&b, m, "/tmp/out/out.thrift")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "included module")

	b.Reset()
	filepathAbs = filepath.Abs
	filepathRel = func(string, string) (string, error) {
		return "", fmt.Errorf("rel failed")
	}
	err = writeIncludes(&b, m, "/tmp/out/out.thrift")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "include path")
}

func TestGenerateThriftFile_TypedefFieldReference(t *testing.T) {
	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.thrift")
	require.NoError(t, os.WriteFile(mainPath, []byte(`
namespace go main

typedef i64 Timestamp

struct Event {
  1: required Timestamp createdAt
}
`), 0o644))

	m, err := Compile(mainPath)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "out.thrift")
	require.NoError(t, m.GenerateThriftFile(outPath))

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "1: required Timestamp createdAt")
}
