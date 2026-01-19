package main

import (
	"fmt"
	"time"

	"github.com/qweviluxx/GopherScanner.git/internal"
)

func main() {
	s := internal.NewScanner("tcp")
	conn, err := s.ScanPort("scanme.nmap.org", 80)
	fmt.Println(conn)
	fmt.Println(err)

	start := time.Now()
	s.ScanRange("scanme.nmap.org", 20, 100)
	duration := time.Since(start)
	fmt.Println(duration)
}
