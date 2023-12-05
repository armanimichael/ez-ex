package ezex

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddTransaction(t *testing.T) {
	id1, err1 := AddTransaction(testDB, Transaction{
		CategoryID:          0,
		PayeeID:             0,
		AccountID:           0,
		AmountInCents:       0,
		TransactionDateUnix: 0,
		UpdateDateUnix:      sql.NullInt64{},
		DeleteDateUnix:      sql.NullInt64{},
		Notes:               sql.NullString{},
	})
	id2, err2 := AddTransaction(testDB, Transaction{
		CategoryID:          0,
		PayeeID:             0,
		AccountID:           0,
		AmountInCents:       0,
		TransactionDateUnix: 0,
		UpdateDateUnix:      sql.NullInt64{},
		DeleteDateUnix:      sql.NullInt64{},
		Notes:               sql.NullString{},
	})

	assert.Nil(t, err1)
	assert.Greater(t, id1, 0)
	assert.Nil(t, err2)
	assert.Greater(t, id2, 0)
}

func TestDeleteTransaction(t *testing.T) {
	id, _ := AddTransaction(testDB, Transaction{
		CategoryID:          0,
		PayeeID:             0,
		AccountID:           0,
		AmountInCents:       0,
		TransactionDateUnix: 0,
		UpdateDateUnix:      sql.NullInt64{},
		DeleteDateUnix:      sql.NullInt64{},
		Notes:               sql.NullString{},
	})
	n, err := DeleteTransaction(testDB, id)

	assert.Greater(t, n, 0)
	assert.Nil(t, err)
}

func TestUpdateTransaction(t *testing.T) {
	id, _ := AddTransaction(testDB, Transaction{
		CategoryID:          0,
		PayeeID:             0,
		AccountID:           0,
		AmountInCents:       0,
		TransactionDateUnix: 0,
		UpdateDateUnix:      sql.NullInt64{},
		DeleteDateUnix:      sql.NullInt64{},
		Notes:               sql.NullString{},
	})
	n, err := UpdateTransaction(testDB, Transaction{
		ID:                  id,
		CategoryID:          1,
		PayeeID:             1,
		AccountID:           1,
		AmountInCents:       1,
		TransactionDateUnix: 1,
		UpdateDateUnix: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
		DeleteDateUnix: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
		Notes: sql.NullString{
			String: "note",
			Valid:  true,
		},
	})

	assert.Nil(t, err)
	assert.Greater(t, n, 0)
}

func TestGetTransactions(t *testing.T) {
	payee1ID, _ := AddPayee(testDB, Payee{
		Name:        "Payee1",
		Description: sql.NullString{},
	})
	payee2ID, _ := AddPayee(testDB, Payee{
		Name:        "Payee2",
		Description: sql.NullString{},
	})
	account1ID, _ := AddAccount(testDB, Account{
		Name:                  "Account1",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})
	account2ID, _ := AddAccount(testDB, Account{
		Name:                  "Account2",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})

	transaction1 := Transaction{
		CategoryID:          0,
		PayeeID:             payee1ID,
		AccountID:           account1ID,
		AmountInCents:       4,
		TransactionDateUnix: time.Date(2023, 11, 26, 0, 0, 0, 0, time.UTC).Unix(),
		UpdateDateUnix: sql.NullInt64{
			Int64: 6,
			Valid: true,
		},
		DeleteDateUnix: sql.NullInt64{},
		Notes: sql.NullString{
			String: "note",
			Valid:  true,
		},
	}
	transaction2 := Transaction{
		CategoryID:          0,
		PayeeID:             payee1ID,
		AccountID:           account1ID,
		AmountInCents:       4,
		TransactionDateUnix: time.Date(2023, 11, 28, 0, 0, 0, 0, time.UTC).Unix(),
		UpdateDateUnix: sql.NullInt64{
			Int64: 6,
			Valid: true,
		},
		DeleteDateUnix: sql.NullInt64{},
		Notes: sql.NullString{
			String: "note",
			Valid:  true,
		},
	}
	transaction3 := Transaction{
		CategoryID:          0,
		PayeeID:             payee2ID,
		AccountID:           account2ID,
		AmountInCents:       -4,
		TransactionDateUnix: -5,
		UpdateDateUnix:      sql.NullInt64{},
		DeleteDateUnix:      sql.NullInt64{},
		Notes:               sql.NullString{},
	}
	transaction4 := Transaction{
		CategoryID:          0,
		PayeeID:             payee2ID,
		AccountID:           account2ID,
		AmountInCents:       -4,
		TransactionDateUnix: -5,
		UpdateDateUnix:      sql.NullInt64{},
		DeleteDateUnix:      sql.NullInt64{},
		Notes:               sql.NullString{},
	}
	transactionDeleted := Transaction{
		CategoryID:          0,
		PayeeID:             payee2ID,
		AccountID:           account2ID,
		AmountInCents:       -4,
		TransactionDateUnix: -5,
		UpdateDateUnix:      sql.NullInt64{},
		DeleteDateUnix: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
		Notes: sql.NullString{},
	}

	_, _ = AddTransaction(testDB, transaction1)
	_, _ = AddTransaction(testDB, transaction2)
	_, _ = AddTransaction(testDB, transaction3)
	_, _ = AddTransaction(testDB, transaction4)
	_, _ = AddTransaction(testDB, transactionDeleted)

	transactions := GetTransactions(
		testDB,
		account1ID,
		time.Date(2023, 11, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 11, 28, 0, 0, 0, 0, time.UTC),
	)

	type TransactionWithoutID struct {
		categoryID          int
		payeeID             int
		accountID           int
		amountInCents       int64
		transactionDateUnix int64
		updateDateUnix      sql.NullInt64
		deleteDateUnix      sql.NullInt64
		notes               sql.NullString
		categoryName        string
		payeeName           string
		accountName         string
	}
	transactionsWithoutID := make([]TransactionWithoutID, len(transactions))
	for i, transaction := range transactions {
		transactionsWithoutID[i] = TransactionWithoutID{
			categoryID:          transaction.CategoryID,
			payeeID:             transaction.PayeeID,
			accountID:           transaction.AccountID,
			amountInCents:       transaction.AmountInCents,
			transactionDateUnix: transaction.TransactionDateUnix,
			updateDateUnix:      transaction.UpdateDateUnix,
			deleteDateUnix:      transaction.DeleteDateUnix,
			notes:               transaction.Notes,
			categoryName:        transaction.CategoryName,
			payeeName:           transaction.PayeeName,
			accountName:         transaction.AccountName,
		}
	}

	assert.Contains(t, transactionsWithoutID, TransactionWithoutID{
		categoryID:          transaction1.CategoryID,
		payeeID:             transaction1.PayeeID,
		accountID:           transaction1.AccountID,
		amountInCents:       transaction1.AmountInCents,
		transactionDateUnix: transaction1.TransactionDateUnix,
		updateDateUnix:      transaction1.UpdateDateUnix,
		deleteDateUnix:      transaction1.DeleteDateUnix,
		notes:               transaction1.Notes,
		categoryName:        "no category",
		payeeName:           "Payee1",
		accountName:         "Account1",
	})
	assert.NotContains(t, transactionsWithoutID, transactionDeleted)
}
