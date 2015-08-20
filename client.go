// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsonrpc2 implements a JSON-RPC 2.0 ClientCodec and ServerCodec
// for the rpc package.
package jsonrpc2

import (
	"encoding/json"
	"io"
	"math"
	"net"
	"net/rpc"
	"reflect"
	"sync"
)

const seqNotify = math.MaxUint64

type clientCodec struct {
	dec *json.Decoder // for reading JSON values
	enc *json.Encoder // for writing JSON values
	c   io.Closer

	// temporary work space
	resp clientResponse

	// JSON-RPC responses include the request id but not the request method.
	// Package rpc expects both.
	// We save the request method in pending when sending a request
	// and then look it up by request ID when filling out the rpc Response.
	mutex   sync.Mutex        // protects pending
	pending map[uint64]string // map request id to method name
}

// NewClientCodec returns a new rpc.ClientCodec using JSON-RPC on conn.
func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{
		dec:     json.NewDecoder(conn),
		enc:     json.NewEncoder(conn),
		c:       conn,
		pending: make(map[uint64]string),
	}
}

type clientRequest struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	Id      *uint64     `json:"id,omitempty"`
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) error {
	// If return error: it will be returned as is for this call.
	// Allow param to be only Array, Slice, Map or Struct.
	// When param is nil or uninitialized Map or Slice - omit "params".
	if param != nil {
		switch k := reflect.TypeOf(param).Kind(); k {
		case reflect.Map:
			if reflect.TypeOf(param).Key().Kind() == reflect.String {
				if reflect.ValueOf(param).IsNil() {
					param = nil
				}
			}
		case reflect.Slice:
			if reflect.ValueOf(param).IsNil() {
				param = nil
			}
		case reflect.Array, reflect.Struct:
		case reflect.Ptr:
			switch k := reflect.TypeOf(param).Elem().Kind(); k {
			case reflect.Map:
				if reflect.TypeOf(param).Elem().Key().Kind() == reflect.String {
					if reflect.ValueOf(param).Elem().IsNil() {
						param = nil
					}
				}
			case reflect.Slice:
				if reflect.ValueOf(param).Elem().IsNil() {
					param = nil
				}
			case reflect.Array, reflect.Struct:
			default:
				return NewError(errInternal.Code, "unsupported param type: Ptr to "+k.String())
			}
		default:
			return NewError(errInternal.Code, "unsupported param type: "+k.String())
		}
	}

	var req clientRequest
	if r.Seq != seqNotify {
		c.mutex.Lock()
		c.pending[r.Seq] = r.ServiceMethod
		c.mutex.Unlock()
		req.Id = &r.Seq
	}
	req.Version = "2.0"
	req.Method = r.ServiceMethod
	req.Params = param
	if err := c.enc.Encode(&req); err != nil {
		return NewError(errInternal.Code, err.Error())
	}
	return nil
}

type clientResponse struct {
	Version string           `json:"jsonrpc"`
	Id      uint64           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   *Error           `json:"error"`
}

func (r *clientResponse) reset() {
	r.Version = ""
	r.Id = 0
	r.Result = nil
	r.Error = nil
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	// If return err:
	// - io.EOF will became ErrShutdown or io.ErrUnexpectedEOF
	// - it will be returned as is for all pending calls
	// - client will be shutdown
	// So, return io.EOF as is, return *Error for all other errors.
	var raw json.RawMessage
	if err := c.dec.Decode(&raw); err != nil {
		if err == io.EOF {
			return err
		}
		return NewError(errInternal.Code, err.Error())
	}

	var rawMap = make(map[string]*json.RawMessage)
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}
	if len(rawMap) != 3 {
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}
	if rawMap["jsonrpc"] == nil {
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}
	if _, ok := rawMap["id"]; !ok {
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}

	c.resp.reset()
	if err := json.Unmarshal(raw, &c.resp); err != nil {
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}
	if c.resp.Version != "2.0" {
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}
	if _, ok := rawMap["result"]; ok && c.resp.Result == nil {
		c.resp.Result = &null
	}
	if c.resp.Result == nil && c.resp.Error == nil || c.resp.Result != nil && c.resp.Error != nil {
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}
	if c.resp.Error != nil {
		if rawMap["error"] == nil {
			return NewError(errInternal.Code, "bad response: "+string(raw))
		}
		rawErrMap := make(map[string]*json.RawMessage)
		if err := json.Unmarshal(*rawMap["error"], &rawErrMap); err != nil {
			return NewError(errInternal.Code, "bad response: "+string(raw))
		}
		if rawErrMap["code"] == nil {
			return NewError(errInternal.Code, "bad response: "+string(raw))
		}
		if rawErrMap["message"] == nil {
			return NewError(errInternal.Code, "bad response: "+string(raw))
		}
		if _, ok := rawErrMap["data"]; len(rawErrMap) == 3 && !ok {
			return NewError(errInternal.Code, "bad response: "+string(raw))
		}
		if len(rawErrMap) > 3 {
			return NewError(errInternal.Code, "bad response: "+string(raw))
		}
	}
	if rawMap["id"] == nil {
		if c.resp.Error != nil {
			return c.resp.Error
		}
		return NewError(errInternal.Code, "bad response: "+string(raw))
	}

	c.mutex.Lock()
	r.ServiceMethod = c.pending[c.resp.Id]
	delete(c.pending, c.resp.Id)
	c.mutex.Unlock()

	r.Error = ""
	r.Seq = c.resp.Id
	if c.resp.Error != nil {
		r.Error = c.resp.Error.Error()
	}
	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	// If x!=nil and return error e:
	// - this call get e.Error() appended to "reading body "
	// - other pending calls get error as is XXX actually other calls
	//   shouldn't be affected by this error at all, so let's at least
	//   provide different error message for other calls
	if x == nil {
		return nil
	}
	if err := json.Unmarshal(*c.resp.Result, x); err != nil {
		e := NewError(errInternal.Code, err.Error())
		e.Data = NewError(errInternal.Code, "some other Call failed to unmarshal Reply")
		return e
	}
	return nil
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}

// Client represents a JSON RPC 2.0 Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
//
// It also provides all methods of net/rpc.Client.
type Client struct {
	*rpc.Client
	codec *clientCodec
}

// Notify try to invoke the named function. It return error only in case
// it wasn't able to send request.
func (c Client) Notify(serviceMethod string, args interface{}) error {
	req := &rpc.Request{
		ServiceMethod: serviceMethod,
		Seq:           seqNotify,
	}
	return c.codec.WriteRequest(req, args)
}

// NewClient returns a new Client to handle requests to the
// set of services at the other end of the connection.
func NewClient(conn io.ReadWriteCloser) *Client {
	codec := NewClientCodec(conn)
	client := rpc.NewClientWithCodec(codec)
	return &Client{client, codec.(*clientCodec)}
}

// Dial connects to a JSON-RPC server at the specified network address.
func Dial(network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), err
}
