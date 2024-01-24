package icmp

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
)

func Ping(ipAddr *net.IPAddr) (bool, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "")
	if err != nil {
		return false, err
	}

	err = conn.SetDeadline(time.Now().Add(time.Second * 1))
	if err != nil {
		return false, err
	}

	defer conn.Close()

	mg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("PINGING"),
		},
	}

	mgb, err := mg.Marshal(nil)
	if err != nil {
		return false, err
	}

	if _, err := conn.WriteTo(mgb, ipAddr); err != nil {
		return false, err
	}

	rb := make([]byte, 1500)
	n, _, err := conn.ReadFrom(rb)
	if err != nil {
		return false, err
	}

	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])

	if err != nil {
		return false, err
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return true, nil
	default:
		return false, nil

	}

}
