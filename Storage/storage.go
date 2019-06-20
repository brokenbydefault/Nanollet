// +build !js

package Storage

import (
	"github.com/shibukawa/configdir"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"encoding/gob"
	"github.com/brokenbydefault/Nanollet/Util"
	"bytes"
	"path/filepath"
)

func init() {
	Engine = &DesktopStorage{
		dir: configdir.New("BrokenByDefault", Configuration.Storage.Folder).QueryFolders(configdir.Global)[0],
	}

	Engine.Load(&PersistentStorage)

	if len(PersistentStorage.Quorum.PublicKeys) != 0 {
		Configuration.Account.Quorum = PersistentStorage.Quorum
	}
}

type DesktopStorage struct {
	dir *configdir.Config
}

func (s *DesktopStorage) Save(from *Persistent) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(from); err != nil {
		panic(err)
	}

	s.dir.WriteFile("storage.nanollet", b.Bytes())
}

func (s *DesktopStorage) Load(to *Persistent) {
	if s.dir.Exists("storage.nanollet") {
		r, err := s.dir.ReadFile("storage.nanollet")
		if err != nil {
			return
		}

		if err := gob.NewDecoder(bytes.NewReader(r)).Decode(to); err != nil {
			panic(err)
		}
	} else if s.dir.Exists("mfa.dat") || s.dir.Exists("wallet.dat") {
		s.loadLegacy(to)
		s.Save(to)
	}
}

func (s *DesktopStorage) Write(name string, data []byte) (path string, err error) {
	return filepath.Join(s.dir.Path, name), s.dir.WriteFile(name, data)
}

func (s *DesktopStorage) Read(name string) ([]byte, error) {
	return s.dir.ReadFile(name)
}

// Backward compatibility for Nanollet 1.0 and Nanollet 2.0
func (s *DesktopStorage) loadLegacy(to *Persistent) {
	if s.dir.Exists("mfa.dat") {
		if file, err := s.dir.ReadFile("mfa.dat"); err == nil {
			if b, ok := Util.SecureHexDecode(string(file)); ok && len(b) == 32 {
				to.AddAllowedKey(Wallet.NewPublicKey(b))
			}
		}
	}

	if s.dir.Exists("wallet.dat") {
		if data, err := s.dir.ReadFile("wallet.dat"); err == nil {
			if seedfy, err := Wallet.ReadSeedFY(string(data)); err == nil {
				to.SeedFY = seedfy
			}
		}
	}
}
