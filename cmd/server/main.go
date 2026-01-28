package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qweviluxx/GopherScanner.git/internal"
	"github.com/qweviluxx/GopherScanner.git/internal/repository"
)

func validation(w http.ResponseWriter, h string, s, e int) bool {

	if h == "" || s < 0 || e < 0 || e < s || s > 65535 || e > 65535 {
		http.Error(w, "Bad params", http.StatusBadRequest)
		return false
	}

	return true
}

func handler(w http.ResponseWriter, r *http.Request, repo repository.Repository) {
	scanner := internal.NewScanner("tcp")

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()

	hostname := params.Get("hostname")

	startPort, err := strconv.Atoi(params.Get("startport"))
	if err != nil {
		http.Error(w, "Parsing param error:", http.StatusBadRequest)
		return
	}

	endPort, err := strconv.Atoi(params.Get("endport"))
	if err != nil {
		http.Error(w, "Parsing param error:", http.StatusBadRequest)
		return
	}

	valid := validation(w, hostname, startPort, endPort)
	if !valid {
		return
	}

	ctx := r.Context()
	ports := scanner.ScanRange(ctx, hostname, startPort, endPort)
	repo.SaveDB(ports, hostname)

	response := &repository.ScanResponse{Hostname: hostname, Ports: ports}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
		return
	}
}

func historyHandler(w http.ResponseWriter, r *http.Request, repo repository.Repository) {
	data, err := repo.Receiver()
	if err != nil {
		http.Error(w, "Failed to get history", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(data))
}

func main() {

	repo, err := repository.New("./scanner.db")
	if err != nil {
		fmt.Println("Failed to connect to DB:", err)
		return
	}

	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) { handler(w, r, repo) })
	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) { historyHandler(w, r, repo) })

	fmt.Println("Starting web-server on port 8080...")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Starting server error:", err)
		return
	}

}
