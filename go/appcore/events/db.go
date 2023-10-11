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

	db := DB{
		databasePath: dbPath,
		sqldb:        sqldb,
	}

	err = db.migrate()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *DB) Close() error {
	return db.sqldb.Close()
}

// migrations can be run on each start because they are incremental and non-destructive
// Future migrations must also be incremental (only append to this, check if not exists)
func (db *DB) migrate() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_name TEXT NOT NULL,
			event_data TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		);

		CREATE TRIGGER IF NOT EXISTS insert_events_created_at 
		AFTER INSERT ON events
		BEGIN
			UPDATE events SET created_at =unixepoch('subsec') WHERE id = NEW.id;
		END;

		CREATE TRIGGER IF NOT EXISTS update_events_updated_at 
		AFTER UPDATE ON events
		BEGIN
			UPDATE events SET updated_at =unixepoch('subsec') WHERE id = NEW.id;
		END;
	`)
	if err != nil {
		return err
	}

	return nil
}
