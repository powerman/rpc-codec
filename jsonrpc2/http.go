package jsonrpc2

import (
	"io"
	"net/http"
	"net/rpc"
)

const contentType = "application/json"

type httpConn struct {
	req     io.Reader
	res     io.Writer
	replied bool
}

func (conn *httpConn) Read(buf []byte) (int, error) {
	return conn.req.Read(buf)
}

func (conn *httpConn) Write(buf []byte) (int, error) {
	conn.replied = true
	return conn.res.Write(buf)
}

func (conn *httpConn) Close() error {
	return nil
}

type httpHandler struct {
	rpc *rpc.Server
}

// HTTPHandler returns handler for HTTP requests which will execute
// incoming RPC using srv. If srv is nil then use rpc.DefaultServer.
//
// Specification: http://www.simple-is-better.org/json-rpc/transport_http.html
// Pipelined Requests/Responses not supported.
// GET Request not supported.
func HTTPHandler(srv *rpc.Server) http.Handler {
	if srv == nil {
		srv = rpc.DefaultServer
	}
	return &httpHandler{srv}
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", contentType)

	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Content-Type") != contentType || req.Header.Get("Accept") != contentType {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	conn := &httpConn{req: req.Body, res: w}
	h.rpc.ServeRequest(NewServerCodec(conn, h.rpc))
	if !conn.replied {
		w.WriteHeader(http.StatusNoContent)
	}
}
