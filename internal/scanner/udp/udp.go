package udp

import (
	"net"
	"time"
)

func Scan(dstAddr *net.UDPAddr, payload []byte, timeOut time.Duration) (bool, error) {
	conn, err := net.DialUDP("udp", nil, dstAddr)
	defer conn.Close()
	if err != nil {
		return false, err
	}

	if err = sendRequest(conn, payload); err != nil {
		return false, err
	}

	if err = waitResponse(conn, timeOut); err == nil {
		return true, nil
	}

	return false, err
}

func sendRequest(conn *net.UDPConn, payload []byte) error {
	_, err := conn.Write(payload)
	return err
}

func waitResponse(conn *net.UDPConn, timeOut time.Duration) error {
	buffer := make([]byte, 1024)
	if err := conn.SetDeadline(time.Now().Add(timeOut)); err != nil {
		return err
	}

	if _, err := conn.Read(buffer); err != nil {
		return err
	}

	return nil
}
