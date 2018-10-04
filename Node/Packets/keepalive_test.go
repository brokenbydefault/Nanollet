package Packets

import (
	"testing"
	"net"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
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

	dst := make([]byte, PackageSize)

	packet := NewKeepAlivePackage(nil)
	packet.Encode(dst)

	for i := 0; i < b.N; i++ {
		dial.Write(dst)
	}

	server.Close()
}

func TestKeepAlivePackage_Decode(t *testing.T) {
	expected := []*Peer.Peer{
		Peer.NewPeer(net.ParseIP("87.65.160.15"), 54000),
		Peer.NewPeer(net.ParseIP("185.45.113.124"), 54000),
		Peer.NewPeer(net.ParseIP("178.128.149.150"), 1024),
		Peer.NewPeer(net.ParseIP("80.25.160.217"), 54000),
		Peer.NewPeer(net.ParseIP("90.229.199.116"), 54000),
		Peer.NewPeer(net.ParseIP("81.169.243.90"), 54000),
		Peer.NewPeer(net.ParseIP("77.20.254.59"), 24946),
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

	for i, peer := range pack.List {
		if !expected[i].UDP.IP.Equal(peer.UDP.IP) {
			t.Errorf("invalid decode, wrong ip. Gets %s expecting %s", peer.UDP.IP, expected[i].UDP.IP)
		}

		if expected[i].UDP.Port != peer.UDP.Port {
			t.Errorf("invalid decode, wrong port. Gets %d expecting %d", peer.UDP.Port, expected[i].UDP.Port)
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

	pack := NewKeepAlivePackage(peers)
	encoded := EncodePacketUDP(*NewHeader(), pack)

	header := new(Header)
	if err := header.Decode(encoded); err != nil {
		t.Error(err)
	}

	depack := new(KeepAlivePackage)
	if err := depack.Decode(header, encoded[HeaderSize:]); err != nil {
		t.Error(err)
	}

	for i, peer := range pack.List {
		if !peers[i].UDP.IP.Equal(peer.UDP.IP) {
			t.Errorf("invalid decode, wrong ip. Gets %s expecting %s", peer.UDP.IP, peers[i].UDP.IP)
		}

		if peers[i].UDP.Port != peer.UDP.Port {
			t.Errorf("invalid decode, wrong port. Gets %d expecting %d", peer.UDP.Port, peers[i].UDP.Port)
		}
	}

}
