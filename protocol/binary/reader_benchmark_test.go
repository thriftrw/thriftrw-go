package binary

import (
	"bytes"
	"io"
	"testing"

	"go.uber.org/thriftrw/wire"
)

func Benchmark_readStruct(b *testing.B) {
	for _, bench := range []struct {
		name      string
		numFields int
	}{
		{
			name:      "small struct with 64 fields",
			numFields: 64,
		},
		{
			name:      "large struct with 512 fields",
			numFields: 512,
		},
	} {
		reader := NewReader(constructThriftStruct(bench.numFields))

		b.Run(bench.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				reader.readStruct(0)
			}
		})
	}
}

// constructThriftStruct constructs a message containing a struct with n bool fields.
func constructThriftStruct(n int) io.ReaderAt {
	var buf bytes.Buffer
	buf.Grow(n*4 + 1)
	boolField := []byte{byte(wire.TBool), 0, 0, 0}

	for i := 0; i < n; i++ {
		buf.Write(boolField)
	}
	// A zero byte to terminate the message
	buf.WriteByte(0)

	return bytes.NewReader(buf.Bytes())
}
