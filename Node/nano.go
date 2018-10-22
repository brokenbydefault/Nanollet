// +build !js

package Node

import (
	"net"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"time"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Storage"
	"errors"
	"context"
	"sync"
)

var (
	ErrTCPNotAvailable = errors.New("this node doesn't support TCP")
)

type Server struct {
	udp *net.UDPConn
	//TCP *net.TCPListener

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
	if srv.udp == nil {
		srv.initUDP()
	}

	c := make(chan RawUDP, 1<<15)

	go func() {
		for {
			b := make([]byte, Packets.HeaderSize+Packets.MessageSize)
			end, dest, err := srv.udp.ReadFromUDP(b)
			if err != nil {
				continue
			}

			c <- RawUDP{
				Raw:    b[:end],
				Source: Peer.NewPeer(dest.IP, dest.Port),
			}
		}
	}()

	return c, nil
}

// SendUDPTo sends a package to specific peer, this peer must accept UDP.
func (srv *Server) SendUDPTo(packet Packets.PacketUDP, dest *Peer.Peer) (err error) {
	if dest == nil {
		return
	}

	if srv.udp == nil {
		srv.initUDP()
	}


	if _, err := srv.udp.WriteTo(Packets.EncodePacketUDP(srv.Header, packet), &net.UDPAddr{IP: dest.IP, Port: dest.Port}); err != nil {
		return err
	}

	return nil
}

// SendUDP sends a package to all known active peer.
func (srv *Server) SendUDP(packet Packets.PacketUDP) (err error) {
	for _, dest := range srv.Peer.GetAll() {
		srv.SendUDPTo(packet, dest)
	}

	// @TODO (inkeliz) add a error handler
	return nil
}

// SendTCPTo sends a package to specific peer, this peer must accept TCP.
func (srv *Server) SendTCPTo(request Packets.PacketTCP, responseType Packets.MessageType, dest *Peer.Peer) (packet Packets.PacketTCP, err error) {
	if dest == nil {
		return packet, ErrTCPNotAvailable
	}

	packet = Packets.DecodeResponsePacketTCP(responseType)

	tcp, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: dest.IP, Port: dest.Port})
	if err != nil {
		return packet, err
	}

	defer tcp.Close()

	tcp.SetReadDeadline(time.Now().Add(10 * time.Second))

	if err := Packets.EncodePacketTCP(srv.Header, request, tcp); err != nil {
		return packet, err
	}

	if err := packet.Decode(nil, tcp); err != nil {
		return packet, err
	}

	return packet, nil
}

// SendTCP sends a package to one random known peer.
func (srv *Server) SendTCP(request Packets.PacketTCP, responseType Packets.MessageType) (packets <-chan Packets.PacketTCP, cancel context.CancelFunc) {
	peers := srv.Peer.GetRandom(0)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	c := make(chan Packets.PacketTCP, len(peers))

	go func() {
		wg := new(sync.WaitGroup)

		for _, peer := range peers {
			if peer.IsActive() && peer.IsKnow() {
				wg.Add(1)
				go func(peer *Peer.Peer) {
					if packet, err := srv.sendTCPToContext(ctx, request, responseType, peer); err == nil && ctx.Err() == nil {
						c <- packet
					} else {
					}
					wg.Done()
				}(peer)
			}
		}

		wg.Wait()
		close(c)
	}()

	return c, cancel
}

func (srv *Server) sendTCPToContext(ctx context.Context, request Packets.PacketTCP, responseType Packets.MessageType, dest *Peer.Peer) (packet Packets.PacketTCP, err error) {
	dialer := new(net.Dialer)

	packet = Packets.DecodeResponsePacketTCP(responseType)

	addr := &net.TCPAddr{IP: dest.IP, Port: dest.Port}
	tcp, err := dialer.DialContext(ctx, "tcp", addr.String())
	if err != nil {
		return packet, err
	}

	tcp.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer tcp.Close()

	if err := Packets.EncodePacketTCP(srv.Header, request, tcp); err != nil {
		return packet, err
	}

	if err := packet.Decode(nil, tcp); err != nil {
		return packet, err
	}

	return packet, nil
}

func (srv *Server) initUDP() {
	udp, err := net.ListenUDP("udp", nil)
	if udp == nil || err != nil {
		panic(err)
	}

	srv.udp = udp
}
