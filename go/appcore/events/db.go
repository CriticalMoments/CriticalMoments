package events

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	databasePath string
}

func NewDB(dataDir string) (*DB, error) {
	if dirInfo, err := os.Stat(dataDir); err != nil || !dirInfo.IsDir() {
		return nil, errors.New("CriticalMoments: Data directory path does not exist")
	}

	dbPath := fmt.Sprintf("%s/critical_moments_db.db", dataDir)
	return &DB{
		databasePath: dbPath,
	}, nil
}

func (db *DB) testDB() (int, error) {
	conn, err := sql.Open("sqlite3", db.databasePath)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	r, err := conn.Query("SELECT 99")
	if err != nil {
		return 0, err
	}
	defer r.Close()
	r.Next()
	var v int
	err = r.Scan(&v)
	if err != nil {
		return 0, err
	}

	return v, nil
}
