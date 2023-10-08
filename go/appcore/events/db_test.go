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
	r, err := db.testDB()
	if err != nil {
		t.Fatal(err)
	}
	if r != 99 {
		t.Fatal("DB test failed")
	}
}
