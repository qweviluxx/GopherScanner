package repository

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) saveDB(ports []int, hostname string) error {
	jsonData, err := json.Marshal(ports)
	if err != nil {
		return err
	}

	jsonString := string(jsonData)

	insrt, err := r.db.Prepare("INSERT INTO history (target, ports) VALUES (?, ?)")
	if err != nil {
		return err
	}
	insrt.Exec(hostname, jsonString)
	insrt.Close()
	return nil
}

func New(path string) (Repository, error) {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return Repository{}, err
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS history (id INTEGER PPRIMARY KEY AUTOINCREMENT, target TEXT, ports TEXT, created_at DATETIME)")
	if err != nil {
		return Repository{}, err
	}
	statement.Exec()

	return Repository{db: database}, nil

}
