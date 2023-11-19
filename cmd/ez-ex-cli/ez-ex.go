package main

import (
	"database/sql"
	ezex "github.com/armanimichael/ez-ex"
	"log"
)

func main() {
	db, err := ezex.OpenDB()
	if err != nil {
		log.Fatalf("Error opening the DB: %s", err)
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			log.Fatalf("Error closing the DB: %s", err)
		}
	}(db)

	if err = ezex.MigrateDB(db); err != nil {
		log.Fatalf("Error migrating the DB: %s", err)
	}
}
