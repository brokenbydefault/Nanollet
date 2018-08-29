package Packets

import (
	"testing"
	"net"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"bytes"
)

func BenchmarkKeepAlivePackage_Encode(b *testing.B) {
	server, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6000,
	})
	if err != nil {
		panic(err)
	}

	dial, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6000,
	})
	if err != nil {
		panic(err)
	}

	packet := NewKeepAlivePackage(nil)

	for i := 0; i < b.N; i++ {
		dial.Write(packet.Encode(nil, nil))
	}

	server.Close()
}

func TestKeepAlivePackage_Decode(t *testing.T) {
	expected := []*Peer.Peer{
		Peer.NewPeer(net.ParseIP("87.65.160.15"), 54000),
		Peer.NewPeer(net.ParseIP("185.45.113.124"), 54000),
		Peer.NewPeer(net.ParseIP("178.128.149.150"), 54000),
		Peer.NewPeer(net.ParseIP("90.229.199.116"), 54000),
		Peer.NewPeer(net.ParseIP("81.169.243.90"), 54000),
		Peer.NewPeer(net.ParseIP("77.20.254.59"), 54000),
		Peer.NewPeer(net.ParseIP("13.59.162.102"), 54000),
	}

	udpMessage := Util.SecureHexMustDecode("52420d0d0702000000000000000000000000ffff5741a00ff0d200000000000000000000ffffb92d717cf0d200000000000000000000ffffb2809596000400000000000000000000ffff5019a0d9f0d200000000000000000000ffff5ae5c774f0d200000000000000000000ffff51a9f35af0d200000000000000000000ffff4d14fe3b726100000000000000000000ffff0d3ba266f0d2")

	header := new(Header)
	err := header.Decode(udpMessage)
	if err != nil {
		t.Error(err)
	}

	pack := new(KeepAlivePackage)
	if err = pack.Decode(header, udpMessage[HeaderSize:]); err != nil {
		t.Error(err)
	}

	for i, peer := range pack {
		if !expected[i].RawIP().Equal(peer.RawIP()) {
			t.Error("invalid decode, wrong ip")
		}

		if bytes.Equal(expected[i].RawPort(), peer.RawPort()) {
			t.Error("invalid decode, wrong port")
		}
	}

}

func TestKeepAlivePackage_Encode(t *testing.T) {
	peers := []*Peer.Peer{
		Peer.NewPeer(net.ParseIP("87.65.160.15"), 54000),
		Peer.NewPeer(net.ParseIP("185.45.113.124"), 54000),
		Peer.NewPeer(net.ParseIP("178.128.149.150"), 54000),
		Peer.NewPeer(net.ParseIP("90.229.199.116"), 54000),
		Peer.NewPeer(net.ParseIP("81.169.243.90"), 54000),
		Peer.NewPeer(net.ParseIP("77.20.254.59"), 54000),
		Peer.NewPeer(net.ParseIP("13.59.162.102"), 54000),
	}

	header := NewHeader()

	pack := NewKeepAlivePackage(peers)
	encoded := pack.Encode(header, nil)

	depack := new(KeepAlivePackage)
	if err := depack.Decode(header, encoded); err != nil {
		t.Error(err)
	}

}
