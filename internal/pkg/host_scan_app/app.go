package host_scan_app

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/host_scan_app/endpoint"
	"github.com/SergeyAkaikin/vulnscan/internal/host_scan_app/params"
)

type HostScanApp struct {
	params params.Parameters
}

func New(parameters params.Parameters) *HostScanApp {
	app := HostScanApp{}
	app.params = parameters
	return &app
}

func (a *HostScanApp) Run() {
	addresses := endpoint.InitAddresses(a.params.Host)
	upHosts := endpoint.PingAddresses(addresses)
	fmt.Println(upHosts, "are up")
	scanners := endpoint.InitScanners(a.params)
	report := endpoint.StartWorkers(upHosts, scanners)
	endpoint.WriteReport(report)
}
