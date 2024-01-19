package scanner

import (
	"github.com/SergeyAkaikin/vulnscan/internal/resolver"
	"github.com/SergeyAkaikin/vulnscan/internal/scanner/tcp"
	"github.com/SergeyAkaikin/vulnscan/internal/scanner/udp"
	"log"
	"time"
)

var TIMEOUT time.Duration = 1 * time.Second

type Scanner interface {
	NetworkType() string
	Scan(ip string, port uint16) (open bool)
}

type Base struct {
	Network string
	Timeout time.Duration
}

type TCPConnect struct {
	Base
}

func NewTCPConnect() Scanner {
	sc := &TCPConnect{}
	sc.Timeout = TIMEOUT
	sc.Network = "tcp"
	return sc
}

func (sc *Base) NetworkType() string {
	return sc.Network
}

func (sc *TCPConnect) Scan(ip string, port uint16) bool {
	dstAddr, err := resolver.ResolveTCPAddr(ip, port)
	if err != nil {
		log.Println(err)
		return false
	}

	open, err := tcp.ConnectScan(dstAddr, sc.Timeout)

	if err != nil {
		log.Println(err)
	}

	return open
}

type TCPSYN struct {
	Base
}

func NewTCPSYN() Scanner {
	sc := &TCPSYN{}
	sc.Network = "tcp"
	sc.Timeout = TIMEOUT
	return sc
}

func (sc *TCPSYN) Scan(ip string, port uint16) bool {
	dstAddr, err := resolver.ResolveTCPAddr(ip, port)
	if err != nil {
		log.Println(err)
		return false
	}

	open, err := tcp.SYNScan(dstAddr, sc.Timeout)

	if err != nil {
		log.Println(err)
	}

	return open
}

type UDP struct {
	Base
	payload []byte
}

func NewUDP() Scanner {
	sc := &UDP{}
	sc.Network = "udp"
	sc.Timeout = TIMEOUT
	return sc
}

func (sc *UDP) Scan(ip string, port uint16) bool {
	dstAddr, err := resolver.ResolveUDPAddr(ip, port)
	if err != nil {
		log.Println(err)
		return false
	}

	open, err := udp.Scan(dstAddr, sc.Timeout, sc.payload)

	if err != nil {
		log.Println(err)
	}

	return open
}
