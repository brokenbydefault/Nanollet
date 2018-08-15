// +build !js

package Connectivity

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brokenbydefault/Nanollet/Config"
	"github.com/brokenbydefault/Nanollet/RPC/internal"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"time"
)

type S struct {
	endpoint  string
	publickey []byte
	wss       *websocket.Conn
}

var Socket = NewSocket()
var ErrConnectionLost = errors.New("impossible to reach the wss server")

func NewSocket() *S {
	return &S{
		endpoint:  "wss://api.nanollet.org",
		publickey: []byte{0x26, 0x97, 0xab, 0x86, 0x23, 0xd3, 0xab, 0x86, 0xee, 0x4c, 0xe3, 0x05, 0xe3, 0x25, 0x9f, 0xfc, 0xef, 0x04, 0x4a, 0x41, 0xc3, 0x59, 0x27, 0x1c, 0x36, 0x59, 0x0f, 0x85, 0xaf, 0x95, 0xee, 0xf4, 0x9f, 0x08, 0x10, 0x80, 0x8f, 0xb0, 0x60, 0x5f, 0xad, 0x55, 0xa0, 0x56, 0x74, 0x12, 0x56, 0x63, 0xee, 0x3a, 0x5c, 0xda, 0x3e, 0xa5, 0xed, 0xee, 0x01, 0x7e, 0x30, 0xec, 0x6c, 0x20, 0xf9, 0x4f},
	}
}

func (c *S) StartWebsocket() (err error) {
	if c.wss != nil {
		c.wss.Close()
	}

	config, _ := websocket.NewConfig(c.endpoint, c.endpoint)
	config.Dialer = &net.Dialer{
		Timeout: 10 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial:     internal.DNSResolver("1.1.1.1:53"),
		},
	}

	config.TlsConfig = &tls.Config{
		InsecureSkipVerify:       Config.Configuration().DebugStatus,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.X25519},
		VerifyPeerCertificate:    internal.VerifyPeerCertificate(c.publickey, c.endpoint),
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}

	c.wss, err = websocket.DialConfig(config)
	return
}

func (c *S) CloseWebsocket() (err error) {
	return c.wss.Close()
}

// @TODO Improve reliability
func (c *S) ReceiveAllMessages(message []byte, response chan []byte) error {
	defer close(response)

	if message != nil {
		resp, err := c.SendRequest(message)
		if err == nil {
			response <- resp
		} else {
			return err
		}
	}

	for {
		var msg []byte
		if err := websocket.Message.Receive(c.wss, &msg); err != nil {
			return ErrConnectionLost
		}

		if Config.Configuration().DebugStatus {
			fmt.Println("Received at:")
			fmt.Println(string(msg))
		}

		response <- msg
	}
}

func (c *S) SendRequest(b []byte) (msg []byte, err error) {
	if c.wss == nil || !c.wss.IsClientConn() {
		if c.StartWebsocket() != nil {
			return nil, ErrConnectionLost
		}
	}

	if Config.Configuration().DebugStatus {
		fmt.Println("Requesting:")
		fmt.Println(string(b))
	}

	_, err = c.wss.Write(b)
	if err != nil {
		c.StartWebsocket()
		return
	}

	err = websocket.Message.Receive(c.wss, &msg)

	if Config.Configuration().DebugStatus {
		fmt.Println("Response:")
		fmt.Println(string(msg))
	}

	return
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
