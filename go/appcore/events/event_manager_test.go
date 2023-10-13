package events

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestConstructor(t *testing.T) {
	dataPath := fmt.Sprintf("/tmp/criticalmoments/test-temp-%v", rand.Int())
	em, err := NewEventManager(dataPath)
	if err == nil {
		t.Fatal("Failed to check path exists")
	}
	if em != nil {
		t.Fatal("set invalid path")
	}

	filePath := fmt.Sprintf("/tmp/cm-test-temp-%v.txt", rand.Int())
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	em, err = NewEventManager(dataPath)
	if err == nil {
		t.Fatal("Failed to check path exists and is dir")
	}
	if em != nil {
		t.Fatal("set invalid path")
	}

	os.MkdirAll(dataPath, os.ModePerm)
	em, err = NewEventManager(dataPath)
	expectedPath := fmt.Sprintf("file:%s/critical_moments_db.db?_journal_mode=WAL&mode=rwc", dataPath)
	if err != nil || em.db.databasePath != expectedPath {
		t.Fatal("Failed to set data path")
	}
}
