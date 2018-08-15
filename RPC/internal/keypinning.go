package internal

import (
	"crypto/subtle"
	"crypto/x509"
	"errors"
	"github.com/brokenbydefault/Nanollet/Config"
	"github.com/brokenbydefault/Nanollet/Util"
	"net/url"
)

var ErrInvalidCert = errors.New("invalid certificate")

func VerifyPeerCertificate(expectedhash []byte, uri string) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		// KeyPinning is ignored when debug is enabled
		if Config.Configuration().DebugStatus {
			return nil
		}

		link, err := url.Parse(uri)
		if err != nil {
			panic("the connection url is invalid")
		}
		for _, cert := range verifiedChains[0] {
			if cert.VerifyHostname(link.Hostname()) == nil {
				if subtle.ConstantTimeCompare(expectedhash, Util.CreateHash(64, cert.RawSubjectPublicKeyInfo)) == 1 {
					return nil
				}
			}
		}

		return ErrInvalidCert
	}
}
