package Packets

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

func TestConfirmACKPackage_Decode(t *testing.T) {
	expected := Block.UniversalBlock{
		Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
		Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
		Previous:       Block.NewBlockHash(Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5")),
		Link:           Block.NewBlockHash(Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604")),
		Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
		DefaultBlock: Block.DefaultBlock{
			Signature: Wallet.NewSignature(Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05")),
			PoW:       Block.NewWork(Util.SecureHexMustDecode("B34A5DE4C3F98B14")),
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

	if expected.Hash() != pack.Hashes[0] {
		t.Error("invalid decode")
	}
}

func TestConfirmACKPackage_Decode_ByHash(t *testing.T) {
	udpMessage, _ := Util.SecureHexDecode("52430e0e07050001e09655b3019ea595c3d30cb951ba7f21d1d538983f1a6c739e664e5581b3b2b9320f05908eaefd5fbbe6f82dd38635e95e509f10ea38ce816befe8b8d2c052110d7407ed7fc360c0becef253fd999004604a3259241e423cfaa1163267ad3e05f4b3e14400000000275bbf469018adfa97e187f166e4f0f48096ad644a13676a2ae9b5308d200830c27017e56d8cc16389d14b6a9117d7944ceba968a7b9dfee884e5c6c893342fb1a572abf9f9a8029d876631743b6703f50b0c170bc862dadf2ac63bf3b712e4124f9fe5bb8d5fd1894145d92d2b6fb4c51388ae46b32162bd9d699b22ebd8beb74b0289f7514d004b2e02f454cc92be9550e8a99c333e052dd1109a9428e38d43d0c08ec599112b9b74ea83d5e3d6ff1284ae948a7c8537a11de686003566e5dccf651def3703fe6e96b89eb10032c5f994f96309b28f66d2628ba46c8810585d974c0494acc911e0ca870ad9d57cdff2df91ec55646f54ea3cef910d5864f87b8dd262a7ff07b358a3dd43d6a409097e608b68ae46870916bb59b86a41a0b485f678b757285ca90c1c7da9e979ce18d7d58bbb261e0613bf311808030e074b56987692c206d7667a0312ccc29fa4bcbaee13b1f570efbdbd39ce37c49891922779ef7bf69f83e699de3ac678e23efbda66a5530b28eb8e3b4b2690ae83ca5c0")

	header := new(Header)
	err := header.Decode(udpMessage)
	if err != nil {
		t.Error(err)
	}

	pack := new(ConfirmACKPackage)
	if err = pack.Decode(header, udpMessage[HeaderSize:]); err != nil {
		t.Error(err)
	}

	if len(pack.Hashes) != 12 {
		t.Error("wrong hashes")
	}
}

func TestConfirmACKPackage_Encode(t *testing.T) {
	tx := &Block.UniversalBlock{
		Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
		Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
		Previous:       Block.NewBlockHash(Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5")),
		Link:           Block.NewBlockHash(Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604")),
		Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
		DefaultBlock: Block.DefaultBlock{
			Signature: Wallet.NewSignature(Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05")),
			PoW:       Block.NewWork(Util.SecureHexMustDecode("B34A5DE4C3F98B14")),
		},
	}

	_, sk, _ := Wallet.GenerateRandomKeyPair()

	pack := NewConfirmACKPackage(&sk, tx)
	encoded := EncodePacketUDP(*NewHeader(), pack)

	header := new(Header)
	if err := header.Decode(encoded); err != nil {
		t.Error(err)
	}

	depack := new(ConfirmACKPackage)
	err := depack.Decode(header, encoded[HeaderSize:])

	if err != nil {
		t.Error(err)
	}

}

func TestConfirmACKPackage_Encode_ByHash(t *testing.T) {

	tx := []Block.Transaction{
		&Block.UniversalBlock{
			Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
			Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
			Previous:       Block.NewBlockHash(Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5")),
			Link:           Block.NewBlockHash(Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604")),
			Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
			DefaultBlock: Block.DefaultBlock{
				Signature: Wallet.NewSignature(Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05")),
				PoW:       Block.NewWork(Util.SecureHexMustDecode("B34A5DE4C3F98B14")),
			},
		},
		&Block.UniversalBlock{
			Account:        Wallet.Address("xrb_3dm8nb5xenwobdk3agbi3soyqpi9jzygwer377btgcn9f5jg7tfjpf7kwas5").MustGetPublicKey(),
			Representative: Wallet.Address("xrb_1beta4nkzb3g6b1a1qhae89earmz3gk3kfrp3f8hyztm8qyjkeyz9kajfutq").MustGetPublicKey(),
			Previous:       Block.NewBlockHash(Util.SecureHexMustDecode("770FFBF8756C883CF78189A1D7E5A1A013E0E11C41324E4215A8AFFC728813E5")),
			Link:           Block.NewBlockHash(Util.SecureHexMustDecode("82F94DE07379887A7B0C822E2F6F1FD7DEE1C5C3D6CB1BCE1947413949AEB604")),
			Balance:        Numbers.NewRawFromBytes(Util.SecureHexMustDecode("00007B426FAB61F00DE36398FF693D50")),
			DefaultBlock: Block.DefaultBlock{
				Signature: Wallet.NewSignature(Util.SecureHexMustDecode("3E59F0496431B3A740279EF0D294133C5EF788BFF0647CBC97CECD21F2746D553BEA1109FF3749E4D128503DD7E50EBAFAEA465DCD0EAD6C7F2871FFE988DC05")),
				PoW:       Block.NewWork(Util.SecureHexMustDecode("B34A5DE4C3F98B14")),
			},
		},
	}

	_, sk, _ := Wallet.GenerateRandomKeyPair()

	pack := NewConfirmACKPackage(&sk, tx...)
	encoded := EncodePacketUDP(*NewHeader(), pack)

	header := new(Header)
	if err := header.Decode(encoded); err != nil {
		t.Error(err)
	}

	depack := new(ConfirmACKPackage)
	err := depack.Decode(header, encoded[HeaderSize:])
	if err != nil {
		t.Error(err)
	}
}
