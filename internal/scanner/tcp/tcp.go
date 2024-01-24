package tcp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"math/rand"
	"net"
	"time"
)

func ConnectScan(dstAddr *net.TCPAddr, timeout time.Duration) (bool, error) {
	conn, err := net.DialTimeout("tcp", dstAddr.String(), timeout)
	if err == nil {
		conn.Close()
		return true, nil
	}

	return false, err
}

func SYNScan(dstAddr *net.TCPAddr, timeout time.Duration) (bool, error) {
	dstIp := dstAddr.IP.To4()
	dstPort := layers.TCPPort(dstAddr.Port)

	srcIp, sPort, err := localIPPort(dstAddr)
	if err != nil {
		return false, err
	}
	srcPort := layers.TCPPort(sPort)

	ip := &layers.IPv4{
		SrcIP:    srcIp,
		DstIP:    dstIp,
		Protocol: layers.IPProtocolTCP,
	}

	tcp := &layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		Seq:     randomUint32(),
		SYN:     true,
		Window:  14600,
	}

	if err := tcp.SetNetworkLayerForChecksum(ip); err != nil {
		return false, err
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		return false, err
	}

	conn, err := net.ListenPacket("ip4:tcp", srcIp.String())
	if err != nil {
		return false, err
	}

	defer conn.Close()

	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIp}); err != nil {
		return false, err
	}

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return false, err
	}

	for {
		bf := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(bf)

		if err != nil {
			return false, err
		} else if addr.String() == dstIp.String() {
			packet := gopacket.NewPacket(bf[:n], layers.LayerTypeTCP, gopacket.Default)

			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == srcPort {
					if tcp.SYN && tcp.ACK {
						return true, nil
					}

				} else {
					return false, nil
				}
			}
		}
	}

}

func localIPPort(dstAddr *net.TCPAddr) (net.IP, int, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", dstAddr.String())
	if err != nil {
		return nil, 0, err
	}

	if conn, err := net.DialUDP("udp", nil, udpAddr); err == nil {
		if udpAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			return udpAddr.IP, udpAddr.Port, nil
		}
	}

	return nil, 0, err
}

func randomUint32() uint32 {
	s := rand.NewSource(time.Now().Unix())
	return rand.New(s).Uint32()
}
