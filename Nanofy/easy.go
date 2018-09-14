package Nanofy

import (
	"io"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func VerifyFromHistory(file io.Reader, pk Wallet.PublicKey, txs []Block.Transaction) bool {
	if len(txs) < 3 {
		return false
	}

	fileHash, err := CreateHash(file)
	if err != nil {
		return false
	}

	for i, tx := range txs {
		dest, _ := tx.GetTarget()

		switch dest {
		case CreatePublicKey(0):
			if NewLegacyVerifierHash(fileHash, &pk, txs[i+1], tx).IsValid() {
				return true
			}
		case CreatePublicKey(1):
			if NewStateVerifierHash(fileHash, &pk, txs[i+2], txs[i+1], tx).IsValid() {
				return true
			}
		}
	}

	return false
}
