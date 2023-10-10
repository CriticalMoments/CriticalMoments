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
	sqldb        *sql.DB
}

func NewDB(dataDir string) (*DB, error) {
	if dirInfo, err := os.Stat(dataDir); err != nil || !dirInfo.IsDir() {
		return nil, errors.New("CriticalMoments: Data directory path does not exist")
	}

	dbPath := fmt.Sprintf("file:%s/critical_moments_db.db?_journal_mode=WAL&mode=rwc", dataDir)

	sqldb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &DB{
		databasePath: dbPath,
		sqldb:        sqldb,
	}, nil
}

func (db *DB) Close() error {
	return db.sqldb.Close()
}
