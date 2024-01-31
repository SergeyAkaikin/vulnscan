package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/host_scan_app/params"
	"github.com/SergeyAkaikin/vulnscan/internal/pkg/host_scan_app"
	"time"
)

func main() {

	start := time.Now()

	parameters := params.Define()
	app := host_scan_app.New(parameters)
	app.Run()

	end := time.Since(start)
	fmt.Println("======", end, "=======")
}
