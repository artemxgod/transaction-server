package sqlstore_test

import (
	"os"
	"testing"
)

var databaseURL string

//	TestMain will be automatically called before testing
//	This will let us to configure before testing

// TestMain helps us to configure our databaseURL in case we call it from another place
func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "user=postgres dbname=tnserver sslmode=disable"
	}

	os.Exit(m.Run())
}