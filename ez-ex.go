package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

//go:embed db/tables.sql
var dbInitScriptSQL string

const dbName = "user-data.db"

func main() {
	home, _ := os.UserHomeDir()
	_ = os.MkdirAll(home+"/.ez-ex", 0700)
	_, err := os.Stat(dbName)
	dbExists := !errors.Is(err, os.ErrNotExist)

	dsName := fmt.Sprintf("file:%s/.ez-ex/%s?_foreign_keys=true", home, dbName)
	db, err := sql.Open("sqlite3", dsName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(db)

	if !dbExists {
		if _, err := db.Exec(dbInitScriptSQL); err != nil {
			log.Fatal(err)
		}
	}
}
