package internal

import (
	"golang.org/x/net/context"
	"net"
)

func DNSResolver(ip string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return new(net.Dialer).DialContext(ctx, "udp", ip)
	}
}
