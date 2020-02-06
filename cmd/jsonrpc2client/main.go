package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/seanhagen/rpc-codec/jsonrpc2"
)

const (
	transportSTDIN = "stdin"
	transportTCP   = "tcp"
	transportHTTP  = "http"
	indent         = "    "
)

//nolint:gochecknoglobals
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

	flag.BoolVar(&cfg.version, "version", false, "print version")
	flag.BoolVar(&cfg.notification, "notification", false, "send notification request")
	flag.StringVar(&cfg.transport, "transport", transportHTTP, "transport (stdin|tcp|http)")
	flag.StringVar(&cfg.tcpAddr, "tcp.addr", "", "service endpoint as host:port")
	flag.StringVar(&cfg.httpEndpoint, "http.endpoint", "", "service endpoint as url")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] method params-as-json\n", cmd)
		flag.PrintDefaults()
	}
	flag.Parse()

	method, paramsJSON := flag.Arg(0), flag.Arg(1)
	var params interface{}
	err := json.Unmarshal([]byte(paramsJSON), &params)

	switch {
	case cfg.version:
		fmt.Println(cmd, ver, runtime.Version())
		os.Exit(0)
	case len(flag.Args()) != 2:
		FatalUsage("")
	case method == "":
		FatalUsage("method: required\n")
	case err != nil:
		FatalUsage("params: %s\n", err)
	case !(cfg.transport == transportSTDIN || cfg.transport == transportTCP || cfg.transport == transportHTTP):
		FatalFlagValue("must be one of: stdin, tcp, http", "transport", cfg.transport)
	case cfg.transport == transportHTTP && !Endpoint(&cfg.httpEndpoint):
		FatalFlagValue("must be endpoint", "http.endpoint", cfg.httpEndpoint)
	case cfg.transport == transportTCP && cfg.tcpAddr == "":
		FatalFlagValue("required", "tcp.addr", cfg.tcpAddr)
	}

	var client *jsonrpc2.Client
	switch cfg.transport {
	case transportHTTP:
		client = jsonrpc2.NewHTTPClient(cfg.httpEndpoint)
	case transportTCP:
		client, err = jsonrpc2.Dial("tcp", cfg.tcpAddr)
		if err != nil {
			log.Fatal(err)
		}
	case transportSTDIN:
		client = jsonrpc2.NewClient(os.Stdin)
	default:
		panic("never here")
	}

	var result interface{}
	if cfg.notification {
		err = client.Notify(method, params)
	} else {
		err = client.Call(method, params, &result)
	}

	if err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
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
