package app

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/params"
	"github.com/SergeyAkaikin/vulnscan/internal/resolver"
	"github.com/SergeyAkaikin/vulnscan/internal/scanner"
	"strings"
	"sync"
)

var scannersMapper = map[string]func() scanner.Scanner{
	"tc": scanner.NewTCPConnect,
	"ts": scanner.NewTCPSYN,
	"u":  scanner.NewUDP,
}

const (
	workersNum = 512
)

var (
	lPort   = uint16(1)
	rPort   = uint16(8192)
	isRange = true
	ports   []uint16
)

type TargetsReport map[string][]Target

type Target struct {
	Port    uint16
	Network string
}

type targetStatus struct {
	port    uint16
	ip      string
	network string
}
type ScannersPipeLine []scanner.Scanner

type openPorts map[string]struct{}

func InitScanners(parameters params.Parameters) ScannersPipeLine {
	pipeLine := make(ScannersPipeLine, 0, len(scannersMapper))
	for key, isSet := range parameters.Enables {
		if isSet {
			pipeLine = append(pipeLine, scannersMapper[key]())
		}
	}
	if parameters.IsRange {
		lPort = parameters.Ports[0]
		rPort = parameters.Ports[1]
	} else {
		isRange = false
		ports = parameters.Ports
	}

	scanner.TIMEOUT = parameters.Timeout

	return pipeLine
}

func InitAddresses(addressValue string) []string {

	if strings.ContainsRune(addressValue, '/') {
		addr := addressValue[:strings.Index(addressValue, "/")]

		if addrs := resolver.Resolve(addr); len(addrs) != 0 {
			mask := addressValue[strings.Index(addressValue, "/"):]
			return resolver.ResolveSubnetAddrs(fmt.Sprintf("%s%s", addrs[0], mask))
		}

		return resolver.ResolveSubnetAddrs(addressValue)
	}

	return resolver.ResolveDNS(addressValue)
}

func PingAddresses(addresses []string) []string {
	workersPool := make(chan struct{}, workersNum)
	defer close(workersPool)

	addressesCh := make(chan string)

	var wg sync.WaitGroup

	for _, address := range addresses {
		workersPool <- struct{}{}
		wg.Add(1)
		go pingWorker(address, workersPool, addressesCh, &wg)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(addressesCh)
	}(&wg)

	return writeAddresses(addressesCh)

}

func pingWorker(address string, workersPool <-chan struct{}, addressesCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	if scanner.Ping(address) {
		addressesCh <- address
	}
	<-workersPool
}

func writeAddresses(addressesCh <-chan string) []string {
	addressesList := make([]string, 0)
	for address := range addressesCh {
		addressesList = append(addressesList, address)
	}
	return addressesList
}

func StartWorkers(addrs []string, scannersPipeLine ScannersPipeLine) TargetsReport {
	workersPool := make(chan struct{}, workersNum)
	defer close(workersPool)

	targetsReport := make(TargetsReport)

	portsCh := make(chan targetStatus)
	defer close(portsCh)
	openPortsList := make(openPorts)

	var wg sync.WaitGroup

	go openPortsWriter(targetsReport, openPortsList, portsCh)

	for _, addr := range addrs {
		for _, currScanner := range scannersPipeLine {
			fmt.Println(addr, "scanning...")
			if isRange || (len(ports) == 0) {
				for port := lPort; port <= rPort; port++ {
					workersPool <- struct{}{}
					wg.Add(1)
					go portsWorker(
						addr,
						port,
						currScanner,
						workersPool,
						openPortsList,
						portsCh,
						&wg,
					)
				}
			} else {
				for _, port := range ports {
					workersPool <- struct{}{}
					wg.Add(1)
					go portsWorker(
						addr,
						port,
						currScanner,
						workersPool,
						openPortsList,
						portsCh,
						&wg,
					)
				}
			}
		}
	}

	wg.Wait()

	return targetsReport
}

func portsWorker(
	addr string,
	port uint16,
	scanner scanner.Scanner,
	workersPool <-chan struct{},
	openPortsList openPorts,
	openPortsCh chan<- targetStatus,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	ipport := fmt.Sprintf("%s:%d", addr, port)
	if _, scanned := openPortsList[ipport]; !scanned && scanner.Scan(addr, port) {
		openPortsCh <- targetStatus{port: port, ip: addr, network: scanner.NetworkType()}
	}

	<-workersPool
}

func openPortsWriter(
	report TargetsReport,
	portsList openPorts,
	portsCh <-chan targetStatus,
) {
	for target := range portsCh {
		key := fmt.Sprintf("%s:%d", target.ip, target.port)
		portsList[key] = struct{}{}
		if targetPortsList, exists := report[target.ip]; !exists {
			targetPortsList := make([]Target, 0)
			targetPortsList = append(targetPortsList, Target{target.port, target.network})
			report[target.ip] = targetPortsList

		} else {
			targetPortsList = append(targetPortsList, Target{target.port, target.network})
			report[target.ip] = targetPortsList

		}

	}
}

func WriteReport(report TargetsReport) {
	fmt.Println("================================================")
	for ip, ports := range report {

		host := ""
		if domains := resolver.ResolveHostName(ip); len(domains) != 0 {
			host = domains[0]
		}
		fmt.Printf("Scan report for %s (%s):\n", ip, host)
		for _, port := range ports {
			fmt.Printf("PORT: %d\\%s (open)\n", port.Port, port.Network)
		}
		fmt.Println("================================================")
	}
}
