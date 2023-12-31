package ezex

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

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

func DeleteAccount(db *sql.DB, id int) (int, error) {
	return dbUpdate(
		db,
		`UPDATE accounts SET delete_date_unix = $date WHERE id = $id`,
		time.Now().Unix(),
		id,
	)
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

func UpdateAccountBalance(db *sql.DB, accountID int, amountInCents int64) (int, error) {
	return dbUpdate(
		db,
		`
		UPDATE	accounts
		SET		balance_in_cents = balance_in_cents - $amount_in_cents
		WHERE	id = $id
		`,
		amountInCents,
		accountID,
	)
}

func GetAccounts(db *sql.DB) []Account {
	return dbGet[Account](
		db,
		`
		SELECT		id,
					name,
					description,
					initial_balance_in_cents,
					balance_in_cents
		FROM 		accounts
		WHERE		delete_date_unix IS NULL
		ORDER BY 	id DESC`,
	)
}

func GetAccount(db *sql.DB, id int) (Account, error) {
	results := dbGet[Account](
		db,
		`
		SELECT		id,
					name,
					description,
					initial_balance_in_cents,
					balance_in_cents
		FROM		accounts
		WHERE		id = $id
		ORDER BY 	id DESC`,
		id,
	)

	if len(results) == 0 {
		return Account{}, errors.New(fmt.Sprintf("no accounts with id: %d", id))
	}

	return results[0], nil
}
