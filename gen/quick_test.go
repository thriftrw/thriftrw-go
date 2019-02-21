// Copyright (c) 2015 Uber Technologies, Inc.
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

package gen

import (
	"encoding"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tl "go.uber.org/thriftrw/gen/internal/tests/collision"
	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	tle "go.uber.org/thriftrw/gen/internal/tests/enum_conflict"
	te "go.uber.org/thriftrw/gen/internal/tests/enums"
	tx "go.uber.org/thriftrw/gen/internal/tests/exceptions"
	ter "go.uber.org/thriftrw/gen/internal/tests/noerror"
	tz "go.uber.org/thriftrw/gen/internal/tests/nozap"
	tf "go.uber.org/thriftrw/gen/internal/tests/services"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	td "go.uber.org/thriftrw/gen/internal/tests/typedefs"
	tu "go.uber.org/thriftrw/gen/internal/tests/unions"
	tul "go.uber.org/thriftrw/gen/internal/tests/uuid_conflict"
	envex "go.uber.org/thriftrw/internal/envelope/exception"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/zap/zapcore"
)

func thriftTypeIsValid(v thriftType) bool {
	// TODO(abg): ToWire + EvaluateValue to validate here means we end
	// up serializing this value twice. We may want to include a
	// Validate method on generated types.

	// We validate while serializing.
	w, err := v.ToWire()
	if err != nil {
		return false
	}

	// Because we evaluate collections lazily, validation issues with items in
	// them won't be known until we try to serialize it or explicitly evaluate
	// the lazy lists with wire.EvaluateValue.
	if err := wire.EvaluateValue(w); err != nil {
		return false
	}

	return true
}

func defaultValueGenerator(typ reflect.Type) func(*testing.T, *rand.Rand) thriftType {
	return func(t *testing.T, rand *rand.Rand) thriftType {
		for {
			// We will keep trying to generate a value until a valid one
			// is found.

			v, ok := quick.Value(typ, rand)
			require.True(t, ok, "failed to generate a value")

			tval := v.Addr().Interface().(thriftType)
			if !thriftTypeIsValid(tval) {
				// Value fails validity check. Try again.
				continue
			}

			return tval
		}
	}
}

// Version fo defaultValueGenerator that sets only one of the fields of a
// union.
func unionValueGenerator(sample interface{}) func(*testing.T, *rand.Rand) thriftType {
	typ := reflect.TypeOf(sample)
	return func(t *testing.T, rand *rand.Rand) thriftType {
		for {
			// We will keep trying to generate a value until a valid one
			// is found.

			v := reflect.New(typ)
			if typ.NumField() == 0 {
				return v.Interface().(thriftType)
			}

			field := typ.Field(rand.Intn(typ.NumField()))
			fieldValue, ok := quick.Value(field.Type, rand)
			require.True(t, ok, "failed to generate a value for field %q", field.Name)

			v.Elem().FieldByIndex(field.Index).Set(fieldValue)

			tval := v.Interface().(thriftType)
			if !thriftTypeIsValid(tval) {
				// Value fails validity check. Try again.
				continue
			}

			return tval
		}
	}
}

// enumValueGenerator builds a generator for random enum values given the
// `*_Values` function for that enum.
func enumValueGenerator(valuesFunc interface{}) func(*testing.T, *rand.Rand) thriftType {
	vfunc := reflect.ValueOf(valuesFunc)
	typ := vfunc.Type().Out(0).Elem() // Foo_Values() []Foo -> Foo
	return func(t *testing.T, rand *rand.Rand) thriftType {
		knownValues := vfunc.Call(nil)[0]

		var giveV reflect.Value
		// Flip a coin to decide whether we're evaluating a known or
		// unknown value.
		if rand.Int()%2 == 0 && knownValues.Len() > 0 {
			// Pick a known value at random
			giveV = knownValues.Index(rand.Intn(knownValues.Len()))
		} else {
			// give = MyEnum($randomValue)
			giveV = reflect.New(typ).Elem()
			giveV.Set(reflect.ValueOf(rand.Int31()).Convert(typ))
		}

		return giveV.Addr().Interface().(thriftType)
	}
}

func keyValueSetValueArgsGenerator() func(*testing.T, *rand.Rand) thriftType {
	keyGenerator := defaultValueGenerator(reflect.TypeOf(tf.Key("")))
	valueGenerator := unionValueGenerator(tu.ArbitraryValue{})
	return func(t *testing.T, rand *rand.Rand) thriftType {
		return &tf.KeyValue_SetValue_Args{
			Key:   keyGenerator(t, rand).(*tf.Key),
			Value: valueGenerator(t, rand).(*tu.ArbitraryValue),
		}
	}
}

func keyValueSetValueV2ArgsGenerator() func(*testing.T, *rand.Rand) thriftType {
	keyGenerator := defaultValueGenerator(reflect.TypeOf(tf.Key("")))
	valueGenerator := unionValueGenerator(tu.ArbitraryValue{})
	return func(t *testing.T, rand *rand.Rand) thriftType {
		return &tf.KeyValue_SetValueV2_Args{
			Key:   *keyGenerator(t, rand).(*tf.Key),
			Value: valueGenerator(t, rand).(*tu.ArbitraryValue),
		}
	}
}

func keyValueGetValueResultGenerator() func(*testing.T, *rand.Rand) thriftType {
	successGenerator := unionValueGenerator(tu.ArbitraryValue{})
	doesNotExistGenerator := defaultValueGenerator(reflect.TypeOf(tx.DoesNotExistException{}))
	return func(t *testing.T, rand *rand.Rand) thriftType {
		var result tf.KeyValue_GetValue_Result
		if rand.Int()%2 == 0 {
			result.Success = successGenerator(t, rand).(*tu.ArbitraryValue)
		} else {
			result.DoesNotExist = doesNotExistGenerator(t, rand).(*tx.DoesNotExistException)
		}
		return &result
	}
}

func keyValueGetManyValuesResultGenerator() func(*testing.T, *rand.Rand) thriftType {
	arbitraryValueGenerator := unionValueGenerator(tu.ArbitraryValue{})
	successGenerator := func(t *testing.T, rand *rand.Rand) []*tu.ArbitraryValue {
		values := make([]*tu.ArbitraryValue, rand.Intn(10))
		for i := range values {
			values[i] = arbitraryValueGenerator(t, rand).(*tu.ArbitraryValue)
		}
		return values
	}
	doesNotExistGenerator := defaultValueGenerator(reflect.TypeOf(tx.DoesNotExistException{}))
	return func(t *testing.T, rand *rand.Rand) thriftType {
		var result tf.KeyValue_GetManyValues_Result
		if rand.Int()%2 == 0 {
			result.Success = successGenerator(t, rand)
		} else {
			result.DoesNotExist = doesNotExistGenerator(t, rand).(*tx.DoesNotExistException)
		}
		return &result
	}
}

func noErrorSetValueArgsGenerator() func(*testing.T, *rand.Rand) thriftType {
	keyGenerator := defaultValueGenerator(reflect.TypeOf(ter.Key("")))
	value := "value"
	return func(t *testing.T, rand *rand.Rand) thriftType {
		return &ter.NoErrorService_SetValue_Args{
			Key:   keyGenerator(t, rand).(*ter.Key),
			Value: &value,
		}
	}
}

func noErrorGetValueResultGenerator() func(*testing.T, *rand.Rand) thriftType {
	successGenerator := defaultValueGenerator(reflect.TypeOf(ter.Key("")))
	noErrorExceptionGenerator := defaultValueGenerator(reflect.TypeOf(ter.NoErrorException{}))
	return func(t *testing.T, rand *rand.Rand) thriftType {
		var result ter.NoErrorService_GetValue_Result
		if rand.Int()%2 == 0 {
			result.Success = successGenerator(t, rand).(*ter.Key)
		} else {
			result.DoesNotExist = noErrorExceptionGenerator(t, rand).(*ter.NoErrorException)
		}
		return &result
	}
}

// Returns true for Go types that are nillable from Thrift's point-of-view.
func isThriftNillable(typ reflect.Type) bool {
	// Only struct types and typedefs of nillable Go native types are
	// nillable.
	switch typ.Kind() {
	case reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct:
		return true
	}
	return false
}

type thriftKind int

const (
	thriftStruct thriftKind = iota + 1
	thriftEnum
	thriftTypedef
)

func TestQuickSuite(t *testing.T) {
	type testCase struct {
		// Sample value of the type to be tested.
		Sample interface{}

		// Specifies how we generate valid values of this type. Defaults to
		// defaultValueGenerator(Type) if unspecified.
		Generator func(*testing.T, *rand.Rand) thriftType

		// If set, logging for this type will not be tested. This is needed
		// for typedefs of primitives which can't implement ArrayMarshaler or
		// ObjectMarshaler.
		NoLog bool

		// If set, the Equals check will not be performed. The check should be
		// disabled for types for which we cannot reliably generate random
		// values at this time: maps with unhashable keys. The randomly
		// generated values will have duplicate keys.
		NoEquals bool
		// TODO(abg): Use a custom generator for these types^

		// Kind of Thrift type we're testing. Some tests only run on certain
		// Thrift kinds.
		Kind thriftKind
	}

	tests := []testCase{
		// structs, unions, and exceptions
		{Sample: envex.TApplicationException{}, Kind: thriftStruct},
		{Sample: tc.ContainersOfContainers{}, NoEquals: true, Kind: thriftStruct},
		{Sample: tc.EnumContainers{}, Kind: thriftStruct},
		{Sample: tc.ListOfConflictingEnums{}, Kind: thriftStruct},
		{Sample: tc.ListOfConflictingUUIDs{}, Kind: thriftStruct},
		{Sample: tc.MapOfBinaryAndString{}, NoEquals: true, Kind: thriftStruct},
		{Sample: tc.PrimitiveContainersRequired{}, Kind: thriftStruct},
		{Sample: tc.PrimitiveContainers{}, Kind: thriftStruct},
		{Sample: td.DefaultPrimitiveTypedef{}, Kind: thriftStruct},
		{Sample: td.Event{}, Kind: thriftStruct},
		{Sample: td.I128{}, Kind: thriftStruct},
		{Sample: td.Transition{}, Kind: thriftStruct},
		{Sample: te.StructWithOptionalEnum{}, Kind: thriftStruct},
		{Sample: ter.NoErrorException{}, Kind: thriftStruct},
		{Sample: ter.NoErrorService_DeleteValue_Args{}, Kind: thriftStruct},
		{Sample: ter.NoErrorService_DeleteValue_Result{}, Kind: thriftStruct},
		{Sample: ter.NoErrorService_GetValue_Args{}, Kind: thriftStruct},
		{
			Sample:    ter.NoErrorService_GetValue_Result{},
			Kind:      thriftStruct,
			Generator: noErrorGetValueResultGenerator(),
		},
		{
			Sample:    ter.NoErrorService_SetValue_Args{},
			Kind:      thriftStruct,
			Generator: noErrorSetValueArgsGenerator(),
		},
		{Sample: ter.NoErrorService_SetValue_Result{}, Kind: thriftStruct},
		{Sample: ter.NoErrorService_Size_Args{}, Kind: thriftStruct},
		{Sample: ter.NoErrorService_Size_Result{}, Kind: thriftStruct},
		{Sample: tf.Cache_Clear_Args{}, Kind: thriftStruct},
		{Sample: tf.Cache_ClearAfter_Args{}, Kind: thriftStruct},
		{Sample: tf.ConflictingNames_SetValue_Args{}, Kind: thriftStruct},
		{Sample: tf.ConflictingNames_SetValue_Result{}, Kind: thriftStruct},
		{Sample: tf.ConflictingNamesSetValueArgs{}, Kind: thriftStruct},
		{Sample: tf.InternalError{}, Kind: thriftStruct},
		{Sample: tf.KeyValue_DeleteValue_Args{}, Kind: thriftStruct},
		{Sample: tf.KeyValue_DeleteValue_Result{}, Kind: thriftStruct},
		{Sample: tf.KeyValue_GetManyValues_Args{}, Kind: thriftStruct},
		{
			Sample:    tf.KeyValue_GetManyValues_Result{},
			Kind:      thriftStruct,
			Generator: keyValueGetManyValuesResultGenerator(),
		},
		{Sample: tf.KeyValue_GetValue_Args{}, Kind: thriftStruct},
		{
			Sample:    tf.KeyValue_GetValue_Result{},
			Kind:      thriftStruct,
			Generator: keyValueGetValueResultGenerator(),
		},
		{
			Sample:    tf.KeyValue_SetValue_Args{},
			Kind:      thriftStruct,
			Generator: keyValueSetValueArgsGenerator(),
		},
		{
			Sample:    tf.KeyValue_SetValueV2_Args{},
			Kind:      thriftStruct,
			Generator: keyValueSetValueV2ArgsGenerator(),
		},
		{Sample: tf.KeyValue_SetValue_Result{}, Kind: thriftStruct},
		{Sample: tf.KeyValue_SetValueV2_Result{}, Kind: thriftStruct},
		{Sample: tf.KeyValue_Size_Args{}, Kind: thriftStruct},
		{Sample: tf.KeyValue_Size_Result{}, Kind: thriftStruct},
		{
			Sample: tf.NonStandardServiceName_NonStandardFunctionName_Args{},
			Kind:   thriftStruct,
		},
		{
			Sample: tf.NonStandardServiceName_NonStandardFunctionName_Result{},
			Kind:   thriftStruct,
		},
		{Sample: tl.AccessorConflict{}, Kind: thriftStruct},
		{Sample: tl.AccessorNoConflict{}, Kind: thriftStruct},
		{Sample: tl.PrimitiveContainers{}, Kind: thriftStruct},
		{Sample: tl.StructCollision2{}, Kind: thriftStruct},
		{Sample: tl.StructCollision{}, Kind: thriftStruct},
		{
			Sample:    tl.UnionCollision2{},
			Generator: unionValueGenerator(tl.UnionCollision2{}),
			Kind:      thriftStruct,
		},
		{
			Sample:    tl.UnionCollision{},
			Generator: unionValueGenerator(tl.UnionCollision{}),
			Kind:      thriftStruct,
		},
		{Sample: tl.WithDefault{}, Kind: thriftStruct},
		{Sample: tle.Records{}, Kind: thriftStruct},
		{Sample: ts.ContactInfo{}, Kind: thriftStruct},
		{Sample: ts.DefaultsStruct{}, Kind: thriftStruct},
		{Sample: ts.Edge{}, Kind: thriftStruct},
		{Sample: ts.EmptyStruct{}, Kind: thriftStruct},
		{Sample: ts.Frame{}, Kind: thriftStruct},
		{Sample: ts.GoTags{}, Kind: thriftStruct},
		{Sample: ts.Graph{}, Kind: thriftStruct},
		{Sample: ts.Node{}, Kind: thriftStruct},
		{Sample: ts.Omit{}, Kind: thriftStruct},
		{Sample: ts.Point{}, Kind: thriftStruct},
		{Sample: ts.PersonalInfo{}, Kind: thriftStruct},
		{Sample: ts.PrimitiveOptionalStruct{}, Kind: thriftStruct},
		{Sample: ts.PrimitiveRequiredStruct{}, Kind: thriftStruct},
		{Sample: ts.Rename{}, Kind: thriftStruct},
		{Sample: ts.Size{}, Kind: thriftStruct},
		{Sample: ts.StructLabels{}, Kind: thriftStruct},
		{Sample: ts.User{}, Kind: thriftStruct},
		{Sample: ts.ZapOptOutStruct{}, Kind: thriftStruct},
		{
			Sample:    tu.ArbitraryValue{},
			Generator: unionValueGenerator(tu.ArbitraryValue{}),
			Kind:      thriftStruct,
		},
		{
			Sample:    tu.Document{},
			Generator: unionValueGenerator(tu.Document{}),
			Kind:      thriftStruct,
		},
		{
			Sample:    tu.EmptyUnion{},
			Generator: unionValueGenerator(tu.EmptyUnion{}),
			Kind:      thriftStruct,
		},
		{Sample: tul.UUIDConflict{}, Kind: thriftStruct},
		{Sample: tx.DoesNotExistException{}, Kind: thriftStruct},
		{Sample: tx.EmptyException{}, Kind: thriftStruct},
		{
			Sample: tz.PrimitiveRequiredStruct{},
			NoLog:  true,
			Kind:   thriftStruct,
		},

		// typedefs
		{Sample: td.BinarySet{}, Kind: thriftTypedef},
		{Sample: td.EdgeMap{}, Kind: thriftTypedef},
		{Sample: td.FrameGroup{}, Kind: thriftTypedef},
		{Sample: td.MyEnum(0), Kind: thriftTypedef},
		{Sample: td.PDF{}, NoLog: true, Kind: thriftTypedef},
		{Sample: td.PointMap{}, Kind: thriftTypedef},
		{Sample: td.State(""), NoLog: true, Kind: thriftTypedef},
		{Sample: td.StateMap{}, Kind: thriftTypedef},
		{Sample: td.Timestamp(0), NoLog: true, Kind: thriftTypedef},
		{Sample: td.UUID{}, Kind: thriftTypedef},
		{Sample: tl.LittlePotatoe(0), NoLog: true, Kind: thriftTypedef},
		{Sample: tl.LittlePotatoe2(0.0), NoLog: true, Kind: thriftTypedef},
		{Sample: tul.UUID(""), NoLog: true, Kind: thriftTypedef},
		{Sample: tz.StringMap{}, NoLog: true, Kind: thriftTypedef},
		{Sample: tz.Primitives{}, NoLog: true, Kind: thriftTypedef},
		{Sample: tz.StringList{}, NoLog: true, Kind: thriftTypedef},

		// enums
		{
			Sample:    envex.ExceptionType(0),
			Generator: enumValueGenerator(envex.ExceptionType_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.EmptyEnum(0),
			Generator: enumValueGenerator(te.EmptyEnum_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.EnumDefault(0),
			Generator: enumValueGenerator(te.EnumDefault_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.EnumWithDuplicateName(0),
			Generator: enumValueGenerator(te.EnumWithDuplicateName_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.EnumWithDuplicateValues(0),
			Generator: enumValueGenerator(te.EnumWithDuplicateValues_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.EnumWithLabel(0),
			Generator: enumValueGenerator(te.EnumWithLabel_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.EnumWithValues(0),
			Generator: enumValueGenerator(te.EnumWithValues_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.LowerCaseEnum(0),
			Generator: enumValueGenerator(te.LowerCaseEnum_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.RecordType(0),
			Generator: enumValueGenerator(te.RecordType_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    te.RecordTypeValues(0),
			Generator: enumValueGenerator(te.RecordTypeValues_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    tl.MyEnum(0),
			Generator: enumValueGenerator(tl.MyEnum_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    tl.MyEnum2(0),
			Generator: enumValueGenerator(tl.MyEnum2_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    tle.RecordType(0),
			Generator: enumValueGenerator(tle.RecordType_Values),
			Kind:      thriftEnum,
		},
		{
			Sample:    tz.EnumDefault(0),
			Generator: enumValueGenerator(tz.EnumDefault_Values),
			Kind:      thriftEnum,
			NoLog:     true,
		},
	}

	// Log the seed so that we can reproduce this if it ever fails.
	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))
	t.Logf("Using seed %v for testing/quick", seed)

	const numValues = 1000 // number of values to test against
	for _, tt := range tests {
		typ := reflect.TypeOf(tt.Sample)
		suite := quickSuite{Type: typ}

		t.Run(typ.Name(), func(t *testing.T) {
			generator := tt.Generator
			if generator == nil {
				generator = defaultValueGenerator(typ)
			}

			values := make([]thriftType, numValues)
			for i := range values {
				values[i] = generator(t, rand)
			}

			t.Run("Thrift", func(t *testing.T) {
				for _, give := range values {
					suite.testThriftRoundTrip(t, give)
				}
			})

			t.Run("String", func(t *testing.T) {
				for _, give := range values {
					suite.testString(t, give)
				}
			})

			if isThriftNillable(typ) {
				t.Run("StringNil", suite.testStringNil)
			}

			switch tt.Kind {
			case thriftEnum:
				t.Run("JSON", func(t *testing.T) {
					for _, give := range values {
						suite.testJSONRoundTrip(t, give)
					}
				})

				t.Run("Text", func(t *testing.T) {
					for _, give := range values {
						suite.testTextRoundtrip(t, give)
					}
				})

				t.Run("Ptr", func(t *testing.T) {
					for _, give := range values {
						suite.testEnumPtr(t, give)
					}
				})

			case thriftStruct:
				t.Run("Accessors/Get", func(t *testing.T) {
					for _, give := range values {
						suite.testGetAccessors(t, give)
					}
				})

				t.Run("Accessors/GetOnNil", func(t *testing.T) {
					suite.testGetAccessorsOnNil(t)
				})

				t.Run("Accessors/IsSet", func(t *testing.T) {
					for _, give := range values {
						suite.testIsSetAccessors(t, give)
					}
				})

				t.Run("Accessors/IsSetOnNil", func(t *testing.T) {
					suite.testIsSetAccessorsOnNil(t)
				})
			}

			if !tt.NoLog {
				t.Run("Zap", func(t *testing.T) {
					for _, give := range values {
						suite.testLogging(t, give)
					}
				})

				if isThriftNillable(typ) {
					t.Run("Zap/Nil", suite.testLoggingNil)
				}
			}

			if !tt.NoEquals {
				t.Run("Equals", func(t *testing.T) {
					for _, give := range values {
						suite.testEquals(t, give)
					}
				})

				if typ.Kind() == reflect.Struct {
					t.Run("EqualsNil", suite.testEqualsNil)
				}
			}

		})
	}
}

type quickSuite struct {
	Type reflect.Type
}

// Builds a new empty value of the underlying type.
//
// For structs, returns a pointer to an empty struct. For containers, an empty
// container (not nil) is returned. For everything else, zero value is
// returned.
func (q *quickSuite) newEmpty() reflect.Value {
	v := reflect.New(q.Type).Elem()
	switch q.Type.Kind() {
	case reflect.Map:
		v.Set(reflect.MakeMap(q.Type))
		return v
	case reflect.Slice:
		v.Set(reflect.MakeSlice(q.Type, 0, 0))
		return v
	case reflect.Struct:
		v.Set(reflect.Zero(q.Type))
		return v.Addr()
	default:
		return v
	}
}

// Same as newEmpty but a pointer to the value is returned.
func (q *quickSuite) newEmptyPtr() reflect.Value {
	if q.Type.Kind() == reflect.Struct {
		return q.newEmpty()
	}
	return q.newEmpty().Addr()
}

// Builds a new nil value of the underlying type.
//
// For structs, a nil pointer to the struct is returned. For containers, a nil
// pointer to the container is returned. Using this with non-nillable types is
// invalid.
func (q *quickSuite) newNil(t *testing.T) (v reflect.Value) {
	defer func() {
		require.True(t, v.IsNil(), "bug: newNil generated non-nil value") // sanity check
	}()

	switch q.Type.Kind() {
	case reflect.Struct:
		return reflect.Zero(reflect.PtrTo(q.Type))
	default:
		return reflect.Zero(q.Type)
	}
}

// Tests that the provided value round-trips successfully with wire.Value.
func (q *quickSuite) testThriftRoundTrip(t *testing.T, give thriftType) {
	w, err := give.ToWire()
	require.NoError(t, err, "failed to Thrift encode %v", give)

	got := q.newEmptyPtr().Interface().(thriftType)
	require.NoError(t, got.FromWire(w), "failed to Thrift decode from %v", w)

	assert.Equal(t, got, give)
}

// Tests that String() works on any valid value of this type.
func (q *quickSuite) testString(t *testing.T, give thriftType) {
	assert.NotPanics(t, func() {
		_ = give.String()
	}, "failed to String %#v", give)
}

// Tests that String does not panic with a nil value of this type.
func (q *quickSuite) testStringNil(t *testing.T) {
	v := q.newNil(t).Interface().(fmt.Stringer)
	assert.NotPanics(t, func() {
		_ = v.String()
	})
}

// For types that support it (enums only at this time), tests that JSON
// representations round-trip successfully.
func (q *quickSuite) testJSONRoundTrip(t *testing.T, giveVal thriftType) {
	give, ok := giveVal.(json.Marshaler)
	require.True(t, ok, "does not implement json.Marshaler")

	bs, err := give.MarshalJSON()
	require.NoError(t, err, "failed to encode %v", give)

	got, ok := q.newEmptyPtr().Interface().(json.Unmarshaler)
	require.True(t, ok, "does not implement json.Unmarshaler")

	require.NoError(t, got.UnmarshalJSON(bs), "failed to decode from %q", bs)
	assert.Equal(t, got, give, "could not round-trip")
}

// For types that support it (enums only at this time), tests that
// encoding.TextMarshaler representations round-trip successfully.
func (q *quickSuite) testTextRoundtrip(t *testing.T, giveVal thriftType) {
	give, ok := giveVal.(encoding.TextMarshaler)
	require.True(t, ok, "does not implement encoding.TextMarshaler")

	bs, err := give.MarshalText()
	require.NoError(t, err, "failed to encode %v", give)

	got, ok := q.newEmptyPtr().Interface().(encoding.TextUnmarshaler)
	require.True(t, ok, "does not implement encoding.TextUnmarshaler")

	require.NoError(t, got.UnmarshalText(bs), "failed to decode from %q", bs)
	assert.Equal(t, got, give, "could not round-trip")
}

// Tests that the object can be logged by Zap.
func (q *quickSuite) testLogging(t *testing.T, give thriftType) {
	enc := zapcore.NewMapObjectEncoder()

	if obj, ok := give.(zapcore.ObjectMarshaler); ok {
		assert.NoErrorf(t, obj.MarshalLogObject(enc), "failed to log %v", give)
		return
	}

	if arr, ok := give.(zapcore.ArrayMarshaler); ok {
		assert.NoErrorf(t, enc.AddArray("values", arr), "failed to log %v", give)
		return
	}

	t.Fatal(
		"Type does not implement zapcore.ObjectMarshaler or zapcore.ArrayMarshaler. "+
			"Did you mean to add NoLog?", q.Type)
}

func (q *quickSuite) testLoggingNil(t *testing.T) {
	enc := zapcore.NewMapObjectEncoder()

	v := q.newNil(t).Interface()

	if obj, ok := v.(zapcore.ObjectMarshaler); ok {
		assert.NoErrorf(t, obj.MarshalLogObject(enc), "failed to log %v", v)
		return
	}

	if arr, ok := v.(zapcore.ArrayMarshaler); ok {
		assert.NoErrorf(t, enc.AddArray("values", arr), "failed to log %v", v)
		return
	}

	t.Fatal(
		"Type does not implement zapcore.ObjectMarshaler or zapcore.ArrayMarshaler. "+
			"Did you mean to add NoLog?", q.Type)
}

// Tests that the v.Equals(v) always returns true.
func (q *quickSuite) testEquals(t *testing.T, giveVal thriftType) {
	give := reflect.ValueOf(giveVal)
	rhs := give

	equals := give.MethodByName("Equals")
	require.True(t, equals.IsValid(), "Type does not implement Equals()")

	if equals.Type().In(0) != rhs.Type() {
		// We were passing the objects around by pointer but
		// we need the value-form here.
		rhs = rhs.Elem()
	}

	assert.True(t,
		equals.Call([]reflect.Value{rhs})[0].Bool(),
		"%v should be equal to itself", giveVal)
}

// Tests that Equals methods work with nil values and receivers.
func (q *quickSuite) testEqualsNil(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		// var x, y *Type
		x := q.newNil(t)
		y := q.newNil(t)

		// x.Equals(y)
		result := x.MethodByName("Equals").
			Call([]reflect.Value{y})[0].
			Interface().(bool)

		assert.True(t, result)
	})

	t.Run("lhs not nil", func(t *testing.T) {
		// var x Type
		// var y *Type
		x := q.newEmpty()
		y := q.newNil(t)

		// x.Equals(y)
		result := x.MethodByName("Equals").
			Call([]reflect.Value{y})[0].
			Interface().(bool)

		assert.False(t, result)
	})

	t.Run("rhs not nil", func(t *testing.T) {
		// var x *Type
		// var y Type
		x := q.newNil(t)
		y := q.newEmpty()

		// x.Equals(y)
		result := x.MethodByName("Equals").
			Call([]reflect.Value{y})[0].
			Interface().(bool)

		assert.False(t, result)
	})
}

// Tests that Ptr methods on enums return the same value back.
func (q *quickSuite) testEnumPtr(t *testing.T, give thriftType) {
	// TODO(abg): should we generate Ptr and _Values for typedefs of enums?
	v := reflect.ValueOf(give)
	ptr := v.MethodByName("Ptr").Call(nil)[0]
	require.Equal(t, reflect.Ptr, ptr.Kind(), "must be a pointer")
	assert.Equal(t, give, ptr.Interface(),
		"pointer must point back to original value")
}

// Tests that each field of a struct has an accessor that returns the same
// value as the field.
func (q *quickSuite) testGetAccessors(t *testing.T, giveVal thriftType) {
	// TODO(abg): should we generate accessors for typedefs of structs?
	give := reflect.ValueOf(giveVal)
	for i := 0; i < q.Type.NumField(); i++ {
		field := q.Type.Field(i)
		fieldValue := give.Elem().FieldByIndex(field.Index)
		accessorValue := give.MethodByName("Get" + field.Name).Call(nil)[0]

		// For optional primitive fields, we use pointers but the accessors
		// return the derefenced value (defaulting to the zero-value of that
		// type). So we'll do the same before comparing results.
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() != reflect.Struct {
			if fieldValue.IsNil() {
				fieldValue = reflect.Zero(field.Type.Elem())
			} else {
				fieldValue = fieldValue.Elem()
			}
		}

		assert.Equal(t, fieldValue.Interface(), accessorValue.Interface())
	}
}

func (q *quickSuite) testIsSetAccessors(t *testing.T, giveVal thriftType) {
	give := reflect.ValueOf(giveVal)
	for i := 0; i < q.Type.NumField(); i++ {
		field := q.Type.Field(i)
		if !isThriftNillable(field.Type) {
			// The field isn't nillable.
			continue
		}

		isSetMethod := give.MethodByName("IsSet" + field.Name)
		if !assert.Truef(t, isSetMethod.IsValid(), "must have an IsSet%v method", field.Name) {
			continue
		}

		fieldValue := give.Elem().FieldByIndex(field.Index)
		isSet, _ := isSetMethod.Call(nil)[0].Interface().(bool)
		assert.Equal(t, !fieldValue.IsNil(), isSet)
	}
}

// Tests accessors on nil structs.
func (q *quickSuite) testGetAccessorsOnNil(t *testing.T) {
	give := q.newNil(t)
	for i := 0; i < q.Type.NumField(); i++ {
		field := q.Type.Field(i)

		assert.NotPanics(t, func() {
			give.MethodByName("Get" + field.Name).Call(nil)
		})
	}
}

func (q *quickSuite) testIsSetAccessorsOnNil(t *testing.T) {
	give := q.newNil(t)
	for i := 0; i < q.Type.NumField(); i++ {
		field := q.Type.Field(i)
		if !isThriftNillable(field.Type) {
			// The field isn't nillable.
			continue
		}

		require.Falsef(t, give.MethodByName("IsSet" + field.Name).Call(nil)[0].Bool(),
			"field %q must be unset", field.Name)
	}
}
