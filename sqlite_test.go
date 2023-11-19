package ezex

import (
	"database/sql"
	"log"
	"os"
	"path"
	"testing"
)

const testDBName = "./test-db.db"

var testDB *sql.DB

func cleanTestDB() {
	pwd, _ := os.Getwd()
	_ = os.Remove(path.Join(pwd, testDBName))
}

func TestMain(m *testing.M) {
	cleanTestDB()

	var err error
	testDB, err = sql.Open("sqlite3", testDBName)
	if err != nil {
		log.Fatal(err)
	}

	if err = MigrateDB(testDB); err != nil {
		log.Fatalf("Error migrating the DB: %s", err)
	}

	code := m.Run()

	defer func(db *sql.DB, code int) {
		if err := testDB.Close(); err != nil {
			log.Fatal(err)
		}

		cleanTestDB()
		os.Exit(code)
	}(testDB, code)
}
