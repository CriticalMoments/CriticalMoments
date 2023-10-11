package events

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

func testBuildTestDb(t *testing.T) *DB {
	dataPath := fmt.Sprintf("/tmp/criticalmoments/test-temp-%v", rand.Int())
	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	db, err := NewDB(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestDBConnection(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	r, err := db.sqldb.Query("SELECT 99")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	r.Next()
	var v int
	err = r.Scan(&v)
	if err != nil {
		t.Fatal(err)
	}

	if v != 99 {
		t.Fatal("DB test failed")
	}
}

func TestDBMigrate(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()
	// Test contructor runs migration
	testSchema(db, t)

	// run again to make sure it doesn't fail
	err := db.migrate()
	if err != nil {
		t.Fatal(err)
	}
	testSchema(db, t)

	// run cold from new DB instance (schema should already exist)
	db.Close()
	db2, err := NewDB(path.Dir(db.databasePath)[5:])
	if err != nil {
		t.Fatal(err)
	}
	testSchema(db2, t)
}

func testSchema(db *DB, t *testing.T) {
	r, err := db.sqldb.Query("SELECT name FROM sqlite_schema WHERE type='table' AND name='events'")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	r.Next()
	var v string
	err = r.Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "events" {
		t.Fatal("DB migation failed")
	}

	r, err = db.sqldb.Query("SELECT name FROM sqlite_schema WHERE type='trigger' AND name='insert_events_created_at'")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	r.Next()
	err = r.Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "insert_events_created_at" {
		t.Fatal("DB migration failed")
	}

	r, err = db.sqldb.Query("SELECT name FROM sqlite_schema WHERE type='trigger' AND name='update_events_updated_at'")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	r.Next()
	err = r.Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "update_events_updated_at" {
		t.Fatal("DB migration failed")
	}
}

func BenchmarkWarmMigrate(b *testing.B) {
	t := &testing.T{}
	db := testBuildTestDb(t)
	defer db.Close()

	for i := 0; i < b.N; i++ {
		err := db.migrate()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreatedAtTrigger(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	// insert a row into events
	_, err := db.sqldb.Exec(`
		INSERT INTO events (event_name, event_data)
		VALUES ('test', 'test')
	`)
	if err != nil {
		t.Fatal(err)
	}

	// select the last row inserted
	r, err := db.sqldb.Query(`
		SELECT created_at, updated_at FROM events
		LIMIT 1
	`)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	r.Next()
	var c float64
	var u float64
	err = r.Scan(&c, &u)
	if err != nil {
		t.Fatal(err)
	}
	now := float64(time.Now().UnixNano()) / float64(time.Second)
	if math.Abs(now-c) > 0.01 {
		t.Fatal("Trigger failed to set created_at")
	}
	if math.Abs(now-u) > 0.01 {
		t.Fatal("Trigger failed to set updated_at")
	}

	// update the row after small delay
	time.Sleep(time.Millisecond * 2)
	_, err = db.sqldb.Exec(`
		UPDATE events SET event_name = 'test2'
		WHERE event_name = 'test'
	`)
	if err != nil {
		t.Fatal(err)
	}

	// check created_at has not changed, but updated_at has
	r, err = db.sqldb.Query(`
		SELECT created_at, updated_at FROM events
		LIMIT 1
	`)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	r.Next()
	var c2 float64
	var u2 float64
	err = r.Scan(&c2, &u2)
	if err != nil {
		t.Fatal(err)
	}
	if c != c2 {
		t.Fatal("created_at changed")
	}
	if u == u2 {
		t.Fatal("Trigger did not change updated_at")
	}
	now = float64(time.Now().UnixNano()) / float64(time.Second)
	if math.Abs(now-u2) > 0.01 {
		t.Fatal("Trigger failed to set updated_at")
	}
}
