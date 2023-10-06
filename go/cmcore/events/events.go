package events

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (int, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	r, err := db.Query("SELECT 99")
	if err != nil {
		return 0, err
	}
	defer r.Close()
	r.Next()
	var v int
	err = r.Scan(&v)
	if err != nil {
		return 0, err
	}

	return v, nil
}
