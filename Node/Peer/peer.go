package Peer

import (
	"net"
	"errors"
	"strconv"
	"time"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
	"crypto/rand"
)

const (
	BetaPort = 54000
	LivePort = 7075
)

const (
	Timeout = 2 * time.Minute
)

var (
	ErrInvalidIP      = errors.New("invalid ip")
	ErrInvalidPort    = errors.New("invalid port")
	ErrIncompleteData = errors.New("invalid IP:PORT, both are needed")
)

type Peer struct {
	IP   []byte
	Port int

	LastSeen  time.Time
	PublicKey Wallet.PublicKey
	Header    [8]byte
	Challenge [32]byte
}

func NewPeer(ip net.IP, port int) (peer *Peer) {
	peer = &Peer{
		IP:       ip,
		Port:     port,
		LastSeen: time.Now(),
	}

	rand.Read(peer.Challenge[:])

	return peer
}

func NewPeersFromString(hosts ...string) (peers []*Peer) {
	for _, host := range hosts {

		iph, porth, err := net.SplitHostPort(host)
		if err != nil {
			continue
		}

		port, err := strconv.Atoi(porth)
		if err != nil {
			continue
		}

		if ip := net.ParseIP(iph); ip != nil {
			peers = append(peers, NewPeer(ip, port))
		} else {
			ips, _ := Util.LookupIP(iph)
			for _, ip := range ips {
				peers = append(peers, NewPeer(ip, port))
			}
		}
	}

	return peers
}

func (p *Peer) IsActive() bool {
	if p == nil {
		return false
	}

	if time.Since(p.LastSeen) > Timeout {
		return false
	}

	return true
}

func (p *Peer) IsKnow() bool {
	if p == nil {
		return false
	}

	if Util.IsEmpty(p.PublicKey[:]) {
		return false
	}

	return true
}
