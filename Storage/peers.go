package Storage

import (
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"net"
	"math/rand"
	"sync"
)

//@TODO (inkeliz) Replace map to sync.Map
var PeerStorage PeersBox

func init() {
	PeerStorage.Add(Configuration.Node.Peers...)
}

type PeersBox struct {
	list *sync.Map
}

func (h *PeersBox) Count() (len int) {
	if h == nil || h.list == nil {
		return 0
	}

	h.list.Range(func(_, _ interface{}) bool {
		len++
		return true
	})

	return len
}

func (h *PeersBox) GetAll() (items []*Peer.Peer) {
	if h == nil || h.list == nil {
		return
	}

	h.list.Range(func(key, value interface{}) bool {
		item, ok := value.(*Peer.Peer)
		if ok {
			items = append(items, item)
		}

		return true
	})

	return items
}

func (h *PeersBox) GetRandom(n int) (items []*Peer.Peer) {
	if h == nil || h.list == nil {
		return
	}

	list := h.GetAll()

	l := len(list)
	if l == 0 {
		return
	}

	if n <= 0 || n > l {
		n = l
	}

	var random = map[int]int{}
	for i := 0; i < n; i++ {
		n := rand.Intn(l)
		random[n] = n
	}

	for i := range random {
		items = append(items, list[i])
	}

	return
}

func (h *PeersBox) Get(ip net.IP) (peer *Peer.Peer, ok bool) {
	if h == nil || h.list == nil {
		return
	}

	val, ok := h.list.Load(string(ip))
	if !ok {
		return nil, false
	}

	peer, ok = val.(*Peer.Peer)

	return peer, ok
}

func (h *PeersBox) IsAllowedIP(ip net.IP) bool {
	if h == nil || h.list == nil {
		return false
	}

	peer, ok := h.Get(ip)
	if !ok {
		return false
	}

	if !peer.IsActive() || !peer.IsKnow() {
		return false
	}

	return true
}

func (h *PeersBox) Add(peers ...*Peer.Peer) (n int) {
	if h == nil {
		return
	}

	if h.list == nil {
		h.list = new(sync.Map)
	}

	for _, peer := range peers {
		if h.Count() < 256 {
			if _, old := h.list.LoadOrStore(string(peer.UDP.IP), peer); !old {
				n++
			}
		}
	}

	return n
}

func (h *PeersBox) Remove(peers ...*Peer.Peer) {
	if h == nil || h.list == nil {
		return
	}

	for _, peer := range peers {
		h.list.Delete(string(peer.UDP.IP))
	}
}
