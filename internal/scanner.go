package internal

import (
	"fmt"
	"net"
	"sync"
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

var timeout time.Duration = 500 * time.Millisecond

func (s *Scanner) ScanPort(hostname string, port int) (bool, error) {
	address := hostname + fmt.Sprintf(":%d", port)
	conn, err := net.DialTimeout(s.Protocol, address, s.timeout)

	if err != nil {
		return false, err
	}
	conn.Close()
	return true, nil
}

func (s *Scanner) ScanRange(hostname string, startPort, endPort int) {
	var wg sync.WaitGroup
	for i := startPort; i <= endPort; i++ {
		wg.Add(1)

		go func(p int) {
			defer wg.Done()
			ok, _ := s.ScanPort(hostname, p)
			if ok {
				fmt.Println("Port", p, "is open")
			}
		}(i)

	}

	wg.Wait()
}
