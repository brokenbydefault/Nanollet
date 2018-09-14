package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"io"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Util"
)

type Legacy struct {
	fileHash       Wallet.PublicKey
	previousBlock  Block.Transaction
	signatureBlock Block.Transaction
	flagBlock      Block.Transaction
	secretKey      *Wallet.SecretKey
	publicKey      *Wallet.PublicKey
}

func NewLegacySigner(file io.Reader, sk *Wallet.SecretKey, previousBlock Block.Transaction) (Nanofier, error) {
	if sk == nil {
		return nil, ErrInsufficientData
	}

	hash, err := CreateHash(file)
	if err != nil {
		return nil, err
	}

	return NewLegacySignerHash(hash, sk, previousBlock), nil
}

func NewLegacySignerHash(fileHash Wallet.PublicKey, sk *Wallet.SecretKey, previousBlock Block.Transaction) Nanofier {
	pk := sk.PublicKey()

	return &Legacy{
		fileHash:      fileHash,
		previousBlock: previousBlock,
		secretKey:     sk,
		publicKey:     &pk,
	}
}

func NewLegacyVerifier(file io.Reader, pk *Wallet.PublicKey, signatureBlock Block.Transaction, flagBlock Block.Transaction) (Nanofier, error) {
	var hash Wallet.PublicKey
	if file != nil {
		h, err := CreateHash(file)
		if err != nil {
			return nil, err
		}

		hash = h
	}

	return NewLegacyVerifierHash(hash, pk, signatureBlock, flagBlock), nil
}

func NewLegacyVerifierHash(fileHash Wallet.PublicKey, pk *Wallet.PublicKey, signatureBlock Block.Transaction, flagBlock Block.Transaction) Nanofier {
	return &Legacy{
		fileHash:       fileHash,
		publicKey:      pk,
		signatureBlock: signatureBlock,
		flagBlock:      flagBlock,
	}
}

func (v *Legacy) NewSigner(file io.Reader, sk *Wallet.SecretKey, previousBlock Block.Transaction) {
	if v == nil {
		return
	}

	if nanofier, err := NewLegacySigner(file, sk, previousBlock); err != nil {
		*v = *nanofier.(*Legacy)
	}

	return
}

func (v *Legacy) NewVerifier(file io.Reader, pk *Wallet.PublicKey, signatureBlock Block.Transaction, flagBlock Block.Transaction) {
	if v == nil {
		return
	}

	if nanofier, err := NewLegacyVerifier(file, pk, signatureBlock, flagBlock); err != nil {
		*v = *nanofier.(*Legacy)
	}

	return
}

// CreateBlocks will return two blocks, the first with the address as the hash of the file and the second one is the
// flag block. It returns a non-nil error if something go wrong.
func (v *Legacy) CreateBlocks() (txs []Block.Transaction, err error) {
	if v == nil || v.secretKey == nil || v.previousBlock == nil {
		return txs, ErrInvalidBlockType
	}

	txs = make([]Block.Transaction, 2)

	txs[0], err = v.CreateSignatureBlock()
	if err != nil {
		return nil, err
	}

	txs[1], err = v.CreateFlagBlock()
	if err != nil {
		return nil, err
	}

	return txs, err
}

func (v *Legacy) CreateSignatureBlock() (tx Block.Transaction, err error) {
	if v == nil || v.previousBlock == nil {
		return nil, ErrInsufficientData
	}

	previous := v.previousBlock.SwitchToUniversalBlock(nil, nil)

	tx, err = Block.CreateUniversalSendBlock(v.secretKey, previous.Representative, previous.Balance, v.Amount(), previous.Hash(), v.fileHash)
	if err != nil {
		return nil, err
	}

	v.signatureBlock = tx
	return tx, nil
}

func (v *Legacy) CreateFlagBlock() (tx Block.Transaction, err error) {
	if v == nil || v.signatureBlock == nil {
		return nil, ErrInsufficientData
	}

	previous := v.signatureBlock.SwitchToUniversalBlock(nil, nil)

	tx, err = Block.CreateUniversalSendBlock(v.secretKey, previous.Representative, previous.Balance, v.Amount(), previous.Hash(), v.Flag())
	if err != nil {
		return nil, err
	}

	v.flagBlock = tx
	return tx, nil
}

// VerifySignature compares two blocks and the file, returns true if the combination of blocks, pk and file is valid.
func (v *Legacy) IsValid() (ok bool) {
	if v == nil {
		return false
	}

	if Util.IsEmpty(v.fileHash[:]) {
		return false
	}

	if !v.IsCorrectlyFormatted() {
		return false
	}

	if dest, _ := v.signatureBlock.GetTarget(); dest != v.fileHash {
		return false
	}

	return true
}

// VerifyBlock lacks the file verification, which can be used to verify if one block is one correct Nanofy transaction,
// in this case we don't care about the hash of file, because we can don't have the file itself.
func (v *Legacy) IsCorrectlyFormatted() (ok bool) {
	if v == nil || v.flagBlock == nil || v.signatureBlock == nil || Util.IsEmpty(v.publicKey[:]) {
		return false
	}

	// The destination of the flag MUST be the Version 0 address
	if destination, _ := v.flagBlock.GetTarget(); destination != v.Flag() {
		return false
	}

	// The previous MUST point to each other, "Flag" previous points to "Sig".
	if v.flagBlock.GetPrevious() != v.signatureBlock.Hash() {
		return false
	}

	// The type MUST be "send" (legacy one)
	if v.signatureBlock.GetType() != Block.Send || v.flagBlock.GetType() != Block.Send {
		return false
	}

	// The signature of the block itself MUST be a correct one
	if !isSignatureValid(v.publicKey, v.flagBlock) || !isSignatureValid(v.publicKey, v.signatureBlock) {
		return false
	}

	// The blocks MUST send only 1 raw.
	if v.signatureBlock.GetBalance().Subtract(v.flagBlock.GetBalance()).Compare(v.Amount()) != 0 {
		return false
	}

	return true
}

func (v *Legacy) Flag() Wallet.PublicKey {
	return CreatePublicKey(0)
}

func (v *Legacy) Amount() *Numbers.RawAmount {
	return Numbers.NewRawFromBytes([]byte{0x01})
}
