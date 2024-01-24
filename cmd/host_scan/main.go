package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/resolver"
	"github.com/SergeyAkaikin/vulnscan/internal/scanner/tcp"
)

func main() {

	//host, port, enableList := params.Define()
	//
	//addresses := app.InitAddresses(host)
	//scanners := app.InitScanners(0, enableList)
	//report := app.StartWorkers(addresses, port, scanners)
	//app.WriteReport(report)

	tcpAddr, err := resolver.ResolveTCPAddr("php.testsparker.com", 80)
	fmt.Println(tcpAddr, err)
	fmt.Println(tcp.ACKPing(tcpAddr))
}
