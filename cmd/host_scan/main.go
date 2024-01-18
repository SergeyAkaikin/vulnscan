package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/scanner/tcp"
	"net"
)

const test = "php.testsparker.com"
const port = "80"

func main() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", test+":"+port)
	fmt.Println(tcp.SYNScan(tcpAddr))
}
