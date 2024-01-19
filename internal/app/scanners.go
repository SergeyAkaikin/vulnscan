package app

import (
	"github.com/SergeyAkaikin/vulnscan/internal/params"
	"github.com/SergeyAkaikin/vulnscan/internal/scanner"
)

var scannersMapper = map[string]func() scanner.Scanner{
	"tc": scanner.NewTCPConnect,
	"ts": scanner.NewTCPSYN,
	"u":  scanner.NewUDP,
}

type ScannersPipeLine []scanner.Scanner

func InitScanners(timeout int, parameters params.EnableList) ScannersPipeLine {
	pipe := make(ScannersPipeLine, 0, len(scannersMapper))
	for key, isSet := range parameters {
		if isSet {
			pipe = append(pipe, scannersMapper[key]())
		}
	}

	return pipe
}
