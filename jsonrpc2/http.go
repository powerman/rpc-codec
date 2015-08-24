package jsonrpc2

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
)

const contentType = "application/json"

type httpServerConn struct {
	req     io.Reader
	res     io.Writer
	replied bool
}

func (conn *httpServerConn) Read(buf []byte) (int, error) {
	return conn.req.Read(buf)
}

func (conn *httpServerConn) Write(buf []byte) (int, error) {
	conn.replied = true
	return conn.res.Write(buf)
}

func (conn *httpServerConn) Close() error {
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

	conn := &httpServerConn{req: req.Body, res: w}
	h.rpc.ServeRequest(NewServerCodec(conn, h.rpc))
	if !conn.replied {
		w.WriteHeader(http.StatusNoContent)
	}
}

type httpClientConn struct {
	url   string
	ready chan io.ReadCloser
	body  io.ReadCloser
}

func (conn *httpClientConn) Read(buf []byte) (int, error) {
	if conn.body == nil {
		conn.body = <-conn.ready
	}
	n, err := conn.body.Read(buf)
	if err == io.EOF {
		conn.body.Close()
		conn.body = nil
		err = nil
		if n == 0 {
			return conn.Read(buf)
		}
	}
	return n, err
}

func (conn *httpClientConn) Write(buf []byte) (int, error) {
	req, err := http.NewRequest("POST", conn.url, bytes.NewReader(buf))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		return 0, fmt.Errorf("bad HTTP Status: %s", resp.Status)
	}
	if resp.Header.Get("Content-Type") != contentType {
		return 0, fmt.Errorf("bad HTTP Content-Type: %s", resp.Header.Get("Content-Type"))
	}
	if resp.StatusCode == http.StatusOK {
		conn.ready <- resp.Body
	}
	return len(buf), nil
}

func (conn *httpClientConn) Close() error {
	return nil
}

// DialHTTP connects to a JSON-RPC 2.0 server using HTTP at the specified url.
func DialHTTP(url string) *Client {
	return NewClient(&httpClientConn{url: url, ready: make(chan io.ReadCloser, 16)})
}
