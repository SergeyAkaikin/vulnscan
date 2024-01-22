package resolver

import (
	"encoding/binary"
	"fmt"
	"net"
)

func Resolve(host string) (addresses []string) {
	addresses, _ = net.LookupHost(host)
	return
}
func ResolveHostName(address string) (hostnames []string) {
	hostnames, _ = net.LookupAddr(address)
	return
}

func ResolveDNS(address string) (addresses []string) {
	names, err := net.LookupAddr(address)
	if err != nil {
		return
	}

	addresses = make([]string, 0, len(names))
	for _, name := range names {
		addresses = append(addresses, Resolve(name)...)
	}

	return
}

func ResolveTCPAddr(ip string, port uint16) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
}

func ResolveUDPAddr(ip string, port uint16) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
}

func ResolveSubnetAddrs(addr string) (addresses []string) {
	_, ipv4Net, err := net.ParseCIDR(addr)
	if err != nil {
		return
	}

	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)
	finish := (start & mask) | (mask ^ 0xffffffff)
	addresses = make([]string, finish-start+1)
	for i, j := start, 0; i <= finish; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		addresses[j] = ip.String()
		j++
	}
	return
}
