package params

import "flag"

type EnableList map[string]bool

func Define() (host string, port uint16, enableList EnableList) {
	enableList = make(EnableList, 4)
	tc := flag.Bool("tc", true, "Enables TCP connection scanning")
	ts := flag.Bool("ts", false, "Enables TCP SYN scanning (tcp connection never completes)")
	u := flag.Bool("u", false, "Enables UDP scanning")
	p := flag.Uint("p", 0, "Specifies port of target")
	flag.Parse()

	host = flag.Arg(0)

	port = uint16(*p)

	if *ts {
		*tc = false
	}
	enableList["tc"] = *tc
	enableList["ts"] = *ts
	enableList["u"] = *u
	return
}
