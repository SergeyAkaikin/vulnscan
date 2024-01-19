package resolver

import (
	"fmt"
	"log/slog"
	"net"
)

func Resolve(host string) (addresses []string) {
	addresses, err := net.LookupHost(host)
	if err != nil {
		slog.Info("Host address resolving problem: %v", "host", host, "error", err)
	}

	return
}

func ResolveDNS(address string) (addresses []string) {
	names, err := net.LookupAddr(address)
	if err != nil {
		slog.Info("Address names resolving problem", "address", address, "error", err)
		return
	}
	addresses = make([]string, 0, len(names))
	for _, name := range names {
		addresses = append(addresses, Resolve(name)...)
	}

	return
}

func ResolveTCPAddr(ip string, port uint16) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", ip, port))
}

func ResolveUDPAddr(ip string, port uint16) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", ip, port))
}
