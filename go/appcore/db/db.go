package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	databasePath string
	sqldb        *sql.DB
	started      bool

	eventManager           *EventManager
	propertyHistoryManager *PropertyHistoryManager
}

func NewDB() *DB {
	db := DB{
		started: false,
	}

	db.eventManager = &EventManager{
		db: &db,
	}
	db.propertyHistoryManager = newPropertyHistoryManager(&db)

	return &db
}

func (db *DB) StartWithPath(dataDir string) error {
	if dirInfo, err := os.Stat(dataDir); err != nil || !dirInfo.IsDir() {
		return errors.New("CriticalMoments: Data directory path does not exist")
	}

	dbPath := fmt.Sprintf("file:%s/critical_moments_db.db?_journal_mode=WAL&mode=rwc", dataDir)

	sqldb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	db.databasePath = dbPath
	db.sqldb = sqldb

	err = db.migrate()
	if err != nil {
		return err
	}

	db.started = true
	return nil
}

func (db *DB) Close() error {
	db.started = false
	return db.sqldb.Close()
}

func (db *DB) EventManager() *EventManager {
	return db.eventManager
}

func (db *DB) PropertyHistoryManager() *PropertyHistoryManager {
	return db.propertyHistoryManager
}

// migrations can be run on each start because they are incremental and non-destructive
// Future migrations must also be incremental (only append to this, check if not exists),
// or we must implement a versioning system
func (db *DB) migrate() error {
	if db.sqldb == nil {
		return errors.New("CriticalMoments: DB not started")
	}

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

		CREATE TABLE IF NOT EXISTS property_history (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			type INTEGER NOT NULL,
			int_value INTEGER,
			text_value TEXT,
			real_value REAL,
			numeric_value NUMERIC,
			sample_type INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		);

		CREATE INDEX IF NOT EXISTS property_history_name_created_at ON property_history (name, created_at);

		CREATE TRIGGER IF NOT EXISTS insert_property_history_created_at
		AFTER INSERT ON property_history
		BEGIN
			UPDATE property_history SET created_at =unixepoch('subsec') WHERE id = NEW.id;
		END;

		CREATE TRIGGER IF NOT EXISTS update_property_history_updated_at
		AFTER UPDATE ON property_history
		BEGIN
			UPDATE property_history SET updated_at =unixepoch('subsec') WHERE id = NEW.id;
		END;
	`)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) InsertEvent(e *datamodel.Event) error {
	if !db.started {
		return errors.New("CriticalMoments: DB not started")
	}

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
	if !db.started {
		return 0, errors.New("CriticalMoments: DB not started")
	}

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
	if !db.started {
		return 0, errors.New("CriticalMoments: DB not started")
	}

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
	if !db.started {
		return nil, errors.New("CriticalMoments: DB not started")
	}

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

type DBPropertyType int

const (
	DBPropertyTypeString DBPropertyType = 1
	DBPropertyTypeInt    DBPropertyType = 2
	DBPropertyTypeFloat  DBPropertyType = 3
	DBPropertyTypeBool   DBPropertyType = 4
	DBPropertyTypeTime   DBPropertyType = 5
)

func DBPropertyTypeIntFromKind(k reflect.Kind) (DBPropertyType, error) {
	switch k {
	case reflect.String:
		return DBPropertyTypeString, nil
	case reflect.Int:
		return DBPropertyTypeInt, nil
	case reflect.Float64:
		return DBPropertyTypeFloat, nil
	case reflect.Bool:
		return DBPropertyTypeBool, nil
	case datamodel.CMTimeKind:
		return DBPropertyTypeTime, nil
	default:
		return 0, errors.New("CriticalMoments: Unsupported property type")
	}
}

func (db *DB) InsertPropertyHistory(name string, value interface{}, sampleType datamodel.CMPropertySampleType) error {
	if !db.started {
		return errors.New("CriticalMoments: DB not started")
	}

	propKind := datamodel.CMTypeFromValue(value)
	dbType, err := DBPropertyTypeIntFromKind(propKind)
	if err != nil {
		return err
	}

	sqlTemplate := ""
	switch dbType {
	case DBPropertyTypeString:
		sqlTemplate = `INSERT INTO property_history (name, type, text_value, sample_type) VALUES (?, ?, ?, ?)`
	case DBPropertyTypeInt:
		sqlTemplate = `INSERT INTO property_history (name, type, int_value, sample_type) VALUES (?, ?, ?, ?)`
	case DBPropertyTypeFloat:
		sqlTemplate = `INSERT INTO property_history (name, type, real_value, sample_type) VALUES (?, ?, ?, ?)`
	case DBPropertyTypeBool:
		sqlTemplate = `INSERT INTO property_history (name, type, numeric_value, sample_type) VALUES (?, ?, ?, ?)`
	case DBPropertyTypeTime:
		// Time stored as microseconds, in int column
		time, ok := value.(time.Time)
		if !ok {
			return errors.New("CriticalMoments: Invalid time")
		}
		value = time.UnixMicro()
		sqlTemplate = `INSERT INTO property_history (name, type, int_value, sample_type) VALUES (?, ?, ?, ?)`
	}
	if sqlTemplate == "" {
		return errors.New("CriticalMoments: Unsupported property type")
	}

	_, err = db.sqldb.Exec(sqlTemplate, name, dbType, value, sampleType)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) LatestPropertyHistory(name string) (interface{}, error) {
	if !db.started {
		return nil, errors.New("CriticalMoments: DB not started")
	}

	var text_value sql.NullString
	var int_value sql.NullInt64
	var real_value sql.NullFloat64
	var numeric_value sql.NullBool
	var dbType sql.NullInt64
	err := db.sqldb.
		QueryRow(`SELECT text_value, int_value, real_value, numeric_value, type FROM property_history WHERE name = ? ORDER BY created_at DESC LIMIT 1`, name).
		Scan(&text_value, &int_value, &real_value, &numeric_value, &dbType)

	if err != nil {
		return nil, err
	}
	if !dbType.Valid {
		return nil, errors.New("CriticalMoments: Property type invalid")
	}

	switch DBPropertyType(dbType.Int64) {
	case DBPropertyTypeString:
		if text_value.Valid {
			return text_value.String, nil
		}
	case DBPropertyTypeInt:
		if int_value.Valid {
			return int_value.Int64, nil
		}
	case DBPropertyTypeFloat:
		if real_value.Valid {
			return real_value.Float64, nil
		}
	case DBPropertyTypeBool:
		if numeric_value.Valid {
			return numeric_value.Bool, nil
		}
	case DBPropertyTypeTime:
		if int_value.Valid {
			return time.UnixMicro(int_value.Int64), nil
		}
	default:
		return nil, errors.New("CriticalMoments: Unsupported property type")
	}

	return nil, errors.New("CriticalMoments: Invalid property value")
}
