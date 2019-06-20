package Storage

import (
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

var Engine FileStorage
var PersistentStorage Persistent

type Persistent struct {
	SeedFY      Wallet.SeedFY
	AllowedKeys []Wallet.PublicKey
	Quorum      Peer.Quorum
}

func (p *Persistent) AddSeedFY(seedfy Wallet.SeedFY) {
	p.SeedFY = seedfy
	Engine.Save(p)
}

func (p *Persistent) AddAllowedKey(newKey Wallet.PublicKey) {
	for _, key := range p.AllowedKeys {
		if key == newKey {
			return
		}
	}

	if p.AllowedKeys == nil {
		p.AllowedKeys = []Wallet.PublicKey{newKey}
	} else {
		p.AllowedKeys = append(p.AllowedKeys, newKey)
	}

	Engine.Save(p)
}

type FileStorage interface {
	Save(*Persistent)
	Load(*Persistent)
	Write(name string, data []byte) (path string, err error)
	Read(name string) (data []byte, err error)
}
