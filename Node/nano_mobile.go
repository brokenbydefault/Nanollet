// +build js, android

package Node

import (
	"net"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Storage"
	"errors"
	"context"
)

var (
	ErrTCPNotAvailable = errors.New("this node doesn't support TCP")
)

type Server struct {
	Peer        *Storage.PeersBox
	Transaction *Storage.TransactionBox
	Header      Packets.Header
}

func NewServer(Header Packets.Header, Peer *Storage.PeersBox, Tx *Storage.TransactionBox) Node {
	return &Server{
		Peer:        Peer,
		Transaction: Tx,
		Header:      Header,
	}
}

func (srv *Server) Peers() *Storage.PeersBox {
	return srv.Peer
}

func (srv *Server) Transactions() *Storage.TransactionBox {
	return srv.Transaction
}

func (srv *Server) ListenTCP() (ch <-chan RawTCP, err error) {
	// no-op
	return nil, nil
}

func (srv *Server) ListenUDP() (ch <-chan RawUDP, err error) {
	c := make(chan RawUDP, 1<<15)
	return c, nil
}

// SendUDPTo sends a package to specific peer, this peer must accept UDP.
func (srv *Server) SendUDPTo(packet Packets.PacketUDP, dest *Peer.Peer) (err error) {

	return nil
}

// SendUDP sends a package to all known active peer.
func (srv *Server) SendUDP(packet Packets.PacketUDP) (err error) {
	return nil
}

// SendTCPTo sends a package to specific peer, this peer must accept TCP.
func (srv *Server) SendTCPTo(request Packets.PacketTCP, responseType Packets.MessageType, dest *Peer.Peer) (packet Packets.PacketTCP, err error) {
	if dest == nil {
		return packet, ErrTCPNotAvailable
	}

	return packet, nil
}

// SendTCP sends a package to one random known peer.
func (srv *Server) SendTCP(request Packets.PacketTCP, responseType Packets.MessageType) (packets <-chan Packets.PacketTCP, cancel context.CancelFunc) {

	return nil, cancel
}

func (srv *Server) sendTCPToContext(ctx context.Context, request Packets.PacketTCP, responseType Packets.MessageType, dest *Peer.Peer) (packet Packets.PacketTCP, err error) {
	return packet, nil
}

func (srv *Server) initUDP() {

}
