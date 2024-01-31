package scanners

import (
	"github.com/SergeyAkaikin/vulnscan/internal/host_scan_app/resolver"
	"github.com/SergeyAkaikin/vulnscan/internal/host_scan_app/scanners/icmp"
	"github.com/SergeyAkaikin/vulnscan/internal/host_scan_app/scanners/tcp"
	"github.com/SergeyAkaikin/vulnscan/internal/user"
)

func Ping(address string) bool {
	ip, err := resolver.ResolveIPAddr(address)
	if err != nil {
		return false
	}

	if up, _ := icmp.Ping(ip); up {
		return true
	}

	tcpAddr80, err := resolver.ResolveTCPAddr(address, 80)
	if err != nil {
		return false
	}
	tcpAddr443, err := resolver.ResolveTCPAddr(address, 443)
	if err != nil {
		return false
	}

	if user.IsPrivileged() {
		if tcp80, _ := tcp.ACKPing(tcpAddr80); tcp80 {
			return true
		}
		if tcp443, _ := tcp.ACKPing(tcpAddr443); tcp443 {
			return true
		}
		return false
	}

	if tcp80, _ := tcp.ConnectScan(tcpAddr80, TIMEOUT); tcp80 {
		return true
	}

	if tcp443, _ := tcp.ConnectScan(tcpAddr443, TIMEOUT); tcp443 {
		return true
	}

	return false
}
