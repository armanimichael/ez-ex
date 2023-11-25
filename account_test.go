package ezex

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddAccount(t *testing.T) {
	id1, err1 := AddAccount(testDB, Account{
		Name:                  "TestAddAccount1",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})
	id2, err2 := AddAccount(testDB, Account{
		Name:                  "TestAddAccount2",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})

	assert.Greater(t, id1, 0)
	assert.Nil(t, err1)
	assert.Greater(t, id2, 0)
	assert.Nil(t, err2)
}

func TestAddAccount_Unique(t *testing.T) {
	_, _ = AddAccount(testDB, Account{
		Name:                  "TestAddAccount_Unique",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})
	_, err := AddAccount(testDB, Account{
		Name:                  "TestAddAccount_Unique",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})

	assert.Error(t, err)
}

func TestDeleteAccount(t *testing.T) {
	id, _ := AddAccount(testDB, Account{
		Name:                  "TestDeleteAccount",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})
	n := DeleteAccount(testDB, id)

	assert.Greater(t, n, 0)
}

func TestUpdateAccount(t *testing.T) {
	id, _ := AddAccount(testDB, Account{
		Name:                  "TestUpdateAccount",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})
	n, err := UpdateAccount(testDB, Account{
		ID:                    id,
		Name:                  "TestUpdateAccount2",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})

	assert.Greater(t, n, 0)
	assert.Nil(t, err)
}

func TestUpdateAccount_Unique(t *testing.T) {
	_, _ = AddAccount(testDB, Account{
		Name:                  "TestUpdateAccount_Unique",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})
	id, _ := AddAccount(testDB, Account{
		Name:                  "TestUpdateAccount_Unique2",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})
	n, err := UpdateAccount(testDB, Account{
		ID:                    id,
		Name:                  "TestUpdateAccount_Unique",
		Description:           sql.NullString{},
		InitialBalanceInCents: 0,
		BalanceInCents:        0,
	})

	assert.Equal(t, n, 0)
	assert.Error(t, err)
}

func TestGetAccounts(t *testing.T) {
	account1 := Account{
		Name: "TestGetAccounts1",
		Description: sql.NullString{
			String: "TestGetAccounts1",
			Valid:  true,
		},
		InitialBalanceInCents: 123,
		BalanceInCents:        1234,
	}
	account2 := Account{
		Name:                  "TestGetAccount2",
		Description:           sql.NullString{},
		InitialBalanceInCents: -123,
		BalanceInCents:        -1234,
	}

	_, _ = AddAccount(testDB, account1)
	_, _ = AddAccount(testDB, account2)

	accounts := GetAccounts(testDB)

	type AccountWithoutID struct {
		name                  string
		description           sql.NullString
		initialBalanceInCents int64
		balanceInCents        int64
	}
	accountsWithoutID := make([]AccountWithoutID, len(accounts))
	for i, account := range accounts {
		accountsWithoutID[i] = AccountWithoutID{
			name:                  account.Name,
			description:           account.Description,
			initialBalanceInCents: account.InitialBalanceInCents,
			balanceInCents:        account.BalanceInCents,
		}
	}

	assert.Contains(t, accountsWithoutID, AccountWithoutID{
		name:                  account1.Name,
		description:           account1.Description,
		initialBalanceInCents: account1.InitialBalanceInCents,
		balanceInCents:        account1.BalanceInCents,
	})
	assert.Contains(t, accountsWithoutID, AccountWithoutID{
		name:                  account2.Name,
		description:           account2.Description,
		initialBalanceInCents: account2.InitialBalanceInCents,
		balanceInCents:        account2.BalanceInCents,
	})
}
