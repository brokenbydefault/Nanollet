package TwoFactor

import (
	"github.com/brokenbydefault/Nanollet/TwoFactor/Ephemeral"
	"github.com/aead/poly1305"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/Inkeliz/blakEd25519"
	"bytes"
	"encoding/binary"
	"github.com/aead/chacha20poly1305"
	"github.com/brokenbydefault/Nanollet/Util"
)

const (
	EnvelopeSize       = 1 + (Ephemeral.X25519Size * 2) + (blakEd25519.PublicKeySize * 2) + blakEd25519.SignatureSize + poly1305.TagSize
	EnvelopeSignedSize = EnvelopeSize - (blakEd25519.SignatureSize + poly1305.TagSize)
)

type Version uint8

type Envelope struct {
	Version  Version
	Sender   Ephemeral.PublicKey
	Receiver Ephemeral.PublicKey
	Capsule  Capsule
}

type Capsule struct {
	Device    Wallet.PublicKey
	Token     Token
	Signature Wallet.Signature

	Tag [poly1305.TagSize]byte
}

func NewEnvelope(device Wallet.PublicKey, sender Ephemeral.PublicKey, receiver Ephemeral.PublicKey, token Token) Envelope {
	return Envelope{
		Version:  1,
		Sender:   sender,
		Receiver: receiver,

		Capsule: Capsule{
			Device: device,
			Token:  token,
		},
	}
}

func (e *Envelope) Sign(device *Wallet.SecretKey) {
	msg, err := e.MarshalBinary()
	if err != nil {
		panic(err)
	}

	e.Capsule.Signature = device.MustSign(msg[:EnvelopeSignedSize])
}

func (e *Envelope) IsValidSignature(allowedDevices []Wallet.PublicKey) bool {
	if allowedDevices == nil {
		allowedDevices = append(allowedDevices, e.Capsule.Device)
	}

	msg, err := e.MarshalBinary()
	if err != nil {
		panic(err)
	}

	for _, device := range allowedDevices {
		if e.Capsule.Device == device {
			return device.IsValidSignature(msg[:EnvelopeSignedSize], &e.Capsule.Signature)
		}
	}

	return false
}

func (e *Envelope) Encrypt(sender *Ephemeral.SecretKey) error {
	key := sender.Exchange(e.Receiver)
	cipher, err := chacha20poly1305.NewXCipher(key[:])
	if err != nil {
		return err
	}

	nonce := Util.CreateHash(24, e.Sender[:], e.Receiver[:])
	msg, err := e.Capsule.MarshalBinary()
	if err != nil {
		return err
	}

	enc := cipher.Seal(msg[:0], nonce, msg[:len(msg)-poly1305.TagSize], nil)

	if err := binary.Read(bytes.NewReader(enc), binary.BigEndian, &e.Capsule); err != nil {
		return err
	}

	return nil
}

func (e *Envelope) Decrypt(receiver *Ephemeral.SecretKey) error {
	key := receiver.Exchange(e.Sender)
	cipher, err := chacha20poly1305.NewXCipher(key[:])
	if err != nil {
		return err
	}

	nonce := Util.CreateHash(24, e.Sender[:], e.Receiver[:])
	msg, err := e.Capsule.MarshalBinary()
	if err != nil {
		return err
	}

	if _, err := cipher.Open(msg[:0], nonce, msg, nil); err != nil {
		return err
	}

	if err := binary.Read(bytes.NewReader(msg), binary.BigEndian, &e.Capsule); err != nil {
		return err
	}

	return nil
}

func (e *Envelope) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, e); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *Envelope) UnmarshalBinary(b []byte) error {
	return binary.Read(bytes.NewReader(b), binary.BigEndian, e)
}

func (c *Capsule) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Capsule) UnmarshalBinary(b []byte) error {
	return binary.Read(bytes.NewReader(b), binary.BigEndian, c)
}
