package NanoAlias

import (
	"errors"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/Node"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"strings"
	"time"
)

type Address string

var (
	ErrNotFound            = errors.New("alias not found")
	ErrNotLocallyFound     = errors.New("alias not found on transaction history")
	ErrInvalidAliasBlock   = errors.New("invalid alias block")
	ErrInvalidAddr         = errors.New("invalid alias")
	ErrUnsupportedBlock    = errors.New("unsupported block")
	ErrUnconfirmedRegister = errors.New("alias block was not confirmed")
)

var (
	Representative = Wallet.Address("xrb_3this7is7an7a1ias33333333333333333333333333333333333x9bet1jy").MustGetPublicKey()
	Amount         = Numbers.NewRawFromBytes([]byte{0x01})
)

// GetPublicKey gets the Ed25519 public-key requesting OpenCAP server, the server should present in the address,
// returns the public-key. It's return an non-nil error if something bad happens.
func (addr Address) GetPublicKey() (pk Wallet.PublicKey, err error) {
	pkAlias, _, err := addr.GetAliasKey()
	if err != nil {
		return pk, err
	}

	chains, err := Node.GetMultiplesHistory(Background.Connection, &pkAlias, nil)
	if err != nil {
		return pk, ErrNotFound
	}

	opensHashes := make([]*Block.BlockHash, 0)
	for _, txs := range chains {
		hash := txs[len(txs)-1].Hash()

		opensHashes = append(opensHashes, &hash)
		Storage.TransactionStorage.Add(txs[len(txs)-1])
	}

	winner, ok := Storage.TransactionStorage.WaitConfirmation(&Storage.Configuration.Account.Quorum, 30*time.Second, opensHashes...)
	if !ok {
		return pk, ErrNotFound
	}

	txOpen, ok := Storage.TransactionStorage.GetByHash(winner)
	if !ok {
		return pk, ErrNotLocallyFound
	}

	if ok := IsValidOpenBlock(txOpen); !ok {
		return pk, ErrInvalidAliasBlock
	}

	_, source := txOpen.GetTarget()

	txSource, err := Node.GetBlock(Background.Connection, &source)
	if err != nil || txSource == nil {
		return pk, ErrInvalidAliasBlock
	}

	if txSource.GetType() != Block.State {
		return pk, ErrInvalidAliasBlock
	}

	return txSource.GetAccount(), nil
}

// MustGetPublicKey is a wrapper from GetPublicKey, which removes the error response and throws panic if error.
func (addr Address) MustGetPublicKey() Wallet.PublicKey {
	pk, err := addr.GetPublicKey()
	if err != nil {
		panic(err)
	}

	return pk
}

// IsValid returns true if the given encoded address have an correct formatted with the OpenCAP format.
func (addr Address) IsValid() bool {
	if rune(addr[0]) != '@' {
		return false
	}

	if len(addr) < 1 || len(addr) > 51 {
		return false
	}

	if _, ok := addr.GetAliasSecretKey(); !ok {
		return false
	}

	return true
}

// GetAliasSecretKey returns the compressed alias secret-key, which uses 5 bit per char.
// It uses a map of:
//
// | CHAR  |     ASCII      |    5BIT    |
// |-------|----------------|------------|
// | NULL  |  0  (00000000) | 0  (00000) |
// | -     | 45  (00101101) | 1  (00001) |
// | .     | 46  (00101110) | 2  (00010) |
// | _     | 95  (01011111) | 3  (00011) |
// | a     | 97  (01100001) | 4  (00100) |
// | b     | 98  (01100010) | 5  (00101) |
// | c     | 99  (01100011) | 6  (00110) |
// | d     | 100 (01100100) | 7  (00111) |
// | e     | 101 (01100101) | 8  (01000) |
// | f     | 102 (01100110) | 9  (01001) |
// | g     | 103 (01100111) | 10 (01010) |
// | h     | 104 (01101000) | 11 (01011) |
// | i     | 105 (01101001) | 12 (01100) |
// | j     | 106 (01101010) | 13 (01101) |
// | k     | 107 (01101011) | 14 (01110) |
// | l     | 108 (01101100) | 15 (01111) |
// | m     | 109 (01101101) | 16 (10000) |
// | n     | 110 (01101110) | 17 (10001) |
// | o     | 111 (01101111) | 18 (10010) |
// | p     | 112 (01110000) | 19 (10011) |
// | q     | 113 (01110001) | 20 (10100) |
// | r     | 114 (01110010) | 21 (10101) |
// | s     | 115 (01110011) | 22 (10110) |
// | t     | 116 (01110100) | 23 (10111) |
// | u     | 117 (01110101) | 24 (11000) |
// | v     | 118 (01110110) | 25 (11001) |
// | w     | 119 (01110111) | 26 (11010) |
// | x     | 120 (01111000) | 27 (11011) |
// | y     | 121 (01111001) | 28 (11100) |
// | z     | 122 (01111010) | 29 (11101) |
//
// This encoding is just to fit more char in the alias and remove similar char from ASCII.
func (addr Address) GetAliasSecretKey() (alias Wallet.SecretKey, ok bool) {
	if len(strings.TrimPrefix(string(addr), "@")) == 0 || len(addr) > 51 {
		return alias, false
	}

	bytePos, offset, e := 0, uint(0), int64(0)
	for _, v := range strings.TrimPrefix(string(addr), "@") {
		v -= (v - 30) >> 8 & 30
		// "-" to "." (45 to 46)   - [44]                     = 1 to 2
		v -= (((44 - v) >> 8) & ((v - 47) >> 8)) & 44
		// "_" (95)                                           = 3
		v -= (((94 - v) >> 8) & ((v - 66) >> 8)) & 92
		// "a" to "z" (97 to 122)                             = 4 to 29
		v -= (((96 - v) >> 8) & ((v - 123) >> 8)) & 93

		e |= ((int64(v-1) >> 8) | int64(29-v)>>8) & 1

		alias[bytePos] |= uint8(uint8(v << 3) >> offset)
		if offset > 3 {
			bytePos++
			alias[bytePos] |= uint8(uint8(v << 3) << (8 - offset))

			offset = 5 - (8 - offset)
		} else if offset == 3 {
			bytePos++
			offset = 0
		} else {
			offset += 5
		}
	}

	if e != 0 {
		return alias, false
	}

	return alias, true
}

func (addr Address) GetAliasKey() (pk Wallet.PublicKey, sk Wallet.SecretKey, err error) {
	sk, ok := addr.GetAliasSecretKey()
	if !ok {
		return pk, sk, ErrInvalidAddr
	}

	return Wallet.CreateKeyPair(sk[:32])
}
