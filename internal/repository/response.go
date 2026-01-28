package repository

type ScanResponse struct {
	Hostname string `json:"hostname"`
	Ports    []int  `json:"ports"`
}
