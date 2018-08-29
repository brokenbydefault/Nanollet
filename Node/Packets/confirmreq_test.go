package Packets

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Util"
	"bytes"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

func TestConfirmReqPackage_Decode(t *testing.T) {
	expected := Util.SecureHexMustDecode("79BEC4E064ED076700A0DDA218FFA7322E4EAE83A5FC503EC968FE929FFAFAAB")
	udpMessage := Util.SecureHexMustDecode("52430d0d070400067cb593df1e997e986bf2b9cf15f65ca48fa72b74ced2820b58f0d210f6ae4ae5dfe37b58b3ca19ce168d6a2d878ac539a0ea51e037687d2c6dfed18c22a7aad87cb593df1e997e986bf2b9cf15f65ca48fa72b74ced2820b58f0d210f6ae4ae5000000087a84c94598b4e5525f000000843481c75fa40484a4b203b290990c0bf916e5df918bc9f2952b1dd99a5928c3ed46be47a5dde7e22f1fd9e15eef5e5793dfc4015b65256a023d6fa2c08e60416fa72111428e17b41591a7413f9fbadfafcd86fafae3e9189fed7709f5facb0be541dc0e9b3041c6")

	header := new(Header)
	err := header.Decode(udpMessage)
	if err != nil {
		t.Error(err)
	}

	pack := new(ConfirmReqPackage)
	if err = pack.Decode(header, udpMessage[HeaderSize:]); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(pack.Transaction.Hash(), expected) {
		t.Error("decode error, invalid block")
	}
}

func TestConfirmReqPackage_Encode(t *testing.T) {
	expected := Util.SecureHexMustDecode("79BEC4E064ED076700A0DDA218FFA7322E4EAE83A5FC503EC968FE929FFAFAAB")
	tx := &Block.UniversalBlock{
		Account:        Wallet.Address("nano_1z7okhhjx8dym3oz7ggh4qu7sb6hnwoqbmpkia7ojw8k45ucwkq73irbtndz").MustGetPublicKey(),
		Representative: Wallet.Address("nano_1z7okhhjx8dym3oz7ggh4qu7sb6hnwoqbmpkia7ojw8k45ucwkq73irbtndz").MustGetPublicKey(),
		Previous:       Util.SecureHexMustDecode("DFE37B58B3CA19CE168D6A2D878AC539A0EA51E037687D2C6DFED18C22A7AAD8"),
		Link:           Util.SecureHexMustDecode("843481C75FA40484A4B203B290990C0BF916E5DF918BC9F2952B1DD99A5928C3"),
		Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("000000087A84C94598B4E5525F000000")),
		DefaultBlock: Block.DefaultBlock{
			Signature: Util.SecureHexMustDecode("ED46BE47A5DDE7E22F1FD9E15EEF5E5793DFC4015B65256A023D6FA2C08E60416FA72111428E17B41591A7413F9FBADFAFCD86FAFAE3E9189FED7709F5FACB0B"),
			PoW:       Util.SecureHexMustDecode("E541DC0E9B3041C6"),
		},
	}

	header := NewHeader()

	pack := NewConfirmReqPackage(tx)
	encoded := pack.Encode(header, nil)
	pack.ModifyHeader(header)

	depack := new(ConfirmReqPackage)
	if err := depack.Decode(header, encoded); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(depack.Transaction.Hash(), expected) {
		t.Error("encode error, invalid block")
	}
}
