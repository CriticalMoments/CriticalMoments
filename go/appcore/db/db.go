package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strings"
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

	// WAL mode for better performance/concurrency
	dbPath := fmt.Sprintf("file:%s/critical_moments_db.db?_journal_mode=WAL&mode=rwc", dataDir)
	sqldb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	// no DB concurrency always a good idea for SQLite. No impact to benchmarks. Do need to be careful to release connections asap
	sqldb.SetMaxOpenConns(1)

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
			name TEXT NOT NULL,
			type INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		);

		CREATE INDEX IF NOT EXISTS events_name_created_at ON events (name, created_at);

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
		INSERT INTO events (name, type)
		VALUES (?, ?)
	`, e.Name, e.EventType)
	if err != nil {
		return err
	}

	return nil
}

const eventCountByNameQuery = `SELECT COUNT(*) FROM events WHERE name = ?`

func (db *DB) EventCountByName(name string) (int, error) {
	if !db.started {
		return 0, errors.New("CriticalMoments: DB not started")
	}

	var count int
	err := db.sqldb.QueryRow(eventCountByNameQuery, name).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const eventCountByNameWithLimitQuery = `SELECT COUNT(*) FROM (SELECT id FROM events WHERE name = ? LIMIT ?)`

func (db *DB) EventCountByNameWithLimit(name string, limit int) (int, error) {
	if !db.started {
		return 0, errors.New("CriticalMoments: DB not started")
	}

	var count int
	err := db.sqldb.QueryRow(eventCountByNameWithLimitQuery, name, limit).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const latestEventTimeByNameQuery = `SELECT created_at FROM events WHERE name = ? ORDER BY created_at DESC LIMIT 1`

func (db *DB) LatestEventTimeByName(name string) (*time.Time, error) {
	if !db.started {
		return nil, errors.New("CriticalMoments: DB not started")
	}

	var epochTime float64
	err := db.sqldb.QueryRow(latestEventTimeByNameQuery, name).Scan(&epochTime)
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

const latestPropHistoryTimeByNameQuery = `SELECT created_at FROM property_history WHERE name = ? ORDER BY created_at DESC LIMIT 1`

func (db *DB) latestPropertyHistoryTime(name string) (*time.Time, error) {
	var epochTime float64
	err := db.sqldb.
		QueryRow(latestPropHistoryTimeByNameQuery, name).
		Scan(&epochTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_, fractionalSeconds := math.Modf(epochTime)
	nanoseconds := int64(fractionalSeconds * 1_000_000_000)
	time := time.Unix(int64(epochTime), nanoseconds)
	return &time, nil
}

var maxTimeBetweenPropertyHistorySamples = time.Minute * 5

const insertPropertyHistorySqlTemplate = `INSERT INTO property_history (name, type, TYPE_VAL, sample_type) VALUES (?, ?, ?, ?)`

func (db *DB) InsertPropertyHistory(name string, value interface{}, sampleType datamodel.CMPropertySampleType) error {
	if !db.started {
		return errors.New("CriticalMoments: DB not started")
	}

	// Check last update time, and skip if it's in last 5 mins
	latestHistoryTime, err := db.latestPropertyHistoryTime(name)
	if err != nil {
		return err
	}
	if latestHistoryTime != nil {
		if time.Now().Before(latestHistoryTime.Add(maxTimeBetweenPropertyHistorySamples)) {
			return nil
		}
	}

	propKind := datamodel.CMTypeFromValue(value)
	dbType, err := DBPropertyTypeIntFromKind(propKind)
	if err != nil {
		return err
	}

	sqlTemplate, value, err := formatSqlForPropHistoryType(value, insertPropertyHistorySqlTemplate)
	if err != nil {
		return err
	}

	_, err = db.sqldb.Exec(sqlTemplate, name, dbType, value, sampleType)
	if err != nil {
		return err
	}

	return nil
}

const latestPropertyHistoryValueByNameQuery = `SELECT text_value, int_value, real_value, numeric_value, type FROM property_history WHERE name = ? ORDER BY created_at DESC LIMIT 1`

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
		QueryRow(latestPropertyHistoryValueByNameQuery, name).
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

const propertyHistoryEverHadValueQuery = `SELECT COUNT(*) FROM property_history WHERE name = ? AND TYPE_VAL = ? LIMIT 1`

func (db *DB) PropertyHistoryEverHadValue(name string, value interface{}) (bool, error) {
	if !db.started {
		return false, errors.New("CriticalMoments: DB not started")
	}

	sqlTemplate, value, err := formatSqlForPropHistoryType(value, propertyHistoryEverHadValueQuery)
	if err != nil {
		return false, err
	}

	var count sql.NullInt64
	err = db.sqldb.QueryRow(sqlTemplate, name, value).Scan(&count)
	if err != nil {
		return false, err
	}
	if !count.Valid {
		return false, errors.New("CriticalMoments: unexpected error")
	}
	return count.Int64 > 0, nil
}

func (db *DB) DbConditionFunctions() map[string]*datamodel.ConditionDynamicFunction {
	return map[string]*datamodel.ConditionDynamicFunction{
		"eventCount": {
			Function: func(params ...any) (any, error) {
				// Parameter type+count checking is done with the Types signature
				count, err := db.EventCountByName(params[0].(string))
				if err != nil {
					return nil, err
				}
				return count, nil
			},
			Types: []any{new(func(string) int)},
		},
		"eventCountWithLimit": {
			Function: func(params ...any) (any, error) {
				// Parameter type+count checking is done the Types signature
				count, err := db.EventCountByNameWithLimit(params[0].(string), params[1].(int))
				if err != nil {
					return nil, err
				}
				return count, nil
			},
			Types: []any{new(func(string, int) int)},
		},
		"propertyHistoryLatestValue": {
			Function: func(params ...any) (any, error) {
				// Parameter type+count checking is done the Types signature
				value, err := db.LatestPropertyHistory(params[0].(string))
				// no rows should return nil
				if err == sql.ErrNoRows {
					return nil, nil
				}
				if err != nil {
					return nil, err
				}
				return value, nil
			},
			Types: []any{new(func(string) interface{})},
		},
		"propertyEverHadValue": {
			Function: func(params ...any) (any, error) {
				// Parameter type+count checking is done the Types signature
				value, err := db.PropertyHistoryEverHadValue(params[0].(string), params[1])
				if err != nil {
					return nil, err
				}
				return value, nil
			},
			Types: []any{new(func(string, interface{}) bool)},
		},
		"stableRand": {
			Function: func(params ...any) (any, error) {
				// Parameter type+count checking is done the Types signature
				return db.StableRandom()
			},
			Types: []any{new(func() int64)},
		},
	}
}

func formatSqlForPropHistoryType(val any, sqlTemplate string) (string, any, error) {
	dbType, err := DBPropertyTypeIntFromKind(datamodel.CMTypeFromValue(val))
	if err != nil {
		return "", nil, err
	}

	column := ""
	switch dbType {
	case DBPropertyTypeString:
		column = "text_value"
	case DBPropertyTypeInt:
		column = "int_value"
	case DBPropertyTypeFloat:
		column = "real_value"
	case DBPropertyTypeBool:
		column = "numeric_value"
	case DBPropertyTypeTime:
		// Time stored as microseconds, in int column
		time, ok := val.(time.Time)
		if !ok {
			return "", nil, errors.New("CriticalMoments: Invalid time")
		}
		val = time.UnixMicro()
		column = "int_value"
	}
	if column == "" {
		return "", nil, errors.New("CriticalMoments: Unsupported property type")
	}

	sql := strings.Replace(sqlTemplate, "TYPE_VAL", column, -1)
	return sql, val, nil
}

func (db *DB) StableRandom() (int64, error) {
	newRandom := rand.Int63()

	r, err := db.sqldb.Exec(`
			INSERT INTO property_history (name, type, int_value, sample_type)
			  SELECT 'stable_random', ?, ?, ?
				WHERE NOT EXISTS (SELECT 1 FROM property_history WHERE name = 'stable_random' LIMIT 1);
	`, DBPropertyTypeInt, newRandom, datamodel.CMPropertySampleTypeDoNotSample)
	if err != nil {
		return 0, err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rows == 1 {
		return newRandom, nil
	}

	var existingRandom sql.NullInt64
	err = db.sqldb.QueryRow(`
		SELECT int_value FROM property_history WHERE name = 'stable_random' ORDER BY created_at LIMIT 1;
		`).Scan(&existingRandom)
	if err != nil {
		return 0, err
	}
	if !existingRandom.Valid {
		return 0, errors.New("CriticalMoments: unexpected error")
	}

	return existingRandom.Int64, nil
}
