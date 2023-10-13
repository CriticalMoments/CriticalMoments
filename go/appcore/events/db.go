package events

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
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
// Future migrations must also be incremental (only append to this, check if not exists),
// or we must implement a versioning system
func (db *DB) migrate() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY,
			event_name TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		);

		CREATE INDEX IF NOT EXISTS events_event_name_created_at ON events (event_name, created_at);

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

func (db *DB) InsertEvent(e *datamodel.Event) error {
	_, err := db.sqldb.Exec(`
		INSERT INTO events (event_name)
		VALUES (?)
	`, e.Name)
	if err != nil {
		return err
	}

	return nil
}

const eventCountByNameQuery = `SELECT COUNT(*) FROM events WHERE event_name = ?`

func (db *DB) EventCountByName(name string) (int, error) {
	r, err := db.sqldb.Query(eventCountByNameQuery, name)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	r.Next()
	var count int
	err = r.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const eventCountByNameWithLimitQuery = `SELECT COUNT(*) FROM (SELECT id FROM events WHERE event_name = ? LIMIT ?)`

func (db *DB) EventCountByNameWithLimit(name string, limit int) (int, error) {
	r, err := db.sqldb.Query(eventCountByNameWithLimitQuery, name, limit)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	r.Next()
	var count int
	err = r.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const latestEventTimeByNameQuery = `SELECT created_at FROM events WHERE event_name = ? ORDER BY created_at DESC LIMIT 1`

func (db *DB) LatestEventTimeByName(name string) (*time.Time, error) {
	r, err := db.sqldb.Query(latestEventTimeByNameQuery, name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	r.Next()
	var epochTime float64
	err = r.Scan(&epochTime)
	if err != nil {
		return nil, err
	}

	_, fractionalSeconds := math.Modf(epochTime)
	nanoseconds := int64(fractionalSeconds * 1_000_000_000)
	time := time.Unix(int64(epochTime), nanoseconds)
	return &time, nil
}
