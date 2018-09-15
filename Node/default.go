package Node

import (
	"net"
	"time"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Storage"
)

func NewServer(header Packets.Header) *Server {
	return &Server{
		Peers:        &Storage.PeerStorage,
		Transactions: &Storage.TransactionStorage,

		Header: header,

		KeepAliveHandler:  defaultKeepAliveHandler,
		HandshakeHandler:  defaultHandshakeHandler,
		PublishHandler:    defaultPublishHandler,
		ConfirmACKHandler: defaultConfirmACKHandler,
	}
}

func defaultKeepAliveHandler(srv *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.KeepAlivePackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if peer, ok := srv.Peers.Get(dest.IP); ok {
		peer.LastSeen = time.Now()
		rHeader.Encode(peer.Header[:])
	}

	for _, peer := range packet.List {
		if srv.Peers.Add(peer) == 1 {
			srv.SendUDPTo(Packets.NewHandshakePackage(peer.Challenge[:], nil), peer)
		}
	}

}

func defaultHandshakeHandler(srv *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.HandshakePackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if rHeader.ExtensionType.Is(Packets.Response) {
		if peer, ok := srv.Peers.Get(dest.IP); ok {
			if packet.PublicKey.IsValidSignature(peer.Challenge[:], &packet.Signature) {
				peer.PublicKey = packet.PublicKey
			}
		}
	}

	if rHeader.ExtensionType.Is(Packets.Challenge) {
		if peer, ok := srv.Peers.Get(dest.IP); !ok || !peer.IsKnow() {
			peer = Peer.NewPeer(dest.IP, dest.Port)

			if srv.Peers.Add(peer) == 1 {
				srv.SendUDPTo(Packets.NewHandshakePackage(peer.Challenge[:], packet.Challenge[:]), peer)
			}
		} else {
			srv.SendUDPTo(Packets.NewHandshakePackage(nil, packet.Challenge[:]), peer)
		}

	}
}

func defaultPublishHandler(srv *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.PushPackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	if !packet.Transaction.IsValidPOW() {
		return
	}

	srv.Transactions.Add(packet.Transaction)
}

func defaultConfirmACKHandler(srv *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.ConfirmACKPackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	for _, hash := range packet.Hashes {
		srv.Transactions.AddVotes(&hash, &packet.PublicKey, Util.BytesToUint(packet.Sequence[:], Util.BigEndian))
	}
}

func defaultConfirmReqHandler(srv *Server, dest *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	//@TODO create support to Confirm_Req
	//NO-OP
}
