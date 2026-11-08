// Jsonrpc2client is a command-line client for JSON-RPC 2.0 services.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

const (
	transportSTDIN = "stdin"
	transportTCP   = "tcp"
	transportHTTP  = "http"
	indent         = "    "
	wantNArg       = 2 // Method and params.
)

//nolint:gochecknoglobals // Flags and build-time values.
var (
	cmd = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	ver string // set by ./release
	cfg struct {
		version      bool
		notification bool
		transport    string
		tcpAddr      string
		httpEndpoint string
	}
)

func main() {
	log.SetFlags(0)
	initFlags()
	flag.Parse()

	if cfg.version {
		fmt.Println(cmd, ver, runtime.Version())
		os.Exit(0)
	}
	method, params := parseArgs()

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	var result any
	if cfg.notification {
		err = client.Notify(method, params)
	} else {
		err = client.Call(method, params, &result)
	}

	if errors.Is(err, rpc.ErrShutdown) || errors.Is(err, io.ErrUnexpectedEOF) {
		log.Fatal(err)
	} else if err != nil {
		rpcerr := jsonrpc2.ServerError(err)
		rpcerrJSON, err := json.MarshalIndent(rpcerr, "", indent)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalf("%s", rpcerrJSON)
	}

	if !cfg.notification {
		resultJSON, err := json.MarshalIndent(result, "", indent)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultJSON)
	}
}

func initFlags() {
	flag.BoolVar(&cfg.version, "version", false, "print version")
	flag.BoolVar(&cfg.notification, "notification", false, "send notification request")
	flag.StringVar(&cfg.transport, "transport", transportHTTP, "transport (stdin|tcp|http)")
	flag.StringVar(&cfg.tcpAddr, "tcp.addr", "", "service endpoint as host:port")
	flag.StringVar(&cfg.httpEndpoint, "http.endpoint", "", "service endpoint as url")
	flag.Usage = func() { //nolint:reassign // Documented way to customize usage message.
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] method params-as-json\n", cmd)
		flag.PrintDefaults()
	}
}

// parseArgs returns method and params from command-line arguments.
// It exits on invalid flags or arguments.
func parseArgs() (method string, params any) {
	method, paramsJSON := flag.Arg(0), flag.Arg(1)
	err := json.Unmarshal([]byte(paramsJSON), &params)

	switch {
	case flag.NArg() != wantNArg:
		FatalUsagef("")
	case method == "":
		FatalUsagef("method: required\n")
	case err != nil:
		FatalUsagef("params: %s\n", err)
	case cfg.transport != transportSTDIN && cfg.transport != transportTCP && cfg.transport != transportHTTP:
		FatalFlagValue("must be one of: stdin, tcp, http", "transport", cfg.transport)
	case cfg.transport == transportHTTP && !Endpoint(&cfg.httpEndpoint):
		FatalFlagValue("must be endpoint", "http.endpoint", cfg.httpEndpoint)
	case cfg.transport == transportTCP && cfg.tcpAddr == "":
		FatalFlagValue("required", "tcp.addr", cfg.tcpAddr)
	}
	return method, params
}

func newClient() (*jsonrpc2.Client, error) {
	switch cfg.transport {
	case transportHTTP:
		return jsonrpc2.NewHTTPClient(cfg.httpEndpoint), nil
	case transportTCP:
		return jsonrpc2.Dial("tcp", cfg.tcpAddr)
	case transportSTDIN:
		return jsonrpc2.NewClient(os.Stdin), nil
	default:
		panic("never here")
	}
}
