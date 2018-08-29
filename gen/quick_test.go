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
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	te "go.uber.org/thriftrw/gen/internal/tests/enums"
	tx "go.uber.org/thriftrw/gen/internal/tests/exceptions"
	tf "go.uber.org/thriftrw/gen/internal/tests/services"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	td "go.uber.org/thriftrw/gen/internal/tests/typedefs"
	tu "go.uber.org/thriftrw/gen/internal/tests/unions"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultValueGenerator(typ reflect.Type) func(*testing.T, *rand.Rand) thriftType {
	return func(t *testing.T, rand *rand.Rand) thriftType {
		for {
			// We will keep trying to generate a value until a valid one
			// is found.

			v, ok := quick.Value(typ, rand)
			require.True(t, ok, "failed to generate a value")

			tval := v.Addr().Interface().(thriftType)

			// TODO(abg): ToWire + EvaluateValue to validate here means we end
			// up serializing this value twice. We may want to include a
			// Validate method on generated types.

			w, err := tval.ToWire()
			if err != nil {
				// Value fails validity check. Try again.
				continue
			}

			// Because we evaluate collections lazily, validation issues
			// with items in them won't be known until we try to serialize
			// it or explicitly evaluate the lazy lists with
			// wire.EvaluateValue.
			if err := wire.EvaluateValue(w); err != nil {
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

func TestQuickRoundTrip(t *testing.T) {
	type testCase struct {
		// Sample value of the type to be tested.
		Sample interface{}

		// Specifies how we generate valid values of this type. Defaults to
		// defaultValueGenerator(Type) if unspecified.
		Generator func(*testing.T, *rand.Rand) thriftType
	}

	// The following types from our tests have been skipped.
	// - unions.ArbitraryValue: Self-reference causes testing/quick to loop
	//   for too long
	// - services.KeyValue_SetValue_Args{}: Accepts an ArbitraryValue
	// - services.KeyValue_SetValueV2_Args: Accepts an ArbitraryValue
	// - services.KeyValue_GetManyValues_Result{}: Produces an ArbitraryValue
	// - services.KeyValue_GetValue_Result{}: Produces an ArbitraryValue

	// TODO(abg): ^Use custom generators to make this not-a-problem.

	tests := []testCase{
		// structs, unions, and exceptions
		{Sample: tc.ContainersOfContainers{}},
		{Sample: tc.EnumContainers{}},
		{Sample: tc.ListOfConflictingEnums{}},
		{Sample: tc.ListOfConflictingUUIDs{}},
		{Sample: tc.MapOfBinaryAndString{}},
		{Sample: tc.PrimitiveContainersRequired{}},
		{Sample: tc.PrimitiveContainers{}},
		{Sample: td.DefaultPrimitiveTypedef{}},
		{Sample: td.Event{}},
		{Sample: td.I128{}},
		{Sample: td.Transition{}},
		{Sample: te.StructWithOptionalEnum{}},
		{Sample: tf.Cache_Clear_Args{}},
		{Sample: tf.Cache_ClearAfter_Args{}},
		{Sample: tf.ConflictingNames_SetValue_Args{}},
		{Sample: tf.ConflictingNames_SetValue_Result{}},
		{Sample: tf.ConflictingNamesSetValueArgs{}},
		{Sample: tf.InternalError{}},
		{Sample: tf.KeyValue_DeleteValue_Args{}},
		{Sample: tf.KeyValue_DeleteValue_Result{}},
		{Sample: tf.KeyValue_GetManyValues_Args{}},
		{Sample: tf.KeyValue_GetValue_Args{}},
		{Sample: tf.KeyValue_SetValue_Result{}},
		{Sample: tf.KeyValue_SetValueV2_Result{}},
		{Sample: tf.KeyValue_Size_Args{}},
		{Sample: tf.KeyValue_Size_Result{}},
		{Sample: tf.NonStandardServiceName_NonStandardFunctionName_Args{}},
		{Sample: tf.NonStandardServiceName_NonStandardFunctionName_Result{}},
		{Sample: ts.ContactInfo{}},
		{Sample: ts.DefaultsStruct{}},
		{Sample: ts.Edge{}},
		{Sample: ts.EmptyStruct{}},
		{Sample: ts.Frame{}},
		{Sample: ts.GoTags{}},
		{Sample: ts.Graph{}},
		{Sample: ts.Node{}},
		{Sample: ts.Omit{}},
		{Sample: ts.Point{}},
		{Sample: ts.PrimitiveOptionalStruct{}},
		{Sample: ts.PrimitiveRequiredStruct{}},
		{Sample: ts.Rename{}},
		{Sample: ts.Size{}},
		{Sample: ts.StructLabels{}},
		{Sample: ts.User{}},
		{Sample: ts.ZapOptOutStruct{}},
		{Sample: tu.Document{}},
		{Sample: tu.EmptyUnion{}},
		{Sample: tx.DoesNotExistException{}},
		{Sample: tx.EmptyException{}},

		// typedefs
		{Sample: td.BinarySet{}},
		{Sample: td.EdgeMap{}},
		{Sample: td.FrameGroup{}},
		{Sample: td.MyEnum(0)},
		{Sample: td.PDF{}},
		{Sample: td.PointMap{}},
		{Sample: td.State("")},
		{Sample: td.StateMap{}},
		{Sample: td.Timestamp(0)},
		{Sample: td.UUID{}},

		// enums
		{
			Sample:    te.EmptyEnum(0),
			Generator: enumValueGenerator(te.EmptyEnum_Values),
		},
		{
			Sample:    te.EnumDefault(0),
			Generator: enumValueGenerator(te.EnumDefault_Values),
		},
		{
			Sample:    te.EnumWithDuplicateName(0),
			Generator: enumValueGenerator(te.EnumWithDuplicateName_Values),
		},
		{
			Sample:    te.EnumWithDuplicateValues(0),
			Generator: enumValueGenerator(te.EnumWithDuplicateValues_Values),
		},
		{
			Sample:    te.EnumWithLabel(0),
			Generator: enumValueGenerator(te.EnumWithLabel_Values),
		},
		{
			Sample:    te.EnumWithValues(0),
			Generator: enumValueGenerator(te.EnumWithValues_Values),
		},
		{
			Sample:    te.LowerCaseEnum(0),
			Generator: enumValueGenerator(te.LowerCaseEnum_Values),
		},
		{
			Sample:    te.RecordType(0),
			Generator: enumValueGenerator(te.RecordType_Values),
		},
		{
			Sample:    te.RecordTypeValues(0),
			Generator: enumValueGenerator(te.RecordTypeValues_Values),
		},
	}

	// Log the seed so that we can reproduce this if it ever fails.
	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))
	t.Logf("Using seed %v for testing/quick", seed)

	const numValues = 1000 // number of values to test against
	for _, tt := range tests {
		typ := reflect.TypeOf(tt.Sample)
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
					w, err := give.ToWire()
					require.NoError(t, err, "failed to Thrift encode %v", give)

					got := reflect.New(typ).Interface().(thriftType)
					require.NoError(t, got.FromWire(w), "failed to Thrift decode from %v", w)

					assert.Equal(t, got, give)
				}
			})
		})
	}
}
