package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) SaveDB(ports []int, hostname string) error {
	jsonData, err := json.Marshal(ports)
	if err != nil {
		return err
	}

	jsonString := string(jsonData)

	insrt, err := r.db.Prepare("INSERT INTO history (target, ports) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer insrt.Close()
	_, err = insrt.Exec(hostname, jsonString)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Receiver() (string, error) {
	rows, err := r.db.Query("SELECT target, ports FROM history")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	responses := []ScanResponse{}
	var portsJSON string

	for rows.Next() {
		r := ScanResponse{}
		err := rows.Scan(&r.Hostname, &portsJSON)
		if err != nil {
			fmt.Println("DB Scan error:", err)
			continue
		}

		err = json.Unmarshal([]byte(portsJSON), &r.Ports)
		if err != nil {
			fmt.Println("JSON Unmarshall error for host", r.Hostname, ":", err)
			continue
		}

		responses = append(responses, r)
	}

	jsonData, err := json.Marshal(responses)
	if err != nil {
		return "", err
	}

	jsonString := string(jsonData)

	return jsonString, nil
}

func New(path string) (Repository, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return Repository{}, err
	}

	database.SetMaxOpenConns(1)
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY AUTOINCREMENT, target TEXT, ports TEXT, created_at DATETIME)")
	if err != nil {
		return Repository{}, err
	}
	statement.Exec()

	return Repository{db: database}, nil

}
