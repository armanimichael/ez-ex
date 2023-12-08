package ezex

import (
	"database/sql"
)

type Payee struct {
	ID          int
	Name        string
	Description sql.NullString
}

func (p Payee) GetName() string {
	return p.Name
}

func AddPayee(db *sql.DB, payee Payee) (int, error) {
	return dbAdd(
		db,
		`INSERT INTO payees (name, description) VALUES ($name, $description)`,
		payee.Name,
		payee.Description,
	)
}

// DeletePayee deletes a payee, trying to delete a payee already in use will be noop
func DeletePayee(db *sql.DB, id int) int {
	return dbDelete(db, `DELETE FROM payees WHERE id = $id`, id)
}

func UpdatePayee(db *sql.DB, payee Payee) (int, error) {
	return dbUpdate(
		db,
		`
		UPDATE  payees
		SET     name 		= $name,
				description = $description
		WHERE   id = $id
		`,
		payee.Name,
		payee.Description,
		payee.ID,
	)
}

func GetPayees(db *sql.DB) []Payee {
	return dbGet[Payee](db, `SELECT id, name, description FROM payees ORDER BY id DESC`)
}
