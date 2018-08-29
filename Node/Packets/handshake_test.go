package Packets

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Util"
	"bytes"
	"crypto/rand"
)

func TestHandshakePackage_Decode(t *testing.T) {
	expected := NewHandshakePackage(Util.SecureHexMustDecode("7929DF5BEBB4A10C2BA5C05D3851A1D989C4071618599DA1249D9D2CFE420BFB"), nil)
	udpMessage := Util.SecureHexMustDecode("52430d0d070a01007929df5bebb4a10c2ba5c05d3851a1d989c4071618599da1249d9d2cfe420bfb")

	header := new(Header)
	err := header.Decode(udpMessage)
	if err != nil {
		t.Error(err)
	}

	pack := new(HandshakePackage)
	if err = pack.Decode(header, udpMessage[HeaderSize:]); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(pack.Challenge[:], expected.Challenge[:]) {
		t.Error("error decode, invalid challenge")
	}

}

func TestHandshakePackage_Decode_3(t *testing.T) {
	expected := NewHandshakePackage(Util.SecureHexMustDecode("C47A5890B727C97323B0089E1F3D8CB559B4CB70107CA3CD0351C6B7B02A7E81"), nil)
	expected.PublicKey = Util.SecureHexMustDecode("8DEC937B67722A7DE508905674DAF6F5F7E99B676B304F0A0FC11966A08213CC")
	expected.Signature = Util.SecureHexMustDecode("C211B16511A9B5005CD484C9447DFA001F6F7916958D8C9348A49D124EAD144D8188EF8C5C70B967FF865CAB826EA7CC53E85C6CEA075258D3BDEE7B32EC7709")

	udpMessage := Util.SecureHexMustDecode("52430d0d070a0300c47a5890b727c97323b0089e1f3d8cb559b4cb70107ca3cd0351c6b7b02a7e818dec937b67722a7de508905674daf6f5f7e99b676b304f0a0fc11966a08213ccc211b16511a9b5005cd484c9447dfa001f6f7916958d8c9348a49d124ead144d8188ef8c5c70b967ff865cab826ea7cc53e85c6cea075258d3bdee7b32ec7709")

	header := new(Header)
	err := header.Decode(udpMessage)
	if err != nil {
		t.Error(err)
	}

	pack := new(HandshakePackage)
	if err = pack.Decode(header, udpMessage[HeaderSize:]); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(pack.Challenge[:], expected.Challenge[:]) {
		t.Error("error decode, invalid challenge")
	}

	if !bytes.Equal(pack.PublicKey, expected.PublicKey) {
		t.Error("error decode, invalid public-key")
	}

	if !bytes.Equal(pack.Signature, expected.Signature) {
		t.Error("error decode, invalid public-key")
	}

}

func TestHandshakePackage_Decode_2(t *testing.T) {
	expected := NewHandshakePackage(nil, nil)
	expected.PublicKey = Util.SecureHexMustDecode("64C49362A7B0101F0434EC6AB06C8A28C93DAF9974F32A88405E220894AE2164")
	expected.Signature = Util.SecureHexMustDecode("498C9B66C11E5CD9AE492A0DBADC633BE2A8F375E51E51D0A6727B47F17F4E7489545204798D03E6D257A2A94A3C71387EBAE1F48C59D817D3CEDE54298AF600")

	udpMessage := Util.SecureHexMustDecode("52430d0d070a020064c49362a7b0101f0434ec6ab06c8a28c93daf9974f32a88405e220894ae2164498c9b66c11e5cd9ae492a0dbadc633be2a8f375e51e51d0a6727b47f17f4e7489545204798d03e6d257a2a94a3c71387ebae1f48c59d817d3cede54298af600")

	header := new(Header)
	err := header.Decode(udpMessage)
	if err != nil {
		t.Error(err)
	}

	pack := new(HandshakePackage)
	if err = pack.Decode(header, udpMessage[HeaderSize:]); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(pack.Challenge[:], expected.Challenge[:]) {
		t.Error("error decode, invalid challenge")
	}

	if !bytes.Equal(pack.PublicKey, expected.PublicKey) {
		t.Error("error decode, invalid public-key")
	}

	if !bytes.Equal(pack.Signature, expected.Signature) {
		t.Error("error decode, invalid public-key")
	}
}

func TestHandshakePackage_Encode(t *testing.T) {
	challenge := make([]byte, 32)
	rand.Read(challenge)

	header := NewHeader()

	pack := NewHandshakePackage(challenge, nil)
	encoded := pack.Encode(header, nil)
	pack.ModifyHeader(header)

	depack := new(HandshakePackage)
	if err := depack.Decode(header, encoded); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(depack.Challenge[:], pack.Challenge[:]) {
		t.Error("error encode, invalid challenge")
	}
}

func TestHandshakePackage_Encode_3(t *testing.T) {
	challenge := make([]byte, 32)
	rand.Read(challenge)

	rchallenge := make([]byte, 32)
	rand.Read(challenge)

	header := NewHeader()

	pack := NewHandshakePackage(challenge, rchallenge)
	encoded := pack.Encode(header, nil)
	pack.ModifyHeader(header)

	depack := new(HandshakePackage)
	if err := depack.Decode(header, encoded); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(depack.Challenge[:], pack.Challenge[:]) {
		t.Error("error encode, invalid challenge")
	}

	if !pack.PublicKey.IsValidSignature(rchallenge, pack.Signature) {
		t.Error("error encode, invalid signature")
	}
}
