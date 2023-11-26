package main

import (
	"database/sql"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
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

	p := tea.NewProgram(initialModel(db))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error running program: %v", err)
		os.Exit(1)
	}
}
