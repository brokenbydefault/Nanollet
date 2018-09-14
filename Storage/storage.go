package Storage

import (
	"github.com/shibukawa/configdir"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"encoding/gob"
	"github.com/brokenbydefault/Nanollet/Util"
	"bytes"
)

var ArbitraryStorage = configdir.New("BrokenByDefault", Configuration.Storage.Folder).QueryFolders(configdir.Global)[0]

var PermanentStorage Persistent

type Persistent struct {
	SeedFY      Wallet.SeedFY
	AllowedKeys []Wallet.PublicKey
}

func init() {
	PermanentStorage.Load()
}

func (p *Persistent) Load() {
	if ArbitraryStorage.Exists("storage.nanollet") {
		r, err := ArbitraryStorage.ReadFile("storage.nanollet")
		if err != nil {
			return
		}

		if err := gob.NewDecoder(bytes.NewReader(r)).Decode(p); err != nil {
			panic(err)
		}
	} else if ArbitraryStorage.Exists("mfa.dat") || ArbitraryStorage.Exists("wallet.dat") {
		p.loadLegacy(ArbitraryStorage)
		p.Save()
	}
}

func (p *Persistent) Save() {

	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(p); err != nil {
		panic(err)
	}

	ArbitraryStorage.WriteFile("storage.nanollet", b.Bytes())
}

// Backward compatibility for Nanollet 1.0 and Nanollet 2.0
func (p *Persistent) loadLegacy(files *configdir.Config) {
	if files.Exists("mfa.dat") {
		if file, err := files.ReadFile("mfa.dat"); err == nil {
			if b, ok := Util.SecureHexDecode(string(file)); ok && len(b) == 32 {
				p.AddAllowedKey(Wallet.NewPublicKey(b))
			}
		}
	}

	if files.Exists("wallet.dat") {
		if data, err := files.ReadFile("wallet.dat"); err == nil {
			if seedfy, err := Wallet.ReadSeedFY(string(data)); err == nil {
				p.SeedFY = seedfy
			}
		}
	}
}

func (p *Persistent) AddSeedFY(seedfy Wallet.SeedFY) {
	p.SeedFY = seedfy
	p.Save()
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

	p.Save()
}
