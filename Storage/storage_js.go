// +build js

package Storage

import (
	"encoding/gob"
	"github.com/brokenbydefault/Nanollet/Util"
	"bytes"
	"github.com/gopherjs/gopherjs/js"
	"github.com/Inkeliz/goco/nativestorage"
)

func init() {
	Engine = &JSStorage{}
	Engine.Load(&PersistentStorage)
}

type JSStorage struct{}

func (s *JSStorage) Save(from *Persistent) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(from); err != nil {
		panic(err)
	}

	// Prevent save binary directly on LocalStorage and so on.
	bHex := Util.SecureHexEncode(b.Bytes())

	switch {
	case js.Global.Get("NativeStorage") != js.Undefined:
		nativestorage.SetItem("storage.nanollet", bHex)
	case js.Global.Get("localStorage") != js.Undefined:
		js.Global.Get("localStorage").Call("setItem", "storage.nanollet", bHex)
	default:
		panic("Unsupported storage available")
	}
}

func (s *JSStorage) Load(to *Persistent) {
	var bHex string

	switch {
	case js.Global.Get("NativeStorage") != js.Undefined:
		bHex, _ = nativestorage.GetString("storage.nanollet")
	case js.Global.Get("localStorage") != js.Undefined:
		obj := js.Global.Get("localStorage").Call("getItem", "storage.nanollet")
		if obj != nil {
			bHex = obj.String()
		}
	default:
		panic("Unsupported storage available")
	}

	if bHex == "" {
		return
	}

	b, _ := Util.SecureHexDecode(bHex)

	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(to); err != nil {
		panic(err)
	}

}

func (s *JSStorage) Write(name string, data []byte) (path string, err error) {
	// No-op
	return "", nil
}

func (s *JSStorage) Read(name string) ([]byte, error) {
	// No-op
	return nil, nil
}
