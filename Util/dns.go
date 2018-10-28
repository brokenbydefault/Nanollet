package Util

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"encoding/json"
	"strings"
	"errors"
	"time"
)

type rr int

func (rr rr) String() string {
	return strconv.Itoa(int(rr))
}

var (
	a    rr = 1
	aaaa rr = 28
	any  rr = 255
	srv  rr = 33
)

type response struct {
	Answer []struct {
		Type rr     `json:"type"`
		Data string `json:"data"`
	} `json:"Answer"`
}

func LookupIP(name string) (ips [][]byte, err error) {
	ips = make([][]byte, 0)

	ipv4, err := lookup(a, name)
	if err != nil {
		return nil, err
	}

	for _, ip := range ipv4.Answer {
		ips = append(ips, net.ParseIP(ip.Data))
	}

	ipv6, err := lookup(aaaa, name)
	if err != nil {
		return nil, err
	}

	for _, ip := range ipv6.Answer {
		ips = append(ips, net.ParseIP(ip.Data))
	}

	return ips, nil
}

func LookupSRV(service, proto, name string) (cname string, addrs []*net.SRV, err error) {
	raw, err := lookup(srv, "_"+service+"._"+proto+"."+name)

	for _, answer := range raw.Answer {
		sa := strings.Split(answer.Data, " ")
		if len(sa) != 4 {
			return cname, addrs, errors.New("invalid response")
		}

		addrs = append(addrs, &net.SRV{
			Target:   sa[3],
			Port:     mustInt16(sa[2]),
			Priority: mustInt16(sa[1]),
			Weight:   mustInt16(sa[0]),
		})
	}

	return cname, addrs, err
}

func lookup(rr rr, host string) (*response, error) {
	client := &http.Client{Timeout: 3 * time.Second}

	raw, err := client.Get("https://1.1.1.1/dns-query?" + url.Values{
		"name": []string{host},
		"type": []string{rr.String()},
		"ct":   []string{"application/dns-json"},
	}.Encode())

	if err != nil {
		return nil, err
	}

	defer raw.Body.Close()

	resp := new(response)
	if err := json.NewDecoder(raw.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func mustInt16(s string) uint16 {
	i, _ := strconv.Atoi(s)
	return uint16(i)
}