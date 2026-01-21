package main

import (
	"fmt"
	"net/http"

	. "github.com/qweviluxx/GopherScanner.git/internal"
)

func handler(w http.ResponseWriter, r *http.Request) {
	scanner := NewScanner("tcp")
	ports := scanner.ScanRange("scanme.nmap.org", 20, 1000)
	fmt.Fprintf(w, "opened ports:%v", ports)

}

func main() {
	http.HandleFunc("/scan", handler)

	fmt.Println("Starting web-server on port 8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Starting server error:", err)
		return
	}
}
