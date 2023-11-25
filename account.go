package ezex

import "database/sql"

type Account struct {
	ID                    int
	Name                  string
	Description           sql.NullString
	InitialBalanceInCents int64
	BalanceInCents        int64
}

// AddAccount creates a new account and returns the new account ID if successful
func AddAccount(db *sql.DB, account Account) (int, error) {
	return dbAdd(
		db,
		`
		INSERT INTO accounts 	(name, description, initial_balance_in_cents, balance_in_cents)
		VALUES					($name, $description, $initial_balance_in_cents, $balance_in_cents)
		`,
		account.Name,
		account.Description,
		account.InitialBalanceInCents,
		account.BalanceInCents,
	)
}

func DeleteAccount(db *sql.DB, id int) int {
	return dbDelete(db, `DELETE FROM accounts WHERE id = $id`, id)
}

func UpdateAccount(db *sql.DB, account Account) (int, error) {
	return dbUpdate(
		db,
		`
		UPDATE	accounts
		SET		name 						= $name,
		     	description 				= $description,
		     	initial_balance_in_cents	= $initial_balance_in_cents,
		     	balance_in_cents			= $balance_in_cents
		WHERE	id = $id
		`,
		account.Name,
		account.Description,
		account.InitialBalanceInCents,
		account.BalanceInCents,
		account.ID,
	)
}

func GetAccounts(db *sql.DB) []Account {
	return dbGet[Account](db, `SELECT id, name, description, initial_balance_in_cents, balance_in_cents FROM accounts ORDER BY name DESC`)
}
