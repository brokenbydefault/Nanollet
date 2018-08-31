package Node

import (
	"net"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"github.com/brokenbydefault/Nanollet/Config"
	"time"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
	"fmt"
)

var Connection Server

type Server struct {
	UDP *net.UDPConn
	TCP *net.TCPListener
}

func init() {
	con, err := net.ListenUDP("udp", Config.Configuration().DefaultUDPNetwork)
	if err != nil {

	}

	Connection.UDP = con
	Connection.KeepAlive(&Storage.PeerStorage)
	Connection.KeepAskVotes(&Storage.TransactionStorage, &Storage.PeerStorage)
	Connection.ReadUDP(&Storage.PeerStorage)
}

func (srv Server) KeepAlive(peers *Storage.PeersBox) {
	go func() {
		for range time.Tick(15 * time.Second) {
			msg := Packets.NewKeepAlivePackage(peers.GetRandom(Packets.KeepAlivePackageNPeers))

			for _, peer := range peers.List {
				if !peer.IsActive() {
					peers.Remove(peer)
				}

				if pk := peer.PublicKey(); Util.IsEmpty(pk[:]) {
					srv.SendUDP(peer, nil, nil, msg)
				} else {
					srv.SendUDP(peer, nil, nil, Packets.NewHandshakePackage(peers.Challenge.Derivative(peer.RawIP()), nil))
				}
			}
		}
	}()
}

func (srv Server) KeepAskVotes(txs *Storage.TransactionsShelter, peers *Storage.PeersBox) {
	go func() {
		for tx := range txs.Unconfirmed.Listen() {

			// If already exist one block after this, we ignore.
			// That's because if the next block gets approved, this block is also valid.
			if _, ok := txs.Unconfirmed.GetByPreviousHash(tx.Hash()); ok {
				continue
			}

			srv.SendMultiByMapUDP(peers.List, nil, nil, Packets.NewConfirmReqPackage(tx))
		}
	}()
}

func (srv Server) SendUDP(dest *Peer.Peer, lHeader *Packets.Header, rHeader *Packets.Header, packet Packets.PacketUDP) (n int, err error) {
	return srv.UDP.WriteTo(Packets.EncodePacketUDP(lHeader, rHeader, packet), dest.UDPAddr())
}

func (srv Server) SendMultiUDP(dests []*Peer.Peer, lHeader *Packets.Header, rHeader *Packets.Header, packet Packets.PacketUDP) {
	for _, dest := range dests {
		srv.UDP.WriteTo(Packets.EncodePacketUDP(lHeader, rHeader, packet), dest.UDPAddr())
	}
}

func (srv Server) SendMultiByMapUDP(dests map[string]*Peer.Peer, lHeader *Packets.Header, rHeader *Packets.Header, packet Packets.PacketUDP) {
	for _, dest := range dests {
		srv.UDP.WriteTo(Packets.EncodePacketUDP(lHeader, rHeader, packet), dest.UDPAddr())
	}
}

func (srv Server) ReadUDP(peers *Storage.PeersBox) {
	go func() {
		for {
			b := make([]byte, Packets.HeaderSize+Packets.MessageSize)
			h := new(Packets.Header)

			end, dest, err := srv.UDP.ReadFromUDP(b)
			if err != nil {
				continue
			}

			go func() {
				if err := h.Decode(b[:end]); err != nil {
					return
				}

				b = b[Packets.HeaderSize:end]
				if peer, ok := peers.Get(dest.IP.String()); ok && peer.IsKnow() && peer.IsActive() {
					switch h.MessageType {
					case Packets.KeepAlive:
						KeepAliveHandler(dest, h, b, peers)
					case Packets.Publish:
						PublishHandler(dest, h, b, Storage.PK, Storage.TransactionStorage)
					case Packets.ConfirmReq:
						// NO-OP
					case Packets.ConfirmACK:
						ConfirmACKHandler(dest, h, b, Storage.TransactionStorage)
					case Packets.NodeHandshake:
						HandshakeHandler(dest, h, b, peers)
					}
				} else {
					switch h.MessageType {
					case Packets.NodeHandshake:
						HandshakeHandler(dest, h, b, peers)
					}
				}

			}()

		}
	}()
}

func KeepAliveHandler(dest *net.UDPAddr, rHeader *Packets.Header, msg []byte, peers *Storage.PeersBox) {
	packet := new(Packets.KeepAlivePackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if peer, ok := peers.Get(dest.IP.String()); ok {
		peer.SetLastSeen(time.Now())
	} else {
		peers.Add(packet.List...)
	}
}

func HandshakeHandler(dest *net.UDPAddr, rHeader *Packets.Header, msg []byte, peers *Storage.PeersBox) {
	packet := new(Packets.HandshakePackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	challenge := peers.Challenge.Derivative(dest.IP)

	if rHeader.ExtensionType.Is(Packets.Response) {
		if packet.PublicKey.IsValidSignature(challenge, packet.Signature) {
			peer, ok := peers.Get(dest.IP.String())
			if !ok {
				peer = Peer.NewPeer(dest.IP, dest.Port)
				peers.Add(peer)
			}

			peer.SetPublicKey(packet.PublicKey)
		}
	}

	if rHeader.ExtensionType.Is(Packets.Challenge) {
		var response *Packets.HandshakePackage
		if peer, ok := peers.Get(dest.IP.String()); !ok || !peer.IsKnow() {
			response = Packets.NewHandshakePackage(challenge, packet.Challenge[:])
		} else {
			response = Packets.NewHandshakePackage(nil, packet.Challenge[:])
		}

		Connection.SendUDP(Peer.NewPeer(dest.IP, dest.Port), nil, rHeader, response)
	}
}

func PublishHandler(dest *net.UDPAddr, rHeader *Packets.Header, msg []byte, pk Wallet.PublicKey, txs Storage.TransactionsShelter) {
	packet := new(Packets.PushPackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if dest, _ := packet.Transaction.GetTarget(); dest != pk {
		return
	}

	if _, _, ok := txs.GetByHash(packet.Transaction.Hash()); ok {
		return
	}

	if _, i, ok := txs.GetByPreviousHash(packet.Transaction.Hash()); ok && i != 0 {
		txs.Confirmed.Add(packet.Transaction)
		return
	}

	txs.Unconfirmed.Add(packet.Transaction)
}

func ConfirmACKHandler(dest *net.UDPAddr, rHeader *Packets.Header, msg []byte, txs Storage.TransactionsShelter) {
	packet := new(Packets.ConfirmACKPackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if packet.Transaction != nil {
		if _, _, ok := txs.GetByHash(packet.Transaction.Hash()); ok {
			txs.Votes.Add(packet.PublicKey, Util.BytesToUint(packet.Sequence[:], Util.BigEndian), packet.Transaction)
		}
	}

	if packet.Hashes != nil {
		for _, hash := range packet.Hashes {
			if tx, _, ok := txs.GetByHash(hash); ok {
				txs.Votes.Add(packet.PublicKey, Util.BytesToUint(packet.Sequence[:], Util.BigEndian), tx)
			}
		}
	}

}
