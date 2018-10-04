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

type KeepAlivePackage struct {
	List []*Peer.Peer
}

func NewKeepAlivePackage(peer []*Peer.Peer) (packet *KeepAlivePackage) {
	return &KeepAlivePackage{
		List: peer,
	}
}

func (p *KeepAlivePackage) Encode(dst []byte) (n int, err error) {
	if p == nil {
		return
	}

	if len(dst) < KeepAlivePackageSize {
		return 0, ErrDestinationLenghtNotEnough
	}

	max := len(p.List)
	if len(p.List) > KeepAlivePackageNPeers {
		max = KeepAlivePackageNPeers
	}

	bi := 0
	for _, peer := range p.List[:max] {
		bi += copy(dst[bi:], peer.UDP.IP)
		bi += copy(dst[bi:], []byte{byte(peer.UDP.Port), byte(peer.UDP.Port << 8)})
	}

	return KeepAlivePackageSize, nil
}

func (p *KeepAlivePackage) Decode(_ *Header, src []byte) (err error) {
	if p == nil {
		return
	}

	// The packet should have at least 18 bytes and multiples by 18.
	l := len(src)
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
		dataPeer := src[bi:be]

		p.List = append(p.List, Peer.NewPeer(dataPeer[:net.IPv6len], int(dataPeer[16])|int(dataPeer[17])<<8))

		bi = be
	}

	return nil
}

func (p *KeepAlivePackage) ModifyHeader(h *Header) {
	h.MessageType = KeepAlive
}
