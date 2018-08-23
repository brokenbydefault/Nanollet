package Packets


import (
	"testing"
	"net"
)

func BenchmarkKeepAlivePackage_Encode(b *testing.B) {
	server, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6000,
	})
	if err != nil {
		panic(err)
	}

	dial, err := 	net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6000,
	})
	if err != nil {
		panic(err)
	}

	packet := NewKeepAlivePackage()

	for i := 0; i < b.N; i++ {
		dial.Write(packet.Encode(ExtensionType(0)))
	}

	server.Close()
}