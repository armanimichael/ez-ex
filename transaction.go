package ezex

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID                  int
	CategoryID          int
	PayeeID             int
	AccountID           int
	AmountInCents       int64
	TransactionDateUnix int64
	UpdateDateUnix      sql.NullInt64
	DeleteDateUnix      sql.NullInt64
	Notes               sql.NullString
}

type TransactionView struct {
	ID                  int
	CategoryID          int
	PayeeID             int
	AccountID           int
	AmountInCents       int64
	TransactionDateUnix int64
	UpdateDateUnix      sql.NullInt64
	DeleteDateUnix      sql.NullInt64
	Notes               sql.NullString
	CategoryName        string
	PayeeName           string
	AccountName         string
}

func AddTransaction(db *sql.DB, transaction Transaction) (int, error) {
	return dbAdd(
		db,
		`
		INSERT INTO transactions	(category_id, payee_id, account_id, amount_in_cents, transaction_date_unix, update_date_unix, delete_date_unix, notes)
		VALUES 						($category_id, $payee_id, $account_id, $amount_in_cents, $transaction_date_unix, $update_date_unix, $delete_date_unix, $notes)
		`,
		transaction.CategoryID,
		transaction.PayeeID,
		transaction.AccountID,
		transaction.AmountInCents,
		transaction.TransactionDateUnix,
		transaction.UpdateDateUnix,
		transaction.DeleteDateUnix,
		transaction.Notes,
	)
}

// DeleteTransaction soft-deletes the transaction and returns the number of affected rows
func DeleteTransaction(db *sql.DB, id int) (int, error) {
	return dbUpdate(
		db,
		`UPDATE transactions SET delete_date_unix = $date WHERE id = $id`,
		time.Now().Unix(),
		id,
	)
}

func UpdateTransaction(db *sql.DB, transaction Transaction) (int, error) {
	return dbUpdate(
		db,
		`
		UPDATE	transactions
		SET		category_id				= $category_id,
				payee_id 				= $payee_id,
				account_id 				= $account_id,
				amount_in_cents 		= $amount_in_cents,
				transaction_date_unix	= $transaction_date_unix,
				update_date_unix 		= $update_date_unix,
				delete_date_unix 		= $delete_date_unix,
				notes 					= $notes
		WHERE	id = $id
		`,
		transaction.CategoryID,
		transaction.PayeeID,
		transaction.AccountID,
		transaction.AmountInCents,
		transaction.TransactionDateUnix,
		transaction.UpdateDateUnix,
		transaction.DeleteDateUnix,
		transaction.Notes,
		transaction.ID,
	)
}

// GetTransactions returns a list of transaction for a given account between minDate and maxDate (excluded)
func GetTransactions(db *sql.DB, accountID int, minDate time.Time, maxDate time.Time) []TransactionView {
	return dbGet[TransactionView](
		db,
		`
		SELECT		t.id,
					t.category_id,
					t.payee_id,
					t.account_id,
					t.amount_in_cents,
					t.transaction_date_unix,
					t.update_date_unix,
					t.delete_date_unix,
					t.notes,
					c.name                      AS CategoryName,
					p.name                      AS PayeeName,
					a.name                      AS AccountName
		FROM		transactions t
		JOIN        accounts a
		ON          a.id = t.account_id
		JOIN        categories c
		ON          c.id = t.category_id
		JOIN        payees p
		ON          p.id = t.payee_id
		WHERE			t.account_id = $accountID
					AND	t.transaction_date_unix >= $minDateUnix
					AND t.transaction_date_unix < $maxDateUnix
					AND t.delete_date_unix IS NULL
		ORDER BY	t.transaction_date_unix DESC
		`,
		accountID,
		minDate.Unix(),
		maxDate.Unix(),
	)
}
