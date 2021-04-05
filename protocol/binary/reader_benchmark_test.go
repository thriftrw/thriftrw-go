package binary

import (
	"bytes"
	"io"
	"testing"

	"go.uber.org/thriftrw/wire"
)

func Benchmark_readStruct(b *testing.B) {
	benchmarks := []struct {
		name      string
		numFields int
	}{
		{name: "small struct with 64 fields", numFields: 64},
		{name: "large struct with 512 fields", numFields: 512},
	}

	for _, bb := range benchmarks {
		b.Run(bb.name, func(b *testing.B) {
			reader := NewReader(constructThriftStruct(bb.numFields))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				reader.readStruct(0)
			}
		})
	}
}

// constructThriftStruct constructs a message containing a struct with n bool fields.
func constructThriftStruct(n int) io.ReaderAt {
	var buf bytes.Buffer
	w := BorrowWriter(&buf)

	fields := make([]wire.Field, 0, n)
	boolVal := wire.NewValueBool(false)
	for i := 0; i < n; i++ {
		// Add a bool field
		fields = append(fields, wire.Field{
			ID:    int16(i),
			Value: boolVal,
		})
	}
	structVal := wire.NewValueStruct(wire.Struct{Fields: fields})

	w.WriteValue(structVal)
	return bytes.NewReader(buf.Bytes())
}
