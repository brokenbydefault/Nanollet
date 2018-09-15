package Util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	opencap "github.com/lane-c-wagner/go-opencap"
)

// LookupResponse represents the return value of an OpenCAP GET addresses
// request
type LookupResponse struct {
	AddressType int    `json:"address_type`
	Address     string `json:address`
}

// LookupAlias gets the address at the given alias
func LookupAlias(alias string) (string, error) {
	_, domain, err := opencap.ValidateAlias(alias)
	if err != nil {
		return "", errors.New("Invalid alias provided")
	}

	host, err := opencap.GetHost(domain)
	if err != nil {
		return "", errors.New("No host found for the given alias")
	}

	respRaw, err := http.Get("https://" + host + "/v1/addresses?alias=" + alias + "&address_type=300")
	if err != nil {
		return "", errors.New("No host found for the given alias")
	}

	defer respRaw.Body.Close()
	body, err := ioutil.ReadAll(respRaw.Body)
	resp := LookupResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", errors.New("Bad data returned from alias host")
	}
	return resp.Address, nil
}
