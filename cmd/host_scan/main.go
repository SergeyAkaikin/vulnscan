package main

import (
	"github.com/SergeyAkaikin/vulnscan/internal/app"
	"github.com/SergeyAkaikin/vulnscan/internal/params"
)

func main() {

	host, port, enableList := params.Define()

	addresses := app.InitAddresses(host)
	scanners := app.InitScanners(0, enableList)
	report := app.StartWorkers(addresses, port, scanners)
	app.WriteReport(report)
	//addr, _ := resolver.ResolveTCPAddr("scanme.nmap.org", 31337)
	//fmt.Println(tcp.ConnectScan(addr, time.Second*2))
}
