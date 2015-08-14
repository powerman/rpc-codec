// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsonrpc2 implements a JSON-RPC 2.0 ClientCodec and ServerCodec
// for the rpc package.
package jsonrpc2

import (
	"encoding/json"
	"io"
	"net"
	"net/rpc"
	"reflect"
	"sync"
)

type clientCodec struct {
	dec *json.Decoder // for reading JSON values
	enc *json.Encoder // for writing JSON values
	c   io.Closer

	// temporary work space
	req  clientRequest
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
	Id      uint64      `json:"id"`
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

	c.mutex.Lock()
	c.pending[r.Seq] = r.ServiceMethod
	c.mutex.Unlock()
	c.req.Version = "2.0"
	c.req.Method = r.ServiceMethod
	c.req.Params = param
	c.req.Id = r.Seq
	if err := c.enc.Encode(&c.req); err != nil {
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
		if len(rawErrMap) < 2 || len(rawErrMap) > 4 {
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
	// - other pending calls get error as is
	if x == nil {
		return nil
	}
	if err := json.Unmarshal(*c.resp.Result, x); err != nil {
		return NewError(errInternal.Code, err.Error())
	}
	return nil
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}

// NewClient returns a new rpc.Client to handle requests to the
// set of services at the other end of the connection.
func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}

// Dial connects to a JSON-RPC server at the specified network address.
func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), err
}
