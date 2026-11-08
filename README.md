# JSON-RPC 2.0

[![License MIT](https://img.shields.io/badge/license-MIT-royalblue.svg)](LICENSE)
[![Go version](https://img.shields.io/github/go-mod/go-version/powerman/rpc-codec?color=blue)](https://go.dev/)
[![Test](https://img.shields.io/github/actions/workflow/status/powerman/rpc-codec/test.yml?label=test)](https://github.com/powerman/rpc-codec/actions/workflows/test.yml)
[![Coverage Status](https://raw.githubusercontent.com/powerman/rpc-codec/gh-badges/coverage.svg)](https://github.com/powerman/rpc-codec/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/powerman/rpc-codec?color=blue)](https://github.com/powerman/rpc-codec/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/powerman/rpc-codec.svg)](https://pkg.go.dev/github.com/powerman/rpc-codec)

![Linux | amd64 arm64 armv7 ppc64le s390x riscv64](https://img.shields.io/badge/Linux-amd64%20arm64%20armv7%20ppc64le%20s390x%20riscv64-royalblue)
![macOS | amd64 arm64](https://img.shields.io/badge/macOS-amd64%20arm64-royalblue)
![Windows | amd64 arm64](https://img.shields.io/badge/Windows-amd64%20arm64-royalblue)

jsonrpc2 is a codec for net/rpc.

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

```console
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
