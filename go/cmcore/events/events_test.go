package events

import "testing"

func TestTestDB(t *testing.T) {
	r, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	if r != 99 {
		t.Fatal("DB test failed")
	}
}
