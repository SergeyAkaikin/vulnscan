package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/app"
	"github.com/SergeyAkaikin/vulnscan/internal/params"
	"time"
)

func main() {

	start := time.Now()
	parameters := params.Define()
	addresses := app.InitAddresses(parameters.Host)
	upHosts := app.PingAddresses(addresses)
	scanners := app.InitScanners(parameters)
	report := app.StartWorkers(upHosts, scanners)
	end := time.Since(start)
	app.WriteReport(report)
	fmt.Println("======", end, "=======")
}
