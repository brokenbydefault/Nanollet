// Package OpenCAP is highly inspired in https://github.com/opencap/go-opencap. This packages follows the same
// construction of Wallet.
package OpenCAP

import (
	"strings"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"net/http"
	"encoding/json"
	"errors"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
)

type Address string

var (
	ErrInvalidAlias       = errors.New("invalid alias provided")
	ErrUnexpectedResponse = errors.New("bad data returned from alias host")
	ErrCAPPNotSupported   = errors.New("the CAPP is not currently supported")
	ErrNotFoundAlias      = errors.New("no host found for the given alias")
)

// GetPublicKey gets the Ed25519 public-key requesting OpenCAP server, the server should present in the address,
// returns the public-key. It's return an non-nil error if something bad happens.
func (addr Address) GetPublicKey() (pk Wallet.PublicKey, err error) {
	host := addr.GetHost()
	if !strings.HasSuffix(host, addr.GetHost()) {
		return pk, ErrCAPPNotSupported
	}

	client := &http.Client{Timeout: 3 * time.Second}

	respRaw, err := client.Get("https://" + addr.LookupHost() + "/v1/addresses?alias=" + addr.GetAlias() + "&address_type=300")
	if err != nil {
		return pk, ErrUnexpectedResponse
	}

	defer respRaw.Body.Close()

	resp := struct {
		AddressType int            `json:"address_type"`
		Address     Wallet.Address `json:"address"`
	}{}

	if err = json.NewDecoder(respRaw.Body).Decode(&resp); err != nil {
		return pk, ErrUnexpectedResponse
	}

	pk, err = resp.Address.GetPublicKey()
	if err != nil || resp.AddressType != 300 {
		return pk, ErrUnexpectedResponse
	}

	return pk, nil
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
	if len(addr) < 5 {
		// The address must be at least a$b.c
		return false
	}

	split := strings.Split(string(addr), "$")
	if len(split) != 2 || Util.HasEmptyString(split) {
		return false
	}

	domain := strings.Split(split[1], ".")
	if len(domain) < 2 || Util.HasEmptyString(domain) {
		return false
	}

	return true
}

// GetAlias extract the existing alias of the address, everything before $.
func (addr Address) GetAlias() string {
	return strings.SplitN(string(addr), "$", 1)[0]
}

// GetHost extract the existing alias of the address, everything after $.
func (addr Address) GetHost() string {
	split := strings.Split(string(addr), "$")
	if len(split) != 2 {
		return ""
	}

	return split[1]
}

// LookupHost lookup the host, by SRV DNS query, based the given host from the address.
func (addr Address) LookupHost() string {
	host := addr.GetHost()
	if host == "" {
		return ""
	}

	_, srvAddrs, err := Util.LookupSRV("opencap", "tcp", host)
	if err != nil || len(srvAddrs) < 1 {
		return ""
	}

	target := strings.TrimRight(srvAddrs[0].Target, ".")
	if len(target) == 0 {
		return ""
	}

	return target
}
