package jsonrpc2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/rpc"
)

const (
	contentType = "application/json"
	// Max size of unneeded response body to read so underlying TCP connection will be re-used.
	maxBodySlurpSize = 32 * 1024
	// Max amount of received responses not yet read by [rpc.Client].
	readyQueueSize = 16
)

var (
	errBadContentType = errors.New("bad HTTP Content-Type")
	errBadStatus      = errors.New("bad HTTP Status")
)

type contextKey int

const httpRequestContextKey contextKey = 0

// HTTPRequestFromContext returns HTTP request related to this RPC (if
// you use HTTPHander to serve JSON RPC 2.0 over HTTP) or nil otherwise.
func HTTPRequestFromContext(ctx context.Context) *http.Request {
	req, _ := ctx.Value(httpRequestContextKey).(*http.Request)
	return req
}

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

func (*httpServerConn) Close() error {
	return nil
}

type httpHandler struct {
	rpc *rpc.Server
}

// HTTPHandler returns handler for HTTP requests which will execute
// incoming JSON-RPC 2.0 over HTTP using srv.
//
// If srv is nil then [rpc.DefaultServer] will be used.
//
// Specification: http://www.simple-is-better.org/json-rpc/transport_http.html
func HTTPHandler(srv *rpc.Server) http.Handler {
	if srv == nil {
		srv = rpc.DefaultServer
	}
	return &httpHandler{srv}
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", contentType)

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mediaType, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if mediaType != contentType || req.Header.Get("Accept") != contentType {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	ctx := context.WithValue(context.Background(), httpRequestContextKey, req)
	conn := &httpServerConn{req: req.Body, res: w}
	_ = h.rpc.ServeRequest(NewServerCodecContext(ctx, conn, h.rpc))
	if !conn.replied {
		w.WriteHeader(http.StatusNoContent)
	}
}

// Doer is an interface for doing HTTP requests.
type Doer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// The DoerFunc type is an adapter to allow the use of ordinary
// functions as HTTP clients. If f is a function with the appropriate
// signature, DoerFunc(f) is a client that calls f.
type DoerFunc func(req *http.Request) (resp *http.Response, err error)

// Do calls f(req).
func (f DoerFunc) Do(req *http.Request) (resp *http.Response, err error) {
	return f(req)
}

type httpClientConn struct {
	url   string
	doer  Doer
	ready chan io.ReadCloser
	body  io.ReadCloser
	close chan struct{}
}

func (conn *httpClientConn) Read(buf []byte) (int, error) {
	if conn.body == nil {
		select {
		case <-conn.close:
			return 0, io.EOF
		case conn.body = <-conn.ready:
		}
	}
	n, err := conn.body.Read(buf)
	if err == io.EOF {
		logIfFail(conn.body.Close)
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
	go conn.roundTrip(b)
	return len(buf), nil
}

// roundTrip sends request b and queues either a response body
// or a response with error (unless b is a Notification) for Read.
func (conn *httpClientConn) roundTrip(b []byte) {
	err := conn.post(b)
	if err == nil {
		return
	}
	var res clientResponse
	if json.Unmarshal(b, &res) == nil && res.ID == nil {
		return // ignore error from Notification
	}
	res.Error = NewError(errInternal.Code, err.Error())
	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(res) //nolint:errchkjson // Result is a valid JSON or nil, Error.Data is nil.
	conn.ready <- io.NopCloser(buf)
}

// post sends request b and queues response body for Read (if any).
func (conn *httpClientConn) post(b []byte) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, conn.url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	resp, err := conn.doer.Do(req) //nolint:bodyclose // Closed by discardBody or by Read.
	if err != nil {
		if resp != nil {
			discardBody(resp)
		}
		return err
	}

	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch {
	case mediaType != contentType || err != nil:
		err = fmt.Errorf("%w: %s", errBadContentType, resp.Header.Get("Content-Type"))
	case resp.StatusCode == http.StatusOK:
		conn.ready <- resp.Body
		return nil
	case resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusAccepted:
		err = nil
	default:
		err = fmt.Errorf("%w: %s", errBadStatus, resp.Status)
	}
	discardBody(resp)
	return err
}

// discardBody closes resp.Body after reading it if small,
// so underlying TCP connection will be re-used.
// No need to check for read errors: if it fails, Transport won't reuse it anyway.
func discardBody(resp *http.Response) {
	if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
		_, _ = io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
	}
	logIfFail(resp.Body.Close)
}

func (conn *httpClientConn) Close() error {
	close(conn.close)
	return nil
}

// NewHTTPClient returns a new Client to handle requests to the
// set of services at the given url.
func NewHTTPClient(url string) *Client {
	return NewCustomHTTPClient(url, nil)
}

// NewCustomHTTPClient returns a new Client to handle requests to the
// set of services at the given url using provided doer (&http.Client{} by
// default).
//
// Use doer to customize HTTP authorization/headers/etc. sent with each
// request (it method Do() will receive already configured POST request
// with url, all required headers and body set according to specification).
func NewCustomHTTPClient(url string, doer Doer) *Client {
	if doer == nil {
		doer = &http.Client{}
	}
	return NewClient(&httpClientConn{
		url:   url,
		doer:  doer,
		ready: make(chan io.ReadCloser, readyQueueSize),
		close: make(chan struct{}),
	})
}
