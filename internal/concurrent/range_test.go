package concurrent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeInvalid(t *testing.T) {
	tests := []struct {
		desc string
		c    interface{}
		f    interface{}
	}{
		{
			desc: "nil values",
		},
		{
			desc: "c: string",
			c:    "foo",
			f:    func(int, string) {},
		},
		{
			desc: "f: 0 args",
			c:    []int{1, 2, 3},
			f:    func() {},
		},
		{
			desc: "f: 1 args",
			c:    []int{1, 2, 3},
			f:    func(int) {},
		},
		{
			desc: "f: 3 args",
			c:    []int{1, 2, 3},
			f:    func(int, int, int) {},
		},
		{
			desc: "f: 2 returns",
			c:    []int{1, 2, 3},
			f:    func(int, int) (string, error) { return "", nil },
		},
		{
			desc: "f: return non-error",
			c:    []int{1, 2},
			f:    func(int, int) string { return "" },
		},
		{
			desc: "f: slice: arg 1 type mismatch",
			c:    []string{},
			f:    func(string, int) {},
		},
		{
			desc: "f: slice: arg 2 type mismatch",
			c:    []string{},
			f:    func(int, int) {},
		},
		{
			desc: "f: map: arg 1 type mismatch",
			c:    map[string]int{},
			f:    func(int, int) {},
		},
		{
			desc: "f: map: arg 2 type mismatch",
			c:    map[string]int{},
			f:    func(string, string) {},
		},
	}

	for _, tt := range tests {
		assert.Panics(t, func() {
			Range(tt.c, tt.f)
		}, tt.desc)
	}
}

func TestRangeSlice(t *testing.T) {
	tests := []struct {
		desc string
		c    interface{}
		f    interface{}

		shouldFail bool
	}{
		{
			desc: "slice: empty",
			c:    []string{},
			f: func(int, string) {
				t.Errorf("unexpected call to function for empty collection")
			},
		},
		{
			desc: "slice: no return",
			c:    []int32{1, 2, 3},
			f: func(i int, v int32) {
				switch i {
				case 0:
					assert.EqualValues(t, 1, v)
				case 1:
					assert.EqualValues(t, 2, v)
				case 2:
					assert.EqualValues(t, 3, v)
				}
			},
		},
		{
			desc: "slice: fail all",
			c:    []int32{1, 2, 3},
			f: func(i int, v int32) error {
				switch i {
				case 0:
					assert.EqualValues(t, 1, v)
				case 1:
					assert.EqualValues(t, 2, v)
				case 2:
					assert.EqualValues(t, 3, v)
				}
				return errors.New("foo")
			},
			shouldFail: true,
		},
		{
			desc: "slice: fail one",
			c:    []string{"hello", "world"},
			f: func(i int, v string) error {
				switch i {
				case 0:
					assert.Equal(t, "hello", v)
				case 1:
					assert.Equal(t, "world", v)
					return errors.New("foo")
				}
				return nil
			},
			shouldFail: true,
		},
		{
			desc: "map: empty",
			c:    map[string]int{},
			f: func(string, int) {
				t.Errorf("unexpected call to function for empty collection")
			},
		},
		{
			desc: "map: no return",
			c: map[string]int{
				"hello": 1,
				"world": 2,
			},
			f: func(k string, v int) {
				switch k {
				case "hello":
					assert.Equal(t, v, 1)
				case "world":
					assert.Equal(t, v, 2)
				}
			},
		},
		{
			desc: "map: fail all",
			c: map[string]int{
				"hello": 1,
				"world": 2,
			},
			f: func(k string, v int) error {
				switch k {
				case "hello":
					assert.Equal(t, v, 1)
				case "world":
					assert.Equal(t, v, 2)
				}
				return errors.New("great sadness")
			},
			shouldFail: true,
		},
		{
			desc: "map: fail one",
			c: map[string]int{
				"hello": 1,
				"world": 2,
			},
			f: func(k string, v int) error {
				switch k {
				case "hello":
					assert.Equal(t, v, 1)
				case "world":
					assert.Equal(t, v, 2)
					return errors.New("great sadness")
				}
				return nil
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		err := Range(tt.c, tt.f)
		if tt.shouldFail {
			assert.Error(t, err, tt.desc)
		} else {
			assert.NoError(t, err, tt.desc)
		}
	}
}
