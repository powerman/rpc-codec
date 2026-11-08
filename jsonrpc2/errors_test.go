package jsonrpc2_test

import (
	"io"
	"net/rpc"
	"testing"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

func TestWrapError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		err error
	}{
		{nil},
		{rpc.ErrShutdown},
		{io.ErrUnexpectedEOF},
	}
	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			err := jsonrpc2.WrapError(tc.err)
			if err != tc.err { //nolint:err113,errorlint // By design.
				t.Errorf("got %v, want %v", err, tc.err)
			}
		})
	}
}
