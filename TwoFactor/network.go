// +build !js

package TwoFactor

import (
	"net"
	"encoding/binary"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/TwoFactor/Ephemeral"
)

func NewRequesterServer(sk *Ephemeral.SecretKey, allowedDevices []Wallet.PublicKey) (Request, <-chan *Envelope) {
	tcp, err := net.ListenTCP("tcp", &net.TCPAddr{IP: getIP()})
	if err != nil {
		panic(err)
	}

	c := make(chan *Envelope, 1)
	go listen(tcp, sk, allowedDevices, c)

	addr := tcp.Addr().(*net.TCPAddr)

	return NewRequester(sk.PublicKey(), addr.IP, addr.Port), c
}

func ReplyRequest(device *Wallet.SecretKey, token Token, request Request) error {
	sk := Ephemeral.NewEphemeral()

	envelope := NewEnvelope(device.PublicKey(), sk.PublicKey(), request.Receiver, token)
	envelope.Sign(device)
	envelope.Encrypt(&sk)

	tcp, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: request.IP[:], Port: int(request.Port)})
	if err != nil {
		return err
	}

	if err := binary.Write(tcp, binary.BigEndian, &envelope); err != nil {
		return err
	}

	return nil
}

func listen(tcp *net.TCPListener, sk *Ephemeral.SecretKey, devices []Wallet.PublicKey, c chan *Envelope) {
	for {
		con, err := tcp.AcceptTCP()
		if err != nil {
			continue
		}

		go readEnvelope(tcp, con, sk, devices, c)
	}

}

func readEnvelope(tcp *net.TCPListener, con *net.TCPConn, sk *Ephemeral.SecretKey, devices []Wallet.PublicKey, c chan *Envelope) {
	envelope := Envelope{}
	defer con.Close()

	if err := binary.Read(con, binary.BigEndian, &envelope); err != nil {
		return
	}

	if envelope.Receiver != sk.PublicKey() {
		return
	}

	if err := envelope.Decrypt(sk); err != nil {
		return
	}

	if !envelope.IsValidSignature(devices) {
		return
	}

	c <- &envelope
	tcp.Close()
}

func getIP() net.IP {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
