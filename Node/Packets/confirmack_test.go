package Packets

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"bytes"
)

func TestConfirmACKPackage_Decode(t *testing.T) {
	tx := Block.UniversalBlock{
		Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
		Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
		Previous:       Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5"),
		Link:           Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604"),
		Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
		DefaultBlock: Block.DefaultBlock{
			Signature: Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05"),
			PoW:       Util.SecureHexMustDecode("B34A5DE4C3F98B14"),
		},
	}

	udpMessage, _ := Util.SecureHexDecode("52420d0d070500062a7a917a56320a4b04116aea8fdeb873b245d74d25b1069dd9edbe82d4b8f03c4b347753ca1dd761449db7adafff49eb95a5ce3f33c2f9e1677540876f0b666d6cb36c32ebe8a62ef8e4328b1d93a0d89150aca87a54f5a1dfd541e58d0411038064e90e00000000ae66a247d653954ae41439300e6bebda078ffcee33012953a72a8768e2e2e9b1770ffbf8756c883cf78189a1d7e5a1a013e0e11c41324e4215a8affc728813e5259a40a92fa42e2240805de8618ec4627f0ba41937160b4cff7f5335fd1933df00007b426fab61f00de36398ff693d5082f94de07379887a7b0c822e2f6f1fd7dee1c5c3d6cb1bce1947413949aeb6043e59f0496431b3a740279ef0d294133c5ef788bff0647cbc97cecd21f2746d553bea1109ff3749e4d128503dd7e50ebafaea465dcd0ead6c7f2871ffe988dc05b34a5de4c3f98b14")

	header := new(Header)
	err := header.Decode(udpMessage)
	if err != nil {
		t.Error(err)
	}

	pack := new(ConfirmACKPackage)
	if err = pack.Decode(header, udpMessage[HeaderSize:]); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(tx.Hash(), pack.Transaction.Hash()) {
		t.Error("invalid decode")
	}
}

func TestConfirmACKPackage_Encode(t *testing.T) {
	tx := &Block.UniversalBlock{
		Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
		Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
		Previous:       Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5"),
		Link:           Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604"),
		Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
		DefaultBlock: Block.DefaultBlock{
			Signature: Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05"),
			PoW:       Util.SecureHexMustDecode("B34A5DE4C3F98B14"),
		},
	}

	_, sk, _ := Wallet.GenerateRandomKeyPair()
	header := NewHeader()

	pack := NewConfirmACKPackage(sk, tx)
	encoded := pack.Encode(header, nil)
	pack.ModifyHeader(header)

	depack := new(ConfirmACKPackage)
	err := depack.Decode(header, encoded)
	if err != nil {
		t.Error(err)
	}

}

func TestConfirmACKPackage_Encode_ByHash(t *testing.T) {

	tx := []Block.Transaction{
		&Block.UniversalBlock{
			Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
			Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
			Previous:       Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5"),
			Link:           Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604"),
			Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
			DefaultBlock: Block.DefaultBlock{
				Signature: Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05"),
				PoW:       Util.SecureHexMustDecode("B34A5DE4C3F98B14"),
			},
		},
		&Block.UniversalBlock{
			Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
			Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
			Previous:       Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5"),
			Link:           Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604"),
			Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
			DefaultBlock: Block.DefaultBlock{
				Signature: Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05"),
				PoW:       Util.SecureHexMustDecode("B34A5DE4C3F98B14"),
			},
		},
	}

	_, sk, _ := Wallet.GenerateRandomKeyPair()
	header := NewHeader()

	pack := NewConfirmACKPackage(sk, tx...)
	encoded := pack.Encode(header, nil)
	pack.ModifyHeader(header)

	depack := new(ConfirmACKPackage)
	err := depack.Decode(header, encoded)
	if err != nil {
		t.Error(err)
	}
}
