package ezex

import (
	"database/sql"
	"reflect"
)

// dbAdd handles insert queries and returns the new entity ID if successful
func dbAdd(db *sql.DB, query string, args ...any) (int, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// dbUpdate handles update queries and returns the number of affected rows
func dbUpdate(db *sql.DB, query string, args ...any) (int, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()

	return int(n), err
}

// dbDelete handles delete queries and returns the number of affected rows
func dbDelete(db *sql.DB, query string, args ...any) int {
	result, _ := db.Exec(query, args...)
	if result == nil {
		return 0
	}

	n, _ := result.RowsAffected()

	return int(n)
}

// dbGet returns a slice of entities given a query
func dbGet[T any](db *sql.DB, query string, args ...any) []T {
	var rows *sql.Rows
	if args == nil {
		rows, _ = db.Query(query)
	} else {
		rows, _ = db.Query(query, args...)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var mappedRows []T
	// Get row struct value = T
	rowsValue := reflect.ValueOf(&mappedRows).Elem()
	rowType := rowsValue.Type().Elem()

	// Get the number of exported fields in the struct representing the DB row (1 field = 1 column)
	columnValues := make([]any, rowType.NumField())

	for rows.Next() {
		// Create new T (row value)
		rowVal := reflect.New(rowType).Elem()

		// Map each T field to the scanned columnValues
		for i := 0; i < rowVal.NumField(); i++ {
			// Note: this works positionally since we're following the struct fields order
			// so the columns order in the query matter
			columnValues[i] = rowVal.Field(i).Addr().Interface()
		}

		_ = rows.Scan(columnValues...)

		// Add row to the list of returned rows
		rowsValue.Set(reflect.Append(rowsValue, rowVal))
	}

	return mappedRows
}
