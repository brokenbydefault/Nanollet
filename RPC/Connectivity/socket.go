package Connectivity

import (
	"golang.org/x/net/websocket"
	"encoding/json"
	"errors"
	"net"
	"time"
	"io"
	"crypto/tls"
	"github.com/brokenbydefault/Nanollet/RPC/PKP"
)

type S struct {
	endpoint  string
	publickey []byte
}

var Socket = S{
	endpoint:  "wss://api.nanollet.org",
	publickey: []byte{0x26, 0x97, 0xab, 0x86, 0x23, 0xd3, 0xab, 0x86, 0xee, 0x4c, 0xe3, 0x05, 0xe3, 0x25, 0x9f, 0xfc, 0xef, 0x04, 0x4a, 0x41, 0xc3, 0x59, 0x27, 0x1c, 0x36, 0x59, 0x0f, 0x85, 0xaf, 0x95, 0xee, 0xf4, 0x9f, 0x08, 0x10, 0x80, 0x8f, 0xb0, 0x60, 0x5f, 0xad, 0x55, 0xa0, 0x56, 0x74, 0x12, 0x56, 0x63, 0xee, 0x3a, 0x5c, 0xda, 0x3e, 0xa5, 0xed, 0xee, 0x01, 0x7e, 0x30, 0xec, 0x6c, 0x20, 0xf9, 0x4f},
}

var wss *websocket.Conn

var ErrConnectionLost = errors.New("impossible to reach the wss server")

func (c S) StartWebsocket() (err error) {
	if wss != nil {
		wss.Close()
	}

	config, _ := websocket.NewConfig(c.endpoint, "wss://api.nanollet.org")
	config.TlsConfig = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.X25519},
		VerifyPeerCertificate:    PKP.VerifyPeerCertificate(c.publickey, c.endpoint),
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}

	config.Dialer = &net.Dialer{
		Timeout: 5 * time.Second,
	}

	wss, err = websocket.DialConfig(config)
	return
}

func (c S) SendRequest(b []byte) (msg []byte, err error) {
	if wss == nil || !wss.IsClientConn() {
		if c.StartWebsocket() != nil {
			return nil, ErrConnectionLost
		}
	}

	_, err = wss.Write(b)
	if err != nil {
		c.StartWebsocket()
		return
	}

	err = websocket.Message.Receive(wss, &msg)
	return
}

func (c S) SendRequestReader(b []byte) (io.ReadCloser, error) {
	return nil, nil
}

func (c S) SendRequestJSON(request interface{}, response interface{}, try ...interface{}) error {
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
