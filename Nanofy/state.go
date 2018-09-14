package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"io"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"errors"
	"github.com/brokenbydefault/Nanollet/Util"
)

type State struct {
	fileHash       Wallet.PublicKey
	previousBlock  Block.Transaction
	signatureBlock Block.Transaction
	flagBlock      Block.Transaction
	secretKey      *Wallet.SecretKey
	publicKey      *Wallet.PublicKey
}

func NewStateSigner(file io.Reader, sk *Wallet.SecretKey, previousBlock Block.Transaction) (Nanofier, error) {
	if sk == nil {
		return nil, ErrInsufficientData
	}

	hash, err := CreateHash(file)
	if err != nil {
		return nil, err
	}

	return NewStateSignerHash(hash, sk, previousBlock), nil
}

func NewStateSignerHash(fileHash Wallet.PublicKey, sk *Wallet.SecretKey, previousBlock Block.Transaction) Nanofier {
	pk := sk.PublicKey()

	return &State{
		fileHash:      fileHash,
		previousBlock: previousBlock,
		publicKey:     &pk,
		secretKey:     sk,
	}

}
func NewStateVerifier(file io.Reader, pk *Wallet.PublicKey, previousBlock Block.Transaction, signatureBlock Block.Transaction, flagBlock Block.Transaction) (Nanofier, error) {
	if pk == nil {
		return nil, ErrInsufficientData
	}

	var hash Wallet.PublicKey
	if file != nil {
		h, err := CreateHash(file)
		if err != nil {
			return nil, err
		}

		hash = h
	}

	return NewStateVerifierHash(hash, pk, previousBlock, signatureBlock, flagBlock), nil
}

func NewStateVerifierHash(fileHash Wallet.PublicKey, pk *Wallet.PublicKey, previousBlock Block.Transaction, signatureBlock Block.Transaction, flagBlock Block.Transaction) Nanofier {
	return &State{
		fileHash:       fileHash,
		previousBlock:  previousBlock,
		signatureBlock: signatureBlock,
		flagBlock:      flagBlock,
		publicKey:      pk,
	}
}

var (
	ErrInvalidBlockType = errors.New("invalid block type")
	ErrInsufficientData = errors.New("insufficient data")
)

func (v *State) NewSigner(file io.Reader, sk *Wallet.SecretKey, previousBlock Block.Transaction) {
	if v == nil {
		return
	}

	if nanofier, err := NewStateSigner(file, sk, previousBlock); err != nil {
		*v = *nanofier.(*State)
	}

	return
}

func (v *State) NewVerifier(file io.Reader, pk *Wallet.PublicKey, previousBlock Block.Transaction, signatureBlock Block.Transaction, flagBlock Block.Transaction) {
	if v == nil {
		return
	}

	if nanofier, err := NewStateVerifier(file, pk, previousBlock, signatureBlock, flagBlock); err != nil {
		*v = *nanofier.(*State)
	}

	return
}

// CreateBlocks will return two blocks, the first with the address as the hash of the file and the second one is the
// flag block. It returns a non-nil error if something go wrong.
func (v *State) CreateBlocks() (txs []Block.Transaction, err error) {
	if v == nil || Util.IsEmpty(v.fileHash[:]) || v.secretKey == nil || v.previousBlock == nil {
		return txs, ErrInsufficientData
	}

	if v.previousBlock.GetType() != Block.State || !isSignatureValid(v.publicKey, v.previousBlock) {
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

func (v *State) CreateSignatureBlock() (tx Block.Transaction, err error) {
	if v == nil || Util.IsEmpty(v.fileHash[:]) || v.previousBlock == nil {
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

func (v *State) CreateFlagBlock() (tx Block.Transaction, err error) {
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
func (v *State) IsValid() (ok bool) {
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

// IsCorrectlyFormatted lacks the file verification. This function can be used to verify if one block is one correct
// Nanofy transaction, in this case we don't care about the file, because it's possible that the user don't know the
// file.

func (v *State) IsCorrectlyFormatted() (ok bool) {
	if v == nil || v.signatureBlock == nil || v.previousBlock == nil || (Util.IsEmpty(v.publicKey[:]) && Util.IsEmpty(v.secretKey[:])) {
		return false
	}

	// The destination of the flag MUST be the Version 1 address
	if destination, _ := v.flagBlock.GetTarget(); destination != v.Flag() {
		return false
	}

	// The previous MUST point to each other, "Flag" previous point to "Sig", "Sig" previous point to "Previous".
	if v.flagBlock.GetPrevious() != v.signatureBlock.Hash() || v.signatureBlock.GetPrevious() != v.previousBlock.Hash() {
		return false
	}

	// The type MUST be "state"
	if v.flagBlock.GetType() != Block.State || v.signatureBlock.GetType() != Block.State || v.previousBlock.GetType() != Block.State {
		return false
	}

	// The signature of the block itself MUST be a correct one
	if !isSignatureValid(v.publicKey, v.flagBlock) || !isSignatureValid(v.publicKey, v.signatureBlock) || !isSignatureValid(v.publicKey, v.previousBlock) {
		return false
	}

	// The blocks MUST send only 1 raw, it implicitly guaranties that is a send operation over state-block.
	if v.signatureBlock.GetBalance().Subtract(v.flagBlock.GetBalance()).Compare(v.Amount()) != 0 ||
		v.previousBlock.GetBalance().Subtract(v.signatureBlock.GetBalance()).Compare(v.Amount()) != 0 {
		return false
	}

	return true
}

func (v State) Flag() Wallet.PublicKey {
	return CreatePublicKey(1)
}

func (v State) Amount() *Numbers.RawAmount {
	return Numbers.NewRawFromBytes([]byte{0x01})
}
