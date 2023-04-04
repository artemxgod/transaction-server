package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

/* TestDB is a fuction that helps us test db functions
It opens database, and returns function that will clear and close db after testing*/
func TestDB(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper() // says our test that it is a help method, no need to test it

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}
		db.Close()
	}
}