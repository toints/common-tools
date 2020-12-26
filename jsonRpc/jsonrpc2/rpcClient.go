package jsonrpc2

import (
	"fmt"
	"github.com/tidwall/gjson"
	"net"
)

type Client struct {
	stream    ObjectStream
	mapMethod map[string]func(string, *Client)
	id        uint64
}

func NewClient(host string, codec ObjectCodec) (*Client, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	stream := NewBufferedStream(conn, codec)

	client := &Client{
		stream:    stream,
		mapMethod: make(map[string]func(string, *Client)),
	}
	go client.readMessages()
	return client, nil
}

func (c *Client) Method(method string, fn func(string, *Client)) *Client {
	c.mapMethod[method] = fn
	return c
}

func (c *Client) readMessages() {
	var err error
	for err == nil {
		var m anyMessage
		err = c.stream.ReadObject(&m)
		if err != nil {
			break
		}
		switch {
		case m.request != nil:
			fmt.Printf("%v\n", string(*m.request.Params))
			ret := gjson.Parse(string(*m.request.Params))
			method := ret.Get("method").String()
			if fn, ok := c.mapMethod[method];ok{
				fn(string(*m.response.Result), c)
			}
			continue
		case m.response != nil:
			fmt.Printf("%v\n", string(*m.response.Result))
			ret := gjson.Parse(string(*m.response.Result))
			method := ret.Get("method").String()
			if fn, ok := c.mapMethod[method];ok{
				fn(string(*m.response.Result), c)
			}
		}

	}
}

func (c *Client) Call(method string, params interface{}) error {
	c.id++
	sendMap := map[string]interface{}{
		"method": method,
		"params": params,
		"id":     c.id,
	}
	return c.stream.WriteObject(sendMap)
}