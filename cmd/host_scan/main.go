package main

import (
	"fmt"
	"gitlab.simbirsoft/verify/s.akaykin/vulnscan/internal/scanner/udp"
	"net"
	"time"
)

const test = "scanme.nmap.org"
const port = "53"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", test+":"+port)
	fmt.Println(udpAddr, err)

	open := udp.Scan(udpAddr, []byte{}, time.Second*10)
	fmt.Println(open)
}
