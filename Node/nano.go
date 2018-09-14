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

type Node interface {
	Start()

	SendUDPTo(packet Packets.PacketUDP, dest *Peer.Peer) (err error)
	SendUDP(packet Packets.PacketUDP) (err error)

	// SendTCPTo sends a package to specific peer, this peer must accept TCP.
	SendTCPTo(request Packets.PacketTCP, responseType Packets.MessageType, dest *Peer.Peer) (packet Packets.PacketTCP, err error)
	SendTCP(request Packets.PacketTCP, responseType Packets.MessageType) (packets <-chan Packets.PacketTCP, cancel context.CancelFunc)
}

type Server struct {
	udp *net.UDPConn
	//TCP *net.TCPListener
	Peers        *Storage.PeersBox
	Transactions *Storage.TransactionBox

	Header Packets.Header

	KeepAliveHandler  func(con *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte)
	HandshakeHandler  func(con *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte)
	PublishHandler    func(con *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte)
	ConfirmACKHandler func(con *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte)
	ConfirmReqHandler func(con *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte)
}

func (srv Server) Start() {
	if srv.KeepAliveHandler == nil {
		srv.KeepAliveHandler = defaultKeepAliveHandler
	}

	if srv.HandshakeHandler == nil {
		srv.HandshakeHandler = defaultHandshakeHandler
	}

	if srv.PublishHandler == nil {
		srv.HandshakeHandler = defaultPublishHandler
	}

	if srv.ConfirmACKHandler == nil {
		srv.ConfirmACKHandler = defaultConfirmACKHandler
	}

	if srv.ConfirmReqHandler == nil {
		srv.ConfirmReqHandler = defaultConfirmReqHandler
	}

	srv.startServer()
	srv.listenUDP()
	srv.keepAlive()
	srv.keepAskVotes()
}

// SendUDPTo sends a package to specific peer, this peer must accept UDP.
func (srv *Server) SendUDPTo(packet Packets.PacketUDP, dest *Peer.Peer) (err error) {
	if srv.udp == nil {
		srv.startServer()
	}

	if _, err := srv.udp.WriteTo(Packets.EncodePacketUDP(srv.Header, packet), dest.UDP); err != nil {
		return err
	}

	return nil
}

// SendUDP sends a package to all known active peer.
func (srv *Server) SendUDP(packet Packets.PacketUDP) (err error) {
	for _, dest := range srv.Peers.GetAll() {
		srv.SendUDPTo(packet, dest)
	}

	// @TODO (inkeliz) add a error handler
	return nil
}

// SendTCPTo sends a package to specific peer, this peer must accept TCP.
func (srv *Server) SendTCPTo(request Packets.PacketTCP, responseType Packets.MessageType, dest *Peer.Peer) (packet Packets.PacketTCP, err error) {
	if dest == nil || dest.UDP == nil {
		return packet, ErrTCPNotAvailable
	}

	packet = Packets.DecodeResponsePacketTCP(responseType)

	tcp, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: dest.UDP.IP, Port: dest.UDP.Port})
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
	peers := srv.Peers.GetRandom(0)

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

	tcp, err := dialer.DialContext(ctx, "tcp", dest.String())
	if err != nil {
		return packet, err
	}

	tcp.SetReadDeadline(time.Now().Add(2 * time.Second))
	defer tcp.Close()

	if err := Packets.EncodePacketTCP(srv.Header, request, tcp); err != nil {
		return packet, err
	}

	if err := packet.Decode(nil, tcp); err != nil {
		return packet, err
	}

	return packet, nil
}

func (srv *Server) keepAlive() {
	msg := func() {
		keepAlivePacket := Packets.NewKeepAlivePackage(srv.Peers.GetRandom(Packets.KeepAlivePackageNPeers))

		for _, peer := range srv.Peers.GetAll() {
			if !peer.IsActive() {
				srv.Peers.Remove(peer)
			}

			if !peer.IsKnow() {
				srv.SendUDPTo(Packets.NewHandshakePackage(peer.Challenge[:], nil), peer)
			}

			srv.SendUDPTo(keepAlivePacket, peer)
		}

	}

	go func() {
		for range time.Tick(15 * time.Second) {
			msg()
		}
	}()

	msg()
}

func (srv *Server) keepAskVotes() {
	go func() {
		for tx := range srv.Transactions.Listen() {

			// If already exist one block after this, we ignore.
			// That's because if the next block gets approved, this block is also valid.
			hash := tx.Hash()
			if _, ok := srv.Transactions.GetByPreviousHash(&hash); ok {
				continue
			}

			srv.SendUDP(Packets.NewConfirmReqPackage(tx.Transaction))
		}
	}()
}

func (srv *Server) listenUDP() {
	go func() {
		for {
			b := make([]byte, Packets.HeaderSize+Packets.MessageSize)
			h := new(Packets.Header)

			end, dest, err := srv.udp.ReadFromUDP(b)
			if err != nil {
				continue
			}

			go func() {
				if err := h.Decode(b[:Packets.HeaderSize]); err != nil {
					return
				}

				b = b[Packets.HeaderSize:end]
				if peer, ok := srv.Peers.Get(dest.IP); ok && peer.IsKnow() && peer.IsActive() {
					switch h.MessageType {
					case Packets.KeepAlive:
						srv.KeepAliveHandler(srv, dest, h, b)
					case Packets.Publish:
						srv.PublishHandler(srv, dest, h, b)
					case Packets.ConfirmReq:
						srv.ConfirmReqHandler(srv, dest, h, b)
					case Packets.ConfirmACK:
						srv.ConfirmACKHandler(srv, dest, h, b)
					case Packets.NodeHandshake:
						srv.HandshakeHandler(srv, dest, h, b)
					}
				} else {
					switch h.MessageType {
					case Packets.NodeHandshake:
						srv.HandshakeHandler(srv, dest, h, b)
					}
				}
			}()

		}
	}()
}

func (srv *Server) startServer() {
	udp, err := net.ListenUDP("udp", nil)
	if udp == nil || err != nil {
		panic(err)
	}

	srv.udp = udp
}
