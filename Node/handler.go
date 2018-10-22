package Node

import (
	"net"
	"time"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Storage"
	"context"
)

type RawUDP struct {
	Raw    []byte
	Source *net.UDPAddr
}

type RawTCP struct {
	Raw    []byte
	Source *net.TCPAddr
}

type Node interface {
	ListenUDP() (ch <-chan RawUDP, err error)
	ListenTCP() (ch <-chan RawTCP, err error)

	SendUDPTo(packet Packets.PacketUDP, dest *Peer.Peer) (err error)
	SendUDP(packet Packets.PacketUDP) (err error)

	SendTCPTo(request Packets.PacketTCP, responseType Packets.MessageType, dest *Peer.Peer) (packet Packets.PacketTCP, err error)
	SendTCP(request Packets.PacketTCP, responseType Packets.MessageType) (packets <-chan Packets.PacketTCP, cancel context.CancelFunc)

	Peers() *Storage.PeersBox
	Transactions() *Storage.TransactionBox
}

type HandlerFunc func(Node Node, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte)
type SenderFunc func(Node Node)

type HandlerGroup struct {
	Node Node

	KeepAliveHandler  HandlerFunc
	HandshakeHandler  HandlerFunc
	PublishHandler    HandlerFunc
	ConfirmACKHandler HandlerFunc
	ConfirmReqHandler HandlerFunc

	KeepAliveSender SenderFunc
	VoteSender      SenderFunc
}

func NewHandler(node Node) *HandlerGroup {
	return &HandlerGroup{
		Node: node,

		KeepAliveHandler:  defaultKeepAliveHandler,
		HandshakeHandler:  defaultHandshakeHandler,
		PublishHandler:    defaultPublishHandler,
		ConfirmACKHandler: defaultConfirmACKHandler,
		ConfirmReqHandler: defaultConfirmReqHandler,

		KeepAliveSender: defaultKeepAliveSender,
		VoteSender:      defaultVoteSender,
	}
}

// Start will listen the Node and use the defined Handlers on the HandlerGroup
func (hg *HandlerGroup) Start() error {
	if err := hg.listenUDP(); err != nil {
		return err
	}

	if err := hg.listenTCP(); err != nil {
		return err
	}

	hg.KeepAliveSender(hg.Node)
	hg.VoteSender(hg.Node)

	return nil
}

func (hg *HandlerGroup) listenTCP() error {
	// NO-OP
	return nil
}

func (hg *HandlerGroup) listenUDP() error {
	ch, err := hg.Node.ListenUDP()
	if err != nil {
		return err
	}

	go func() {
		for raw := range ch {
			r := raw
			go hg.readUDP(r)
		}
	}()

	return nil
}

func (hg *HandlerGroup) readUDP(raw RawUDP) {
	header := new(Packets.Header)

	if err := header.Decode(raw.Raw[:Packets.HeaderSize]); err != nil {
		return
	}

	raw.Raw = raw.Raw[Packets.HeaderSize:]
	if peer, ok := hg.Node.Peers().Get(raw.Source.IP); ok && peer.IsKnow() && peer.IsActive() {
		switch header.MessageType {
		case Packets.KeepAlive:
			hg.KeepAliveHandler(hg.Node, raw.Source, header,raw.Raw)
		case Packets.Publish:
			hg.PublishHandler(hg.Node, raw.Source, header,raw.Raw)
		case Packets.ConfirmReq:
			hg.ConfirmReqHandler(hg.Node, raw.Source, header,raw.Raw)
		case Packets.ConfirmACK:
			hg.ConfirmACKHandler(hg.Node, raw.Source, header,raw.Raw)
		case Packets.NodeHandshake:
			hg.HandshakeHandler(hg.Node, raw.Source, header,raw.Raw)
		}
	} else {
		switch header.MessageType {
		case Packets.Publish:
			hg.PublishHandler(hg.Node, raw.Source, header,raw.Raw)
		case Packets.ConfirmACK:
			hg.ConfirmACKHandler(hg.Node, raw.Source, header,raw.Raw)
		case Packets.NodeHandshake:
			hg.HandshakeHandler(hg.Node, raw.Source, header,raw.Raw)
		}
	}
}

func defaultKeepAliveHandler(node Node, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.KeepAlivePackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if peer, ok := node.Peers().Get(dest.IP); ok {
		peer.LastSeen = time.Now()
		rHeader.Encode(peer.Header[:])
	}

	for _, peer := range packet.List {
		if node.Peers().Add(peer) == 1 {
			node.SendUDPTo(Packets.NewHandshakePackage(peer.Challenge[:], nil), peer)
		}
	}

}

func defaultHandshakeHandler(node Node, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.HandshakePackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if rHeader.ExtensionType.Is(Packets.Response) {
		if peer, ok := node.Peers().Get(dest.IP); ok {
			if packet.PublicKey.IsValidSignature(peer.Challenge[:], &packet.Signature) {
				peer.PublicKey = packet.PublicKey
			}
		}
	}

	if rHeader.ExtensionType.Is(Packets.Challenge) {

		if peer, ok := node.Peers().Get(dest.IP); !ok || !peer.IsKnow() {
			peer = Peer.NewPeer(dest.IP, dest.Port)

			if node.Peers().Add(peer) == 1 {
				node.SendUDPTo(Packets.NewHandshakePackage(peer.Challenge[:], packet.Challenge[:]), peer)
			}
		} else {
			node.SendUDPTo(Packets.NewHandshakePackage(nil, packet.Challenge[:]), peer)
		}

	}
}

func defaultPublishHandler(node Node, _ *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.PushPackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if !packet.Transaction.IsValidPOW() {
		return
	}

	node.Transactions().Add(packet.Transaction)
}

func defaultConfirmACKHandler(node Node, _ *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.ConfirmACKPackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	for _, hash := range packet.Hashes {
		node.Transactions().AddVotes(&hash, &packet.PublicKey, Util.BytesToUint(packet.Sequence[:], Util.BigEndian))
	}
}

func defaultConfirmReqHandler(_ Node, _ *net.UDPAddr, _ *Packets.Header, _ []byte) {
	//@TODO create support to Confirm_Req
	//NO-OP
}

func defaultKeepAliveSender(node Node) {
	go func() {
		for range time.Tick(1 * time.Second) {

			keepAlivePacket := Packets.NewKeepAlivePackage(node.Peers().GetRandom(Packets.KeepAlivePackageNPeers))

			all := node.Peers().GetAll()
			for _, peer := range all {
				if !peer.IsActive() {
					node.Peers().Remove(peer)
				}

				if !peer.IsKnow() {
					node.SendUDPTo(Packets.NewHandshakePackage(peer.Challenge[:], nil), peer)
				}

				node.SendUDPTo(keepAlivePacket, peer)
			}
		}
	}()
}

func defaultVoteSender(node Node) {
	go func() {
		for tx := range node.Transactions().Listen() {

			// If already exist one block after this, we ignore.
			// That's because if the next block gets approved, this block is also valid.
			hash := tx.Hash()
			if _, ok := node.Transactions().GetByPreviousHash(&hash); ok {
				continue
			}

			node.SendUDP(Packets.NewConfirmReqPackage(tx.Transaction))
		}
	}()
}
