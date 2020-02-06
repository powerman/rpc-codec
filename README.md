# JSON-RPC 2.0 [![GoDoc](https://godoc.org/github.com/seanhagen/rpc-codec/jsonrpc2?status.svg)](http://godoc.org/github.com/powerman/rpc-codec/jsonrpc2) 

jsonrpc2 is a codec for net/rpc.

Forked to add some handling for RPC services that don't perfectly adhere to the spec.

Implements [JSON-RPC 2.0](http://www.jsonrpc.org/specification) and
[JSON-RPC 2.0 Transport: HTTP](http://www.simple-is-better.org/json-rpc/transport_http.html)
specifications with following limitations:

- Client: Batch Request not supported.
- HTTP Client&Server: Pipelined Requests/Responses not supported.
- HTTP Client&Server: GET Request not supported.

Also provides command-line tool `jsonrpc2client`.


## Installation

```sh
go get github.com/powerman/rpc-codec/...
```

Or, if you need only `jsonrpc2client` tool, then download binary for your
OS manually from
[releases](https://github.com/powerman/rpc-codec/releases) or run this to
install the latest version:

```sh
curl -sfL $(curl -s https://api.github.com/repos/powerman/rpc-codec/releases/latest | grep -i /jsonrpc2client-$(uname -s)-$(uname -m)\" | cut -d\" -f4) | sudo install /dev/stdin /usr/local/bin/jsonrpc2client
```

## Usage

### jsonrpc2client

```
$ jsonrpc2client -h
Usage: jsonrpc2client [flags] method params-as-json
  -http.endpoint string
        service endpoint as url
  -notification
        send notification request
  -tcp.addr string
        service endpoint as host:port
  -transport string
        transport (stdin|tcp|http) (default "http")
  -version
        print version
$ jsonrpc2client -http.endpoint https://example.com/rpc method.name '{"namedArg1":"value"}'
$ jsonrpc2client -http.endpoint https://example.com/rpc method.name '["positionalArg1"]'
```
