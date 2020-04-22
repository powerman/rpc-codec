package jsonrpc2

import (
	"io"
	"net/rpc"
	"testing"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		err error
	}{
		{nil},
		{rpc.ErrShutdown},
		{io.ErrUnexpectedEOF},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(t *testing.T) {
			err := WrapError(tc.err)
			if err != tc.err {
				t.Errorf("got %v, want %v", err, tc.err)
			}
		})
	}

}
