package jsonrpc2

import (
	"encoding/json"
	"net"
	"net/rpc"
)

var jErrRequest = json.RawMessage(`{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"invalid request"}}`) // Const.

// JSONRPC2 is an internal RPC service used to process batch requests.
type JSONRPC2 struct{}

// BatchArg is a param for internal RPC JSONRPC2.Batch.
type BatchArg struct {
	Ctx

	srv  *rpc.Server
	reqs []*json.RawMessage
}

// Batch is an internal RPC method used to process batch requests.
func (JSONRPC2) Batch(arg BatchArg, replies *[]*json.RawMessage) (err error) {
	cli, srv := net.Pipe()
	defer logIfFail(cli.Close)
	go arg.srv.ServeCodec(NewServerCodecContext(arg.Context(), srv, arg.srv))

	replyc := make(chan *json.RawMessage, len(arg.reqs))
	donec := make(chan struct{}, 1)

	go func() {
		*replies = collectReplies(json.NewDecoder(cli), replyc, len(arg.reqs))
		donec <- struct{}{}
	}()

	var testreq serverRequest
	for _, req := range arg.reqs {
		if req == nil || json.Unmarshal(*req, &testreq) != nil {
			replyc <- &jErrRequest
			continue
		}
		if testreq.ID != nil {
			replyc <- nil
		}
		_, err = cli.Write(append(*req, '\n'))
		if err != nil {
			break
		}
	}

	close(replyc)
	<-donec
	return err
}

// collectReplies returns replies received from replyc in the same order,
// with each nil reply replaced by the next reply decoded from dec.
func collectReplies(dec *json.Decoder, replyc <-chan *json.RawMessage, size int) []*json.RawMessage {
	replies := make([]*json.RawMessage, 0, size)
	for reply := range replyc {
		if reply == nil {
			reply = new(json.RawMessage)
			if dec.Decode(reply) != nil {
				reply = &jErrRequest
			}
		}
		replies = append(replies, reply)
	}
	return replies
}
