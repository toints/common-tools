package jsonRpc

import (
	"testing"
	"time"

	"github.com/trying2016/common-tools/jsonRpc/jsonrpc2"
)

func TestRpcClient(t *testing.T) {
	client, err := jsonrpc2.NewClient("127.0.0.1:9017", VarintObjectCodec{})
	if err != nil {
		return
	}
	client.Call("login", "")
	for {
		client.Call("keepalived", "")
		time.Sleep(time.Second * 60)
	}
}
