package params

import (
	"flag"
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/user"
	"os"
	"strconv"
	"strings"
	"time"
)

type EnableList map[string]bool
type Parameters struct {
	Host    string
	Ports   []uint16
	IsRange bool
	Enables EnableList
	Timeout time.Duration
}

func Define() Parameters {
	params := Parameters{Enables: make(EnableList, 4)}
	tc := flag.Bool("tc", true, "Enables TCP connection scanning")
	ts := flag.Bool("ts", false, "Enables TCP SYN scanning (tcp connection never completes)")
	u := flag.Bool("u", false, "Enables UDP scanning")
	p := flag.String("p", "", "Specifies port of target")
	timeout := flag.Duration("timeout", time.Second, "Sets timeout for requests")
	flag.Parse()

	if *ts {
		if !user.IsPrivileged() {
			fmt.Println("You must be privileged user for using TCP SYN scanning")
			os.Exit(2)
		}

		*tc = false

		if user.IsWindows() {
			*ts = false
			*tc = true
		}
	}

	params.Host = flag.Arg(0)
	params.Enables["tc"] = *tc
	params.Enables["ts"] = *ts
	params.Enables["u"] = *u

	ports, isRange := parsePorts(*p)
	params.IsRange = isRange
	params.Ports = ports

	params.Timeout = *timeout

	return params
}

func parsePorts(ports string) (portsList []uint16, isRange bool) {
	var portStrs []string

	if strings.Contains(ports, "-") {
		portStrs = strings.Split(ports, "-")
		if len(portStrs) != 2 {
			fmt.Println("Ports range should be in format (number)-(number)")
			os.Exit(2)
		}
		isRange = true
	} else {
		portStrs = strings.Split(ports, ",")

		if portStrs[0] == "" {
			return
		}
	}

	for _, portStr := range portStrs {
		port, err := strconv.ParseUint(portStr, 0, 16)
		if err != nil {
			fmt.Println("Ports should be positive numbers")
			os.Exit(2)
		}
		portsList = append(portsList, uint16(port))
	}

	return
}
