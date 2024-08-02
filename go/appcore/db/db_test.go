package db

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
	db := NewDB()
	err = db.StartWithPath(dataPath)
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
	err := migrate(db.sqldb)
	if err != nil {
		t.Fatal(err)
	}
	testSchema(db, t)

	// run cold from new DB instance (schema should already exist)
	db.Close()
	db2 := NewDB()
	err = db2.StartWithPath(path.Dir(db.databasePath)[5:])
	if err != nil {
		t.Fatal(err)
	}
	testSchema(db2, t)
}

func testSchema(db *DB, t *testing.T) {
	var v string
	err := db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='table' AND name='events'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "events" {
		t.Fatal("DB migation failed")
	}

	err = db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='index' AND name='events_name_created_at'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "events_name_created_at" {
		t.Fatal("DB migration failed")
	}

	err = db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='trigger' AND name='insert_events_created_at'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "insert_events_created_at" {
		t.Fatal("DB migration failed")
	}

	err = db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='trigger' AND name='update_events_updated_at'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "update_events_updated_at" {
		t.Fatal("DB migration failed")
	}

	err = db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='table' AND name='property_history'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "property_history" {
		t.Fatal("DB migation failed")
	}

	err = db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='index' AND name='property_history_name_created_at'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "property_history_name_created_at" {
		t.Fatal("DB migration failed")
	}

	err = db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='trigger' AND name='insert_property_history_created_at'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "insert_property_history_created_at" {
		t.Fatal("DB migration failed")
	}

	err = db.sqldb.QueryRow("SELECT name FROM sqlite_schema WHERE type='trigger' AND name='update_property_history_updated_at'").Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if v != "update_property_history_updated_at" {
		t.Fatal("DB migration failed")
	}
}

func BenchmarkWarmMigrate(b *testing.B) {
	t := &testing.T{}
	db := testBuildTestDb(t)
	defer db.Close()

	for i := 0; i < b.N; i++ {
		err := migrate(db.sqldb)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreatedAtTrigger(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	// insert a row into table
	_, err := db.sqldb.Exec(`
		INSERT INTO events (name, type)
		VALUES ('test', 0)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// select the last row inserted
	var c float64
	var u float64
	err = db.sqldb.QueryRow(`
		SELECT created_at, updated_at FROM events
		LIMIT 1
	`).Scan(&c, &u)
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
		UPDATE events SET name = 'test2'
		WHERE name = 'test'
	`)
	if err != nil {
		t.Fatal(err)
	}

	// check created_at has not changed, but updated_at has
	var c2 float64
	var u2 float64
	err = db.sqldb.QueryRow(`
		SELECT created_at, updated_at FROM events
		LIMIT 1
	`).Scan(&c2, &u2)
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

	// request before rows exist
	ct, err := db.LatestEventTimeByName("event_not_in_db")
	if err != nil {
		t.Fatal(err)
	}
	if ct != nil {
		t.Fatal("LatestEventTimeByName didn't return nil")
	}

	// insert a row into events
	e, err := datamodel.NewCustomEventWithName("test")
	if err != nil {
		t.Fatal(err)
	}
	err = db.InsertEvent(e)
	if err != nil {
		t.Fatal(err)
	}

	ct, err = db.LatestEventTimeByName(e.Name)
	if err != nil {
		t.Fatal(err)
	}
	if ct == nil {
		t.Fatal("LatestEventTimeByName returned nil")
	}
	if math.Abs(time.Since(*ct).Seconds()) > 0.01 {
		t.Fatal("LatestEventTimeByName returned wrong time")
	}
	ft, err := db.FirstEventTimeByName(e.Name)
	if err != nil {
		t.Fatal(err)
	}
	if ft.Compare(*ct) != 0 {
		t.Fatal("First and latest event time should return same result when there's only one event")
	}

	// insert another event
	time.Sleep(time.Millisecond * 2)
	err = db.InsertEvent(e)
	if err != nil {
		t.Fatal(err)
	}
	ct2, err := db.LatestEventTimeByName(e.Name)
	if err != nil {
		t.Fatal(err)
	}
	ft2, _ := db.FirstEventTimeByName(e.Name)
	if ft2.Compare(*ct) != 0 {
		t.Fatal("FirstEventTimeByName should return same result as the first latest")
	}
	// Confirm latest is sorting correctly
	if ct2.Compare(*ct) != 1 {
		t.Fatal("LatestEventTimeByName returned older time")
	}

	// confirm type is set
	var name string
	var eventType int
	err = db.sqldb.QueryRow(`
		SELECT name, type FROM events
		LIMIT 1
	`).Scan(&name, &eventType)
	if err != nil {
		t.Fatal(err)
	}
	if name != "test" || eventType != int(datamodel.EventTypeCustom) {
		t.Fatal("insert failed")
	}
}

func TestEventCount(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	for i := 1; i < 10; i++ {
		// insert a row into events
		e, err := datamodel.NewCustomEventWithName("test")
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
	testSqlExplainIncludes(latestEventTimeByNameQuery, "USING COVERING INDEX events_name_created_at", t, "test")
}

func TestEventCountLimitUsesIndex(t *testing.T) {
	testSqlExplainIncludes(eventCountByNameWithLimitQuery, "USING COVERING INDEX events_name_created_at", t, "test", 5) // add_test_count
}

func TestEventCountUsesIndex(t *testing.T) {
	testSqlExplainIncludes(eventCountByNameQuery, "USING COVERING INDEX events_name_created_at", t, "test") // add_test_count
}

func TestPropertyQueriesIndex(t *testing.T) {
	testSqlExplainIncludes(latestPropHistoryTimeByNameQuery, "USING COVERING INDEX property_history_name_created_at", t, "test")          // add_test_count
	testSqlExplainIncludes(latestPropertyHistoryValueByNameQuery, "USING INDEX property_history_name_created_at", t, "test", 1, "val", 1) // add_test_count
	everHadSql := strings.Replace(propertyHistoryEverHadValueQuery, "TYPE_VAL", "text_value", -1)
	testSqlExplainIncludes(everHadSql, "USING INDEX property_history_name_created_at", t, "test", "val") // add_test_count
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

func TestCreatedAtTriggerPropHistory(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	// insert a row into table
	_, err := db.sqldb.Exec(`
		INSERT INTO property_history (name, type, text_value, sample_type)
		VALUES ('test', ?, 'val', 1)
	`, DBPropertyTypeString)
	if err != nil {
		t.Fatal(err)
	}

	// select the last row inserted
	var c float64
	var u float64
	err = db.sqldb.QueryRow(`
		SELECT created_at, updated_at FROM property_history
		LIMIT 1
	`).Scan(&c, &u)
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
		UPDATE property_history SET text_value = 'val2'
		WHERE name = 'test'
	`)
	if err != nil {
		t.Fatal(err)
	}

	// check created_at has not changed, but updated_at has
	var c2 float64
	var u2 float64
	err = db.sqldb.QueryRow(`
		SELECT created_at, updated_at FROM property_history 
		LIMIT 1
	`).Scan(&c2, &u2)
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

func TestInsertAndRetrievePropHistory(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	// insert a row into property history
	// Other types tested in property_hisotry_manager_test.go
	err := db.InsertPropertyHistory("testx", "valx", 1)
	if err != nil {
		t.Fatal(err)
	}

	// retrieve and verify
	var n string
	var v string
	var s int
	err = db.sqldb.QueryRow(`
		SELECT name, text_value, sample_type FROM property_history 
		LIMIT 1
	`).Scan(&n, &v, &s)
	if err != nil {
		t.Fatal(err)
	}
	if n != "testx" || v != "valx" || s != 1 {
		t.Fatal("retrieve failed")
	}

	// retrieve and verify with helper
	// All types tested in property_registry_test.go
	rv, err := db.LatestPropertyHistory("testx")
	if err != nil {
		t.Fatal(err)
	}
	if rv != "valx" {
		t.Fatal("retrieve failed")
	}

	// check if it has ever had value
	has, err := db.PropertyHistoryEverHadValue("testx", "valx")
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal("PropertyHistoryEverHadValue failed")
	}
	has, err = db.PropertyHistoryEverHadValue("testx", "wrong value")
	if err != nil {
		t.Fatal(err)
	}
	if has {
		t.Fatal("PropertyHistoryEverHadValue failed")
	}
}

func TestPropHistoryRateLimit(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	delay := time.Millisecond * 15
	original := maxTimeBetweenPropertyHistorySamples
	maxTimeBetweenPropertyHistorySamples = delay
	defer func() {
		maxTimeBetweenPropertyHistorySamples = original
	}()

	err := db.InsertPropertyHistory("test", "val1", datamodel.CMPropertySampleTypeOnUse)
	if err != nil {
		t.Fatal(err)
	}
	err = db.InsertPropertyHistory("test", "val2", datamodel.CMPropertySampleTypeOnUse)
	if err != nil {
		t.Fatal(err)
	}

	// Immediate write should be rate limited, and only store first value
	rv, err := db.LatestPropertyHistory("test")
	if err != nil {
		t.Fatal(err)
	}
	if rv != "val1" {
		t.Fatal("set failed to rate limit second write")
	}

	// Delayed, writes should work again
	time.Sleep(delay + time.Millisecond)
	err = db.InsertPropertyHistory("test", "val3", datamodel.CMPropertySampleTypeOnUse)
	if err != nil {
		t.Fatal(err)
	}
	rv, err = db.LatestPropertyHistory("test")
	if err != nil {
		t.Fatal(err)
	}
	if rv != "val3" {
		t.Fatal("set failed to block second write")
	}
}

func TestStableRandom(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	initialRand, err := db.StableRandom()
	if err != nil {
		t.Fatal(err)
	}
	if initialRand == 0 {
		t.Fatal("StableRandom returned 0")
	}

	nextRand, err := db.StableRandom()
	if err != nil {
		t.Fatal(err)
	}
	if initialRand != nextRand {
		t.Fatal("StableRandom returned different values")
	}
}

func createTestDb(path string) (*DB, error) {
	db := NewDB()
	err := db.StartWithPath(path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestPathChecks(t *testing.T) {
	dataPath := fmt.Sprintf("/tmp/criticalmoments/test-temp-%v", rand.Int())
	db, err := createTestDb(dataPath)
	if err == nil {
		t.Fatal("Failed to check path exists")
	}
	if db != nil {
		t.Fatal("set invalid path")
	}

	filePath := fmt.Sprintf("/tmp/cm-test-temp-%v.txt", rand.Int())
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	db, err = createTestDb(dataPath)
	if err == nil {
		t.Fatal("Failed to check path exists and is dir")
	}
	if db != nil {
		t.Fatal("set invalid path")
	}

	os.MkdirAll(dataPath, os.ModePerm)
	db, err = createTestDb(dataPath)
	expectedPath := fmt.Sprintf("file:%s/critical_moments_db.db?_journal_mode=WAL&mode=rwc", dataPath)
	if err != nil || db.databasePath != expectedPath {
		t.Fatal("Failed to set data path")
	}
}

// Wild issue: if timestamps rounded to .0 seconds, they are returned as time.Time, and if not they are returned as float64
func TestTimestampRoundingAndLatestPropHistory(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	// insert a row into table
	_, err := db.sqldb.Exec(`
		INSERT INTO property_history (name, type, text_value, sample_type)
		VALUES ('test', ?, 'val', 1)
	`, DBPropertyTypeString)
	if err != nil {
		t.Fatal(err)
	}

	ct, err := db.latestPropertyHistoryTime("test")
	if err != nil {
		t.Fatal(err)
	}
	if ct == nil {
		t.Fatal("LatestPropHistoryTime returned nil")
	}
	if math.Abs(time.Since(*ct).Seconds()) > 0.01 {
		t.Fatal("LatestPropHistoryTime returned wrong time")
	}

	// update the row to exact time, no milliseconds
	// This previously caused the return type to change to time.Time, and errored
	time.Sleep(time.Millisecond * 2)
	_, err = db.sqldb.Exec(`
		UPDATE property_history SET created_at = 1710791550
		WHERE name = 'test'
	`)
	if err != nil {
		t.Fatal(err)
	}

	ct, err = db.latestPropertyHistoryTime("test")
	if err != nil {
		t.Fatal(err)
	}
	if ct == nil {
		t.Fatal("LatestPropHistoryTime returned nil")
	}
	if math.Abs(float64(ct.Unix())-1710791550.0) > 0.1 {
		t.Fatal("LatestPropHistoryTime returned wrong time")
	}
}

func TestAllEventTimesByName(t *testing.T) {
	db := testBuildTestDb(t)
	defer db.Close()

	// Should be empty to start
	times, err := db.AllEventTimesByName("test")
	if err != nil {
		t.Fatal(err)
	}
	if len(times) != 0 {
		t.Fatal("Expected empty list")
	}

	// Insert a few events
	startTime := time.Now()
	for i := 0; i < 10; i++ {
		// insert a row into events
		e, err := datamodel.NewCustomEventWithName("test")
		if err != nil {
			t.Fatal(err)
		}
		err = db.InsertEvent(e)
		if err != nil {
			t.Fatal(err)
		}
	}

	times, err = db.AllEventTimesByName("test")
	if err != nil {
		t.Fatal(err)
	}
	if len(times) != 10 {
		t.Fatal("Expected 10 events")
	}
	// Check they are assending
	var priorTime = startTime.Add(-1 * time.Second) // start time is inclusive
	for i, tm := range times {
		if tm.Before(priorTime) {
			t.Fatal("Expected events in ascending order. Got ", tm, " before ", priorTime, " for event ", i)
		}
		if startTime.Sub(tm).Abs() > time.Second {
			t.Fatal("Expected events within 1 second of start time. Got ", tm, " for event ", i)
		}
		priorTime = tm
	}
}
