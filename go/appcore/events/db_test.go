package events

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
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

func TestTestDB(t *testing.T) {
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
