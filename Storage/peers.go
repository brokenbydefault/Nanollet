package Storage

import (
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"net"
	"math/rand"
	"github.com/brokenbydefault/Nanollet/Util"
)

var (
	PeerStorage = PeersBox{
		List:      make(map[string]*Peer.Peer),
		Challenge: Peer.NewChallenge(),
	}
)

type PeersBox struct {
	List      map[string]*Peer.Peer
	Challenge Peer.Challenge
}

func (h *PeersBox) GetRandom(n int) (items []*Peer.Peer) {
	if h == nil {
		return
	}

	l := len(h.List)
	if l == 0 {
		return
	}

	var random = map[int]int{}
	for i := 0; i < n; i++ {
		n := rand.Intn(l - 1)
		random[n] = n
	}

	var i = 0
	for _, item := range h.List {
		if _, ok := random[i]; ok {
			items = append(items, item)
		}
		i++
	}

	return
}

func (h *PeersBox) Get(s string) (peer *Peer.Peer, ok bool) {
	if h == nil {
		return
	}

	peer, ok = h.List[s]
	return
}

func (h *PeersBox) IsAllowedIP(dest net.IP) bool {
	if h == nil {
		return false
	}

	peer, ok := h.List[dest.String()]
	if !ok {
		return false
	}

	pk := peer.PublicKey()
	if !peer.IsActive() || Util.IsEmpty(pk[:]) {
		return false
	}

	return true
}

func (h *PeersBox) Add(peers ...*Peer.Peer) {
	if h == nil {
		return
	}

	for _, peer := range peers {
		if len(h.List) < 64 {
			h.List[peer.RawIP().String()] = peer
		}
	}
}

func (h *PeersBox) Remove(peers ...*Peer.Peer) {
	if h == nil {
		return
	}

	for _, peer := range peers {
		delete(h.List, peer.RawIP().String())
	}
}
