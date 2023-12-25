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

const DefaultDBName = "user-data.db"
const UserDataDir = ".ez-ex"

type appOptions struct {
	dbName string
}

type OptionsBuilder = func(*appOptions)

// WithDBName sets user-data DB name
func WithDBName(name string) OptionsBuilder {
	return func(o *appOptions) {
		o.dbName = name
	}
}

func OpenDB(opts ...OptionsBuilder) (*sql.DB, error) {
	options := appOptions{
		dbName: DefaultDBName,
	}
	for _, option := range opts {
		option(&options)
	}

	home, _ := os.UserHomeDir()
	_ = os.MkdirAll(path.Join(home, UserDataDir), 0700)

	dsName := fmt.Sprintf("file:%s/.ez-ex/%s?_foreign_keys=true", home, options.dbName)
	return sql.Open("sqlite3", dsName)
}

func MigrateDB(db *sql.DB) error {
	_, err := db.Exec(dbInitScriptSQL)
	return err
}
