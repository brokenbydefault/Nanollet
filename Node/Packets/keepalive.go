package Packets

import (
	"net"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
)

const (
	PeerSize = 18

	KeepAlivePackageNPeers = 8
	KeepAlivePackageSize   = KeepAlivePackageNPeers * PeerSize
)

type KeepAlivePackage [KeepAlivePackageNPeers]*Peer.Peer

func NewKeepAlivePackage(peer []*Peer.Peer) (packet *KeepAlivePackage) {
	packet = new(KeepAlivePackage)

	// Maximum should be KeepAlivePackageNPeers or less than KeepAlivePackageNPeers, if not have enough peers.
	// It ignores other peers, keeping only the firsts.
	max := len(peer)
	if max > KeepAlivePackageNPeers {
		max = KeepAlivePackageNPeers
	}

	for i, peer := range peer[:max] {
		packet[i] = peer
	}

	return packet
}

func (p *KeepAlivePackage) Encode(lHeader *Header, rHeader *Header) (data []byte) {
	if p == nil {
		return
	}

	data = make([]byte, KeepAlivePackageSize)

	bi := 0
	for _, peer := range *p {
		bi += copy(data[bi:], peer.RawIP())
		bi += copy(data[bi:], peer.RawPort())
	}

	return data
}

func (p *KeepAlivePackage) Decode(rHeader *Header, data []byte) (err error) {
	if p == nil {
		return
	}

	// The packet should have at least 18 bytes and multiples by 18.
	l := len(data)
	if l%PeerSize != 0 {
		return ErrInvalidMessageSize
	}

	// Maximum should be KeepAlivePackageNPeers, or less than KeepAlivePackageNPeers if not have enough peers.
	// It ignores other peers, keeping only the firsts.
	max := l / PeerSize
	if KeepAlivePackageNPeers > max {
		max = KeepAlivePackageNPeers
	}

	bi := 0
	for i := 0; i < max; i++ {
		be := bi + PeerSize
		dataPeer := data[bi:be]

		p[i] = Peer.NewPeer(dataPeer[:net.IPv6len], int(data[16])|int(data[17])<<8)

		bi = be
	}

	return nil
}

func (p *KeepAlivePackage) ModifyHeader(h *Header) {
	h.MessageType = KeepAlive
}
