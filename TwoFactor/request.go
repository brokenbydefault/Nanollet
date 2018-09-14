package TwoFactor

import (
	"github.com/brokenbydefault/Nanollet/TwoFactor/Ephemeral"
	"net"
	"github.com/skip2/go-qrcode"
	"image/color"
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/binary"
	"bytes"
)

type Request struct {
	Version  Version
	IP       [net.IPv6len]byte
	Port     int64
	Receiver Ephemeral.PublicKey
}

func NewRequester(pk Ephemeral.PublicKey, ip net.IP, port int) (requester Request) {
	return Request{
		Version:  1,
		IP:       toIPV6(ip),
		Port:     int64(port),
		Receiver: pk,
	}
}

func (r *Request) QRCode(size int, color color.Color) (png []byte, err error) {
	qr, err := qrcode.New(r.String(), qrcode.Highest)
	if err != nil {
		panic(err)
	}
	qr.BackgroundColor = color

	return qr.PNG(size)
}

func (r *Request) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Request) UnmarshalBinary(b []byte) (error) {
	if err := binary.Read(bytes.NewReader(b), binary.BigEndian, r); err != nil {
		return err
	}

	return nil
}

func (r *Request) String() string {
	b, _ := r.MarshalBinary()

	return Util.SecureHexEncode(b)
}

func toIPV6(ip net.IP) (ipV6 [net.IPv6len]byte) {
	copy(ipV6[:], ip.To16())
	return ipV6
}
