package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/scanner/udp"
	"net"
	"time"
)

const test = "php.testsparker.com"
const port = "2611"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", test+":"+port)
	fmt.Println(udpAddr, err)

	open := udp.Scan(udpAddr, []byte{}, time.Second*1)
	fmt.Println(open)
}
