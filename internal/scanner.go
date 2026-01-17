package internal

import (
	"fmt"
	"net"
	"time"
)

type Scanner struct {
	Protocol string
	timeout  time.Duration
}

func NewScanner(protocol string) *Scanner {
	return &Scanner{
		Protocol: protocol,
		timeout:  1 * time.Second,
	}
}

var timeout time.Duration = 1 * time.Second

func ScanPort(protocol, hostname string, port int) (bool, error) {
	address := hostname + fmt.Sprintf(":%d", port)
	conn, err := net.DialTimeout(protocol, address, timeout)

	if err != nil {
		fmt.Printf("TCP connection error: %v", err)
		return false, err
	}
	conn.Close()
	return true, nil
}
