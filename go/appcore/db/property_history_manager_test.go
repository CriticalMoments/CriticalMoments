package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func testPropertyHistoryManager(t testing.TB) (*DB, *PropertyHistoryManager, error) {
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

	return db, newPropertyHistoryManager(db), nil
}

// Test all permutations of property types with sample_type, and when they can be set (before/after startup)
func TestHistoryManager(t *testing.T) {
	dataPath := fmt.Sprintf("/tmp/criticalmoments/test-temp-%v", rand.Int())
	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	db := NewDB()
	phm := db.PropertyHistoryManager()

	baseProps := map[string]interface{}{
		"Text1":     "custom1val",
		"Text2":     "custom2val",
		"TextEmpty": "",
		"Int":       123,
		"IntZero":   0,
		"IntNeg":    -984594854958,
		"Float":     123.456,
		"FloatZero": 0.0,
		"FloatNeg":  -3498534958349.934,
		"BoolTrue":  true,
		"BoolFalse": false,
		"TimeNow":   time.Now(),
		"OtherTime": time.Now().Add(time.Hour),
	}
	customProps := map[string]interface{}{}
	startupProps := map[string]interface{}{}
	beforeStartUseProps := map[string]interface{}{}
	afterStartUseProps := map[string]interface{}{}
	afterStartCustomProps := map[string]interface{}{}
	allSets := map[string]map[string]interface{}{
		"custom":       customProps,
		"startup":      startupProps,
		"before":       beforeStartUseProps,
		"after":        afterStartUseProps,
		"after_custom": afterStartCustomProps,
	}
	for k, v := range baseProps {
		for setName, set := range allSets {
			set[setName+"_"+k] = v
		}
	}

	for k, v := range customProps {
		err = phm.CustomPropertySet(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	for k, v := range beforeStartUseProps {
		err := phm.UpdateHistoryForPropertyAccessed(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = db.StartWithPath(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	err = phm.TrackPropertyHistoryForStartup(startupProps)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range afterStartUseProps {
		err := phm.UpdateHistoryForPropertyAccessed(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	for k, v := range afterStartCustomProps {
		err = phm.CustomPropertySet(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	// retrieve and verify raw from DB
	r, err := db.sqldb.Query(`
		SELECT name, type, text_value, int_value, real_value, numeric_value, sample_type FROM property_history 
	`)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	hasRow := r.Next()
	var n string
	var rowtype int
	var tv sql.NullString
	var rv sql.NullFloat64
	var iv sql.NullInt64
	var bv sql.NullBool
	var s int
	for hasRow {
		err = r.Scan(&n, &rowtype, &tv, &iv, &rv, &bv, &s)
		if err != nil {
			t.Fatal(err)
		}
		var val interface{}
		for _, set := range allSets {
			v, ok := set[n]
			if ok {
				val = v
				delete(set, n)
				break
			}
		}
		if val == nil {
			t.Fatalf("unexpected property %s", n)
		}
		if rowtype == int(DBPropertyTypeString) {
			if !tv.Valid || val != tv.String {
				t.Fatalf("expected %s, got %s", val, tv.String)
			}
			if rv.Valid || iv.Valid || bv.Valid {
				t.Fatalf("unexpected value")
			}
		} else if rowtype == int(DBPropertyTypeInt) {
			if !iv.Valid || val != int(iv.Int64) {
				t.Fatalf("expected %d, got %d", val, iv.Int64)
			}
			if rv.Valid || tv.Valid || bv.Valid {
				t.Fatalf("unexpected value")
			}
		} else if rowtype == int(DBPropertyTypeFloat) {
			if !rv.Valid || val != rv.Float64 {
				t.Fatalf("expected %f, got %f", val, rv.Float64)
			}
			if iv.Valid || tv.Valid || bv.Valid {
				t.Fatalf("unexpected value")
			}
		} else if rowtype == int(DBPropertyTypeBool) {
			if !bv.Valid || val != bv.Bool {
				t.Fatalf("expected %f, got %v", val, bv.Bool)
			}
			if iv.Valid || tv.Valid || rv.Valid {
				t.Fatalf("unexpected value")
			}
		} else if rowtype == int(DBPropertyTypeTime) {
			vtime, ok := val.(time.Time)
			if !ok {
				t.Fatalf("expected time.Time, got %T", val)
			}
			if !iv.Valid || vtime.UnixMicro() != iv.Int64 {
				t.Fatalf("expected %d, got %v", vtime.UnixMicro(), iv.Int64)
			}
			if rv.Valid || tv.Valid || bv.Valid {
				t.Fatalf("unexpected value")
			}
		} else {
			t.Fatalf("unexpected type %d", rowtype)
		}
		hasRow = r.Next()
	}

	for _, set := range allSets {
		if len(set) != 0 {
			t.Fatalf("expected all properties to be found")
		}
	}
}

// Benchmark writes. Could batch startup into one transaction, but this is to verify it's not necessary
// Result: 35k writes/sec and we expect about 30 total so simple approach is great!
func Benchmark(b *testing.B) {
	_, phm, err := testPropertyHistoryManager(b)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		x := rand.Int()
		err = phm.UpdateHistoryForPropertyAccessed(fmt.Sprintf("test%d", x), x)
		if err != nil {
			b.Fatal(err)
		}
	}
}
