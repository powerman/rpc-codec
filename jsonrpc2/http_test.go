package jsonrpc2_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestHTTPServer(t *testing.T) {
	const jBad = `{}`
	const jSum = `{"jsonrpc":"2.0","id":0,"method":"Svc.Sum","params":[3,5]}`
	const jNotify = `{"jsonrpc":"2.0","method":"Svc.Sum","params":[3,5]}`
	const jRes = `{"jsonrpc":"2.0","id":0,"result":8}`
	const jErr = `{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"Invalid request"}}`
	const jParse = `{"jsonrpc":"2.0","id":null,"error":{"code":-32700,"message":"Parse error"}}`
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
		{"POST", contentType, contentType, "", http.StatusOK, jParse},
		{"POST", contentType, contentType, " ", http.StatusOK, jParse},
		{"POST", contentType, contentType, "{", http.StatusOK, jParse},
		{"POST", contentType, contentType, `{"jsonrpc":"2.0",`, http.StatusOK, jParse},
		{"POST", contentType + "; charset=utf-8", contentType, jSum, http.StatusOK, jRes},
	}

	ts := httptest.NewServer(jsonrpc2.HTTPHandler(nil))
	// Don't close because of https://github.com/golang/go/issues/12262
	// defer ts.Close()

	for _, c := range cases {
		req, err := http.NewRequest(c.method, ts.URL, strings.NewReader(c.body))
		if err != nil {
			t.Errorf("NewRequest(%s %s), err = %v", c.method, ts.URL, err)
		}
		if c.contentType != "" {
			req.Header.Add("Content-Type", c.contentType)
		}
		if c.accept != "" {
			req.Header.Add("Accept", c.accept)
		}
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			t.Errorf("Do(%s %s), err = %v", c.method, ts.URL, err)
		}
		if resp.StatusCode != c.code {
			t.Errorf("Do(%s %s), status = %v, want = %v", c.method, ts.URL, resp.StatusCode, c.code)
		}
		if resp.Header.Get("Content-Type") != contentType {
			t.Errorf("Do(%s %s), Content-Type = %q, want = %q", c.method, ts.URL, resp.Header.Get("Content-Type"), contentType)
		}
		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("ReadAll(), err = %v", err)
		}
		if c.reply == "" {
			if len(got) != 0 {
				t.Errorf("Do(%s %s)\nexp: %#q\ngot: %#q", c.method, ts.URL, c.reply, string(bytes.TrimRight(got, "\n")))
			}
		} else {
			var jgot, jwant interface{}
			if err := json.Unmarshal(got, &jgot); err != nil {
				t.Errorf("Do(%s %s), output err = %v\ngot: %#q", c.method, ts.URL, err, string(bytes.TrimRight(got, "\n")))
			}
			if err := json.Unmarshal([]byte(c.reply), &jwant); err != nil {
				t.Errorf("Do(%s %s), expect err = %v\nexp: %#q", c.method, ts.URL, err, c.reply)
			}
			if !reflect.DeepEqual(jgot, jwant) {
				t.Errorf("Do(%s %s)\nexp: %#q\ngot: %#q", c.method, ts.URL, c.reply, string(bytes.TrimRight(got, "\n")))
			}
		}
	}
}

func TestHTTPClient(t *testing.T) {
	ts := httptest.NewServer(jsonrpc2.HTTPHandler(nil))
	// Don't close because of https://github.com/golang/go/issues/12262
	// defer ts.Close()
	client := jsonrpc2.NewHTTPClient(ts.URL)
	defer client.Close()

	var in [2]int
	var got, want int

	in, want = [2]int{1, 2}, 3
	err := client.Call("Svc.Sum", in, &got)
	if err != nil {
		t.Errorf("Call(%v), err = %v", in, err)
	}
	if got != want {
		t.Errorf("Call(%v) = %v, want = %v", in, got, want)
	}

	in = [2]int{2, 3}
	err = client.Notify("Svc.Sum", in)
	if err != nil {
		t.Errorf("Notify(%v), err = %v", in, err)
	}

	in, want = [2]int{3, 4}, 7
	err = client.Call("Svc.Sum", in, &got)
	if err != nil {
		t.Errorf("Call(%v), err = %v", in, err)
	}
	if got != want {
		t.Errorf("Call(%v) = %v, want = %v", in, got, want)
	}
}

type ContentTypeHanlder string

func (h ContentTypeHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", string(h))
	io.WriteString(w, `{"jsonrpc":"2.0","id":0,"result":8}`)
}

func TestHTTPClientContentType(t *testing.T) {
	const contentType = "application/json"
	cases := []struct {
		contentType string
		errorString string
	}{
		{contentType, "<nil>"},
		{contentType + "; charset=utf-8", "<nil>"},
		{contentType + "fail", `{"code":-32603,"message":"bad HTTP Content-Type: application/jsonfail"}`},
	}

	for _, c := range cases {
		ts := httptest.NewServer(ContentTypeHanlder(c.contentType))

		client := jsonrpc2.NewHTTPClient(ts.URL)
		in := []string{"ads"}
		var got int
		err := client.Call("Svc.Sum", in, &got)
		actualErrString := fmt.Sprintf("%v", err)
		if actualErrString != c.errorString {
			t.Errorf("Unexpected result. exp: %#q\ngot: %#q", actualErrString, c.errorString)
		}
	}
}
