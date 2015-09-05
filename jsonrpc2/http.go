package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
//   - Pipelined Requests/Responses not supported.
//   - GET Request not supported.
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
	b := make([]byte, len(buf))
	copy(b, buf)
	go func() {
		req, err := http.NewRequest("POST", conn.url, bytes.NewReader(b))
		if err == nil {
			req.Header.Add("Content-Type", contentType)
			req.Header.Add("Accept", contentType)
			var resp *http.Response
			resp, err = (&http.Client{}).Do(req)
			if err != nil {
			} else if resp.Header.Get("Content-Type") != contentType {
				err = fmt.Errorf("bad HTTP Content-Type: %s", resp.Header.Get("Content-Type"))
			} else if resp.StatusCode == http.StatusOK {
				conn.ready <- resp.Body
				return
			} else if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusAccepted {
				resp.Body.Close() // is it should be read to EOF first?
				return
			} else {
				err = fmt.Errorf("bad HTTP Status: %s", resp.Status)
			}
			if resp != nil {
				resp.Body.Close() // is it should be read to EOF first?
			}
		}
		var res clientResponse
		if json.Unmarshal(b, &res) == nil && res.ID == nil {
			return // ignore error from Notification
		}
		res.Error = NewError(errInternal.Code, err.Error())
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(res)
		conn.ready <- ioutil.NopCloser(buf)
	}()
	return len(buf), nil
}

func (conn *httpClientConn) Close() error {
	return nil
}

// NewHTTPClient returns a new Client to handle requests to the
// set of services at the given url.
func NewHTTPClient(url string) *Client {
	return NewClient(&httpClientConn{url: url, ready: make(chan io.ReadCloser, 16)})
}
