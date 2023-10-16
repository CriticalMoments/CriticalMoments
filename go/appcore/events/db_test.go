package events

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
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

	r, err = db.sqldb.Query("SELECT name FROM sqlite_schema WHERE type='index' AND name='events_event_name_created_at'")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	r.Next()
	err = r.Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "events_event_name_created_at" {
		t.Fatal("DB migration failed")
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
		INSERT INTO events (event_name)
		VALUES ('test')
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

func TestInsertAndRetrieve(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	// insert a row into events
	e, err := datamodel.NewEventWithName("test")
	if err != nil {
		t.Fatal(err)
	}
	err = db.InsertEvent(e)
	if err != nil {
		t.Fatal(err)
	}

	ct, err := db.LatestEventTimeByName("test")
	if err != nil {
		t.Fatal(err)
	}
	if t == nil {
		t.Fatal("LatestEventTimeByName returned nil")
	}
	if math.Abs(time.Since(*ct).Seconds()) > 0.01 {
		t.Fatal("LatestEventTimeByName returned wrong time")
	}

	// insert another event
	time.Sleep(time.Millisecond * 2)
	err = db.InsertEvent(e)
	if err != nil {
		t.Fatal(err)
	}
	ct2, err := db.LatestEventTimeByName("test")
	if err != nil {
		t.Fatal(err)
	}
	// Confirm latest is sorting correctly
	if ct2.Compare(*ct) != 1 {
		t.Fatal("LatestEventTimeByName returned older time")
	}
}

func TestEventCount(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	for i := 1; i < 10; i++ {
		// insert a row into events
		e, err := datamodel.NewEventWithName("test")
		if err != nil {
			t.Fatal(err)
		}
		err = db.InsertEvent(e)
		if err != nil {
			t.Fatal(err)
		}
		// count events
		count, err := db.EventCountByName("test")
		if err != nil {
			t.Fatal(err)
		}
		if count != i {
			t.Fatal("EventCountByName returned wrong count")
		}
	}

	// Check count with limit
	count, err := db.EventCountByNameWithLimit("test", 100)
	if err != nil {
		t.Fatal(err)
	}
	if count != 9 {
		t.Fatal("EventCountByNameWithLimit returned wrong count")
	}

	count, err = db.EventCountByNameWithLimit("test", 5)
	if err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Fatal("EventCountByNameWithLimit returned count past limit")
	}
}

func TestLatestEventUsesIndex(t *testing.T) {
	testSqlExplainIncludes(latestEventTimeByNameQuery, "USING COVERING INDEX events_event_name_created_at", t, "test")
}

func TestEventCountLimitUsesIndex(t *testing.T) {
	testSqlExplainIncludes(eventCountByNameWithLimitQuery, "USING COVERING INDEX events_event_name_created_at", t, "test", 5)
}

func TestEventCountUsesIndex(t *testing.T) {
	testSqlExplainIncludes(eventCountByNameQuery, "USING COVERING INDEX events_event_name_created_at", t, "test")
}

func testSqlExplainIncludes(sql string, expectedExplain string, t *testing.T, args ...any) {
	db := testBuildTestDb(t)
	defer db.Close()

	explainSql := fmt.Sprintf("EXPLAIN QUERY PLAN %s", sql)
	r, err := db.sqldb.Query(explainSql, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	foundIndex := false
	for r.Next() {
		var er1 string
		var er2 string
		var er3 string
		var er4 string
		err = r.Scan(&er1, &er2, &er3, &er4)
		if err != nil {
			t.Fatal(err)
		}
		// SQLite does not guaruntee this string is consistent across versions, so this string check
		// may need to be updated if the SQLite version changes. Still good to have this check to ensure
		// our perf design never breaks silently (SQL footgun #732)
		if strings.Contains(er4, expectedExplain) {
			foundIndex = true
		}
	}
	if !foundIndex {
		t.Fatal("query does not use index: ", sql)
	}

}