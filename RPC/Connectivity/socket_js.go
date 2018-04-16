// +build js

package Connectivity

import (
	"encoding/json"
	"errors"
	"github.com/gopherjs/websocket"
	"io"
	"net"
)

type S struct {
	endpoint string
	wss      net.Conn
}

var Socket = NewSocket()
var ErrConnectionLost = errors.New("impossible to reach the wss server")

func NewSocket() *S {
	return &S{
		endpoint: "wss://api.nanollet.org",
	}
}

func (c *S) StartWebsocket() error {
	wss, err := websocket.Dial(c.endpoint)
	if err != nil {
		wss = nil
		return ErrConnectionLost
	}

	c.wss = wss
	return nil
}

func (c *S) SendRequest(b []byte) (msg []byte, err error) {
	if c.wss == nil {
		if err = c.StartWebsocket(); err != nil {
			return
		}
	}

	_, err = c.wss.Write(b)
	if err != nil {
		c.wss = nil
		err = ErrConnectionLost
		return
	}

	msg = make([]byte, 2048)
	n, err := c.wss.Read(msg)
	if err != nil {
		c.wss = nil
	}

	return msg[:n], err
}

func (c *S) SendRequestReader(b []byte) (io.ReadCloser, error) {
	return nil, nil
}

func (c *S) SendRequestJSON(request interface{}, response interface{}, try ...interface{}) error {
	jsn, err := json.Marshal(request)
	if err != nil {
		return err
	}

	res, err := c.SendRequest(jsn)
	if err != nil {
		return err
	}

	if len(try) > 0 {
		json.Unmarshal(res, try[0])
	}

	return json.Unmarshal(res, response)
}
