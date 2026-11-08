package jsonrpc2_test

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
	"slices"
	"testing"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

type contextKey int

var remoteAddrContextKey contextKey

func init() {
	_ = rpc.Register(&CtxSvc{})
}

// CtxSvc is an RPC service for testing.
type CtxSvc struct{}

func (*CtxSvc) Sum(vals [2]int, res *int) error {
	*res = vals[0] + vals[1]
	return nil
}

// type NameArg struct{ Fname, Lname string }
// type NameRes struct{ Name string }

func (*CtxSvc) Name(t NameArg, res *NameRes) error {
	*res = NameRes{t.Fname + " " + t.Lname}
	return nil
}

// type NameArgContext struct {
// 	Fname, Lname string
// 	jsonrpc2.Ctx
// }

type NameResCtx struct {
	Name           string
	TCPRemoteAddr  string
	HTTPRemoteAddr string
}

// Method with named params and TCP context.
func (*CtxSvc) NameCtx(t NameArgContext, res *NameResCtx) error {
	*res = NameResCtx{Name: t.Fname + " " + t.Lname}
	if val := t.Context().Value(remoteAddrContextKey); val != nil {
		res.TCPRemoteAddr, _, _ = net.SplitHostPort(val.(*net.TCPAddr).String())
	}
	if val := jsonrpc2.HTTPRequestFromContext(t.Context()); val != nil {
		res.HTTPRemoteAddr, _, _ = net.SplitHostPort(val.RemoteAddr)
	}
	return nil
}

func TestNoContext(t *testing.T) {
	t.Parallel()
	var c jsonrpc2.Ctx
	if c.Context() != context.Background() { //nolint:usetesting // Must be exactly context.Background.
		t.Fatal("empty context")
	}
}

// - for each of these servers:
//   - TCP server without context
//   - TCP server with context
//   - HTTP server with context
//
// - call these methods:
//   - Sum()
//   - Name()
//   - NameCtx()
//   - TODO batch call all
func TestContext(t *testing.T) {
	t.Parallel()
	var lc net.ListenConfig

	// Server provide a TCP transport without context.
	serverTCPNoCtx, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { serverTCPNoCtx.Close() })
	go func() {
		for {
			conn, err := serverTCPNoCtx.Accept()
			if err != nil {
				return
			}
			go jsonrpc2.ServeConn(conn)
		}
	}()

	// Server provide a TCP transport with context.
	serverTCP, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { serverTCP.Close() })
	go func() {
		for {
			conn, err := serverTCP.Accept()
			if err != nil {
				return
			}
			ctx := context.WithValue(t.Context(), remoteAddrContextKey, conn.RemoteAddr())
			go jsonrpc2.ServeConnContext(ctx, conn)
		}
	}()

	// Server provide a HTTP transport with context.
	mux := http.NewServeMux()
	mux.Handle("/", jsonrpc2.HTTPHandler(nil))
	serverHTTP, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { serverHTTP.Close() })
	go http.Serve(serverHTTP, mux)

	// Client use TCP transport without context.
	clientTCPNoCtx, err := jsonrpc2.Dial("tcp", serverTCPNoCtx.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { clientTCPNoCtx.Close() })

	// Client use TCP transport with context.
	clientTCP, err := jsonrpc2.Dial("tcp", serverTCP.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { clientTCP.Close() })

	// Client use HTTP transport with context.
	clientHTTP := jsonrpc2.NewHTTPClient("http://" + serverHTTP.Addr().String() + "/")
	t.Cleanup(func() { clientHTTP.Close() })

	cases := []struct {
		client     *jsonrpc2.Client
		clientName string
		method     string
		arg        any
		want       any
	}{
		{
			clientTCPNoCtx, "clientTCPNoCtx",
			"CtxSvc.Sum",
			[2]int{3, 5},
			8.0,
		},
		{
			clientTCP, "clientTCP",
			"CtxSvc.Sum",
			[2]int{3, 5},
			8.0,
		},
		{
			clientHTTP, "clientHTTP",
			"CtxSvc.Sum",
			[2]int{3, 5},
			8.0,
		},
		{
			clientTCPNoCtx, "clientTCPNoCtx",
			"CtxSvc.Name",
			NameArg{"First", "Last"},
			map[string]any{
				"Name": "First Last",
			},
		},
		{
			clientTCP, "clientTCP",
			"CtxSvc.Name",
			NameArg{"First", "Last"},
			map[string]any{
				"Name": "First Last",
			},
		},
		{
			clientHTTP, "clientHTTP",
			"CtxSvc.Name",
			NameArg{"First", "Last"},
			map[string]any{
				"Name": "First Last",
			},
		},
		{
			clientTCPNoCtx, "clientTCPNoCtx",
			"CtxSvc.NameCtx",
			NameArg{"First", "Last"},
			map[string]any{
				"Name":           "First Last",
				"TCPRemoteAddr":  "",
				"HTTPRemoteAddr": "",
			},
		},
		{
			clientTCP, "clientTCP",
			"CtxSvc.NameCtx",
			NameArg{"First", "Last"},
			map[string]any{
				"Name":           "First Last",
				"TCPRemoteAddr":  "127.0.0.1",
				"HTTPRemoteAddr": "",
			},
		},
		{
			clientHTTP, "clientHTTP",
			"CtxSvc.NameCtx",
			NameArg{"First", "Last"},
			map[string]any{
				"Name":           "First Last",
				"TCPRemoteAddr":  "",
				"HTTPRemoteAddr": "127.0.0.1",
			},
		},
	}
	for _, v := range cases {
		var res any
		err := v.client.Call(v.method, v.arg, &res)
		if err != nil {
			t.Errorf("%s.Call(%q) = %v", v.clientName, v.method, err)
		}
		if !reflect.DeepEqual(v.want, res) {
			t.Errorf("%s.Call(%q):\n\n\texp: %#v\n\n\tgot: %#v\n\n", v.clientName, v.method, v.want, res)
		}
	}
}

func TestContextBatch(t *testing.T) {
	t.Parallel()
	req := `[
		{"jsonrpc":"2.0","id":0,"method":"CtxSvc.NameCtx","params":{}},
		{"jsonrpc":"2.0","id":1,"method":"CtxSvc.NameCtx","params":{"Fname":"First","Lname":"Last"}}
		]`
	want := []map[string]any{
		{
			"jsonrpc": "2.0",
			"id":      0.0,
			"result": map[string]any{
				"Name":           " ",
				"TCPRemoteAddr":  "127.0.0.1",
				"HTTPRemoteAddr": "",
			},
		},
		{
			"jsonrpc": "2.0",
			"id":      1.0,
			"result": map[string]any{
				"Name":           "First Last",
				"TCPRemoteAddr":  "127.0.0.1",
				"HTTPRemoteAddr": "",
			},
		},
	}
	buf := bytes.NewBufferString(req)
	rpc.ServeRequest(jsonrpc2.NewServerCodecContext(
		context.WithValue(t.Context(), remoteAddrContextKey, &net.TCPAddr{IP: net.ParseIP("127.0.0.1")}),
		&bufReadWriteCloser{buf},
		nil,
	))
	var res []map[string]any
	err := json.Unmarshal(buf.Bytes(), &res)
	if err != nil {
		t.Fatal(err)
	}
	slices.SortFunc(res, func(a, b map[string]any) int { return cmp.Compare(a["id"].(float64), b["id"].(float64)) })
	if !reflect.DeepEqual(want, res) {
		t.Errorf("%s:\n\n\texp: %#v\n\n\tgot: %#v\n\n", req, want, res)
	}
}

func TestContextBugNoContextWithoutParams(t *testing.T) {
	t.Parallel()
	cases := []struct {
		req  string
		want any
	}{
		{
			`{"jsonrpc":"2.0","id":0,"method":"CtxSvc.NameCtx"}`,
			map[string]any{
				"Name":           " ",
				"TCPRemoteAddr":  "127.0.0.1",
				"HTTPRemoteAddr": "",
			},
		},
		{
			`{"jsonrpc":"2.0","id":0,"method":"CtxSvc.NameCtx","params":{}}`,
			map[string]any{
				"Name":           " ",
				"TCPRemoteAddr":  "127.0.0.1",
				"HTTPRemoteAddr": "",
			},
		},
		{
			`{"jsonrpc":"2.0","id":0,"method":"CtxSvc.NameCtx","params":{"Fname":"First","Lname":"Last"}}`,
			map[string]any{
				"Name":           "First Last",
				"TCPRemoteAddr":  "127.0.0.1",
				"HTTPRemoteAddr": "",
			},
		},
	}
	for _, v := range cases {
		buf := bytes.NewBufferString(v.req)
		rpc.ServeRequest(jsonrpc2.NewServerCodecContext(
			context.WithValue(t.Context(), remoteAddrContextKey, &net.TCPAddr{IP: net.ParseIP("127.0.0.1")}),
			&bufReadWriteCloser{buf},
			nil,
		))
		var res map[string]any
		err := json.Unmarshal(buf.Bytes(), &res)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v.want, res["result"]) {
			t.Errorf("%q:\n\n\texp: %#v\n\n\tgot: %#v\n\n", v.req, v.want, res["result"])
		}
	}
}

type bufReadWriteCloser struct {
	*bytes.Buffer
}

func (*bufReadWriteCloser) Close() error { return nil }
