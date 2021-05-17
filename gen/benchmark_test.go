package gen_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/wire"
)

type thriftType interface {
	ToWire() (wire.Value, error)
	FromWire(wire.Value) error
}

func BenchmarkRoundTrip(b *testing.B) {
	type benchCase struct {
		name string
		give thriftType
	}

	benchmarks := []benchCase{
		{
			name: "PrimitiveOptionalStruct",
			give: &ts.PrimitiveOptionalStruct{
				BoolField:   ptr.Bool(true),
				ByteField:   ptr.Int8(42),
				Int16Field:  ptr.Int16(123),
				Int32Field:  ptr.Int32(1234),
				Int64Field:  ptr.Int64(123456),
				DoubleField: ptr.Float64(math.Pi),
				StringField: ptr.String("foo"),
				BinaryField: []byte("bar"),
			},
		},
		{
			name: "Graph",
			give: &ts.Graph{
				Edges: []*ts.Edge{
					{
						StartPoint: &ts.Point{X: 1.0, Y: 2.0},
						EndPoint:   &ts.Point{X: 3.0, Y: 4.0},
					},
					{
						StartPoint: &ts.Point{X: 5.0, Y: 6.0},
						EndPoint:   &ts.Point{X: 7.0, Y: 8.0},
					},
					{
						StartPoint: &ts.Point{X: 9.0, Y: 10.0},
						EndPoint:   &ts.Point{X: 11.0, Y: 12.0},
					},
				},
			},
		},
		{
			name: "ContainersOfContainers",
			give: &tc.ContainersOfContainers{
				ListOfLists: [][]int32{
					int32range(1, 10),
					int32range(2, 20),
					int32range(3, 30),
					int32range(4, 40),
					int32range(5, 50),
				},
				ListOfSets: []map[int32]struct{}{
					int32set(int32range(6, 60)...),
					int32set(int32range(7, 70)...),
					int32set(int32range(8, 80)...),
					int32set(int32range(9, 90)...),
					int32set(int32range(10, 100)...),
				},
				ListOfMaps: []map[int32]int32{
					int32multiply(42, int32range(5, 10)...),
					int32multiply(43, int32range(6, 20)...),
					int32multiply(44, int32range(7, 30)...),
				},
				SetOfSets: []map[string]struct{}{
					stringset("foo", "bar", "baz", "qux", "quux"),
					stringset("bar", "baz", "qux", "quux"),
					stringset("baz", "qux", "quux"),
					stringset("qux", "quux"),
					stringset("quux"),
					stringset(),
				},
				SetOfLists: [][]string{
					{"foo", "bar", "baz", "qux", "quux"},
					{"bar", "baz", "qux", "quux"},
					{"baz", "qux", "quux"},
					{"qux", "quux"},
					{"quux"},
					{},
				},
				SetOfMaps: []map[string]string{
					{"foo": "bar"},
					{"bar": "baz"},
					{"baz": "qux"},
					{"qux": "quux"},
					{"quux": "foo"},
				},
			},
		},
	}

	benchmarkEncode := func(b *testing.B, bb benchCase) {
		var buff bytes.Buffer

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buff.Reset()

			w, err := bb.give.ToWire()
			require.NoError(b, err, "ToWire")
			require.NoError(b, protocol.Binary.Encode(w, &buff), "Encode")
		}
	}

	benchmarkDecode := func(b *testing.B, bb benchCase) {
		var buff bytes.Buffer
		w, err := bb.give.ToWire()
		require.NoError(b, err, "ToWire")
		require.NoError(b, protocol.Binary.Encode(w, &buff), "Encode")

		r := bytes.NewReader(buff.Bytes())

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r.Seek(0, 0)

			wval, err := protocol.Binary.Decode(r, wire.TStruct)
			require.NoError(b, err, "Decode")

			require.NoError(b, bb.give.FromWire(wval), "FromWire")
		}
	}

	for _, bb := range benchmarks {
		b.Run(bb.name, func(b *testing.B) {
			b.Run("Encode", func(b *testing.B) {
				benchmarkEncode(b, bb)
			})

			b.Run("Decode", func(b *testing.B) {
				benchmarkDecode(b, bb)
			})
		})
	}
}

// Generates a slice representing the range [from, to).
func int32range(from, to int32) []int32 {
	if from > to {
		from, to = to, from
	}

	out := make([]int32, 0, to-from)
	for i := from; i < to; i++ {
		out = append(out, i)
	}
	return out
}

func int32set(i ...int32) map[int32]struct{} {
	o := make(map[int32]struct{}, len(i))
	for _, x := range i {
		o[x] = struct{}{}
	}
	return o
}

// maps provided numbers to the result of multipliying them with the provided
// number.
func int32multiply(with int32, nums ...int32) map[int32]int32 {
	o := make(map[int32]int32, len(nums))
	for _, x := range nums {
		o[x] = x * with
	}
	return o
}

func stringset(ss ...string) map[string]struct{} {
	o := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		o[s] = struct{}{}
	}
	return o
}
