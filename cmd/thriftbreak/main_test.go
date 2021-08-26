package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThriftBreak(t *testing.T) {
	t.Parallel()
	t.Run("wrong flag", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--out_file=tests"})
		require.Error(t, err)
		assert.EqualError(t, err, "flag provided but not defined: -out_file")
	})
	t.Run("wrong file name", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/something.thrift"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no such file")
	})
	t.Run("invalid thrift", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/invalid.thrift"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse file")
	})
	t.Run("invalid thrift for from_file", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/v1.thrift","--from_file=tests/invalid.thrift"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse file")
	})
	t.Run("missing to_file", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--from_file=tests/something.thrift"})
		require.Error(t, err)
		assert.Equal(t, err.Error(), "must provide an updated Thrift file")
	})
	t.Run("integration test", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/v2.thrift", "--from_file=tests/v1.thrift"})
		require.Error(t, err)

		assert.Equal(t, "removing method methodA in service Foo is not backwards compatible;" +
			" deleting service Bar is not backwards compatible;" +
			" changing an optional field B in AddedRequiredField to required is not backwards compatible;" +
			" adding a required field C to AddedRequiredField is not backwards compatible",
			err.Error())
	})
	t.Run("integration test", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/v3.thrift", "--from_file=tests/v1.thrift"})
		require.Error(t, err)

		assert.Equal(t,  "removing method methodA in service Foo is not backwards compatible", err.Error())
	})
}
