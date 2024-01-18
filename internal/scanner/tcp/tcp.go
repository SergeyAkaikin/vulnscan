package tcp

import (
	"net"
)

func ConnectScan(tcpAddr *net.TCPAddr) (open bool) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err == nil {
		open = true
		conn.Close()
	}

	return
}

func SYNScan(addr *net.TCPAddr) (open bool) {
	return
}
