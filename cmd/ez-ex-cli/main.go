package main

import (
	"database/sql"
	"flag"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	customLogger "github.com/armanimichael/ez-ex/cmd/ez-ex-cli/logger"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

var logger customLogger.Logger

func main() {
	dbName := flag.String(
		"db-name",
		ezex.DefaultDBName,
		"App DB filename (if not present, one will be created)",
	)
	logLevel := flag.Int(
		"log-level",
		5,
		"Application log level (trace = 0, debug = 1, info = 2, warn = 3, error = 4, fatal = 5, none = 6)",
	)
	flag.Parse()

	logger = customLogger.NewFileLogger(*logLevel)
	defer func(logger customLogger.Logger) {
		_ = logger.Close()
	}(logger)

	var opts []ezex.OptionsBuilder
	if *dbName != "" {
		opts = append(opts, ezex.WithDBName(*dbName))
	}

	db, err := ezex.OpenDB(opts...)
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
