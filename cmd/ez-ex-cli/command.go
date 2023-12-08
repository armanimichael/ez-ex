package main

import (
	"database/sql"
	"errors"
	ezex "github.com/armanimichael/ez-ex"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type switchModelMsg = struct {
	modelID   int
	accountID int
}

type createNewAccountMsg = struct {
	newAccount ezex.Account
	err        error
}

type deleteAccountMsg = struct {
	deletedID    int
	deletedIndex int
	err          error
}

type switchTransactionsMonthMsg = struct {
	month        time.Month
	year         int
	transactions []ezex.TransactionView
}

type createNewTransactionMsg = struct {
	transactions []ezex.TransactionView
	newPayee     ezex.Payee
	newCategory  ezex.Category
	err          error
}

func createNewAccountCmd(db *sql.DB, account ezex.Account) tea.Cmd {
	return func() tea.Msg {
		id, err := ezex.AddAccount(db, account)
		account.ID = id

		return createNewAccountMsg{
			newAccount: account,
			err:        err,
		}
	}
}

func createNewTransactionCmd(db *sql.DB, transaction ezex.Transaction, payee ezex.Payee, category ezex.Category) tea.Cmd {
	return func() tea.Msg {
		if payee.ID == 0 {
			id, err := ezex.AddPayee(db, payee)
			if err != nil {
				return createNewTransactionMsg{err: err}
			}
			payee.ID = id
			transaction.PayeeID = id
		}
		if category.ID == 0 && category.Name != "" {
			id, err := ezex.AddCategory(db, category)
			if err != nil {
				return createNewTransactionMsg{err: err}
			}
			category.ID = id
			transaction.CategoryID = id
		}

		if _, err := ezex.AddTransaction(db, transaction); err != nil {
			return createNewTransactionMsg{err: err}
		}

		now := time.Now()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)
		transactions := ezex.GetTransactions(db, transaction.AccountID, monthStart, monthEnd)

		return createNewTransactionMsg{
			transactions: transactions,
			newPayee:     payee,
			newCategory:  category,
			err:          nil,
		}
	}
}

func deleteAccountCmd(db *sql.DB, id int, index int) tea.Cmd {
	return func() tea.Msg {
		n := ezex.DeleteAccount(db, id)

		var err error
		if n == 0 {
			err = errors.New("cannot delete account with active transactions")
		}

		return deleteAccountMsg{
			deletedID:    id,
			deletedIndex: index,
			err:          err,
		}
	}
}

func switchModelCmd(id int, accountID int) tea.Cmd {
	return func() tea.Msg {
		return switchModelMsg{
			modelID:   id,
			accountID: accountID,
		}
	}
}

func switchTransactionsMonthCmd(db *sql.DB, accountID int, year int, month time.Month) tea.Cmd {
	return func() tea.Msg {
		monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		monthEnd := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
		transactions := ezex.GetTransactions(db, accountID, monthStart, monthEnd)

		return switchTransactionsMonthMsg{
			month:        monthStart.Month(),
			year:         monthStart.Year(),
			transactions: transactions,
		}
	}
}
