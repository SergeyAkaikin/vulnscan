package udp

import (
	"log/slog"
	"net"
	"time"
)

func Scan(udpAddr *net.UDPAddr, payload []byte, timeOut time.Duration) (open bool) {
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()
	if err != nil {
		return false
	}

	if err = sendRequest(conn, payload); err != nil {
		slog.Info("Sending udp request problem", "error", err)
		return
	}

	if err = waitResponse(conn, timeOut); err == nil {
		open = true
	}

	return
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
