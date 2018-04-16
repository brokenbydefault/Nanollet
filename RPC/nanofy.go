package RPCClient

import (
	"github.com/brokenbydefault/Nanofy/nanofytypes"
	"github.com/brokenbydefault/Nanollet/RPC/rpctypes"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func GetBlockByFile(c rpctypes.Connection, filekey Wallet.PublicKey, pubkey Wallet.PublicKey) (resp nanofytypes.Response, err error) {
	req := nanofytypes.RequestByFile{
		FileKey: filekey,
		PubKey:  pubkey,
		DefaultRequest: nanofytypes.DefaultRequest{
			Action: "file",
			App:    "nanofy",
		},
	}

	err = c.SendRequestJSON(&req, &resp)
	return
}
