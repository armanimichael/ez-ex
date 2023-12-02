package ezex

import (
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

//go:embed db/tables.sql
var dbInitScriptSQL string

const dbName = "user-data.db"
const UserDataDir = ".ez-ex"

func OpenDB() (*sql.DB, error) {
	home, _ := os.UserHomeDir()
	_ = os.MkdirAll(path.Join(home, UserDataDir), 0700)

	dsName := fmt.Sprintf("file:%s/.ez-ex/%s?_foreign_keys=true", home, dbName)
	return sql.Open("sqlite3", dsName)
}

func MigrateDB(db *sql.DB) error {
	_, err := db.Exec(dbInitScriptSQL)
	return err
}
