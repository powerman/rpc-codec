package jsonrpc2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
	"strings"
	"testing"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Svc is an RPC service for testing.
type Svc struct{}

func (*Svc) Sum(vals [2]int, res *int) error {
	*res = vals[0] + vals[1]
	return nil
}

func init() {
	_ = rpc.Register(&Svc{})
}

var addr = getAddr()
var path = "/"
var url = "http://" + addr + path

func getAddr() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().String()
}

func init() {
	http.Handle(path, jsonrpc2.HTTPHandler(nil))
	go http.ListenAndServe(addr, nil)
}

func TestHTTPServer(t *testing.T) {
	const jBad = `{}`
	const jSum = `{"jsonrpc":"2.0","id":0,"method":"Svc.Sum","params":[3,5]}`
	const jNotify = `{"jsonrpc":"2.0","method":"Svc.Sum","params":[3,5]}`
	const jRes = `{"jsonrpc":"2.0","id":0,"result":8}`
	const jErr = `{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"Invalid request"}}`
	const contentType = "application/json"

	cases := []struct {
		method      string
		contentType string
		accept      string
		body        string
		code        int
		reply       string
	}{
		{"GET", "", "", "", http.StatusMethodNotAllowed, ""},
		{"POST", contentType, "", jSum, http.StatusUnsupportedMediaType, ""},
		{"POST", "text/json", contentType, jSum, http.StatusUnsupportedMediaType, ""},
		{"PUT", contentType, contentType, jSum, http.StatusMethodNotAllowed, ""},
		{"POST", contentType, contentType, jNotify, http.StatusNoContent, ""},
		{"POST", contentType, contentType, jSum, http.StatusOK, jRes},
		{"POST", contentType, contentType, jBad, http.StatusOK, jErr},
	}

	for _, c := range cases {
		req, err := http.NewRequest(c.method, url, strings.NewReader(c.body))
		if err != nil {
			t.Errorf("NewRequest(%s %s), err = %v", c.method, url, err)
		}
		if c.contentType != "" {
			req.Header.Add("Content-Type", c.contentType)
		}
		if c.accept != "" {
			req.Header.Add("Accept", c.accept)
		}
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			t.Errorf("Do(%s %s), err = %v", c.method, url, err)
		}
		if resp.StatusCode != c.code {
			t.Errorf("Do(%s %s), status = %v, want = %v", c.method, url, resp.StatusCode, c.code)
		}
		if resp.Header.Get("Content-Type") != contentType {
			t.Errorf("Do(%s %s), Content-Type = %q, want = %q", c.method, url, resp.Header.Get("Content-Type"), contentType)
		}
		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("ReadAll(), err = %v", err)
		}
		if c.reply == "" {
			if len(got) != 0 {
				t.Errorf("Do(%s %s)\nexp: %#q\ngot: %#q", c.method, url, c.reply, string(bytes.TrimRight(got, "\n")))
			}
		} else {
			var jgot, jwant interface{}
			if err := json.Unmarshal(got, &jgot); err != nil {
				t.Errorf("Do(%s %s), output err = %v\ngot: %#q", c.method, url, err, string(bytes.TrimRight(got, "\n")))
			}
			if err := json.Unmarshal([]byte(c.reply), &jwant); err != nil {
				t.Errorf("Do(%s %s), expect err = %v\nexp: %#q", c.method, url, err, c.reply)
			}
			if !reflect.DeepEqual(jgot, jwant) {
				t.Errorf("Do(%s %s)\nexp: %#q\ngot: %#q", c.method, url, c.reply, string(bytes.TrimRight(got, "\n")))
			}
		}
	}
}
