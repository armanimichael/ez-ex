package main

import (
	"database/sql"
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
	transactions  []ezex.TransactionView
	newPayee      ezex.Payee
	newCategory   ezex.Category
	amountInCents int64
	err           error
}

type deleteTransactionMsg = struct {
	deletedID    int
	deletedIndex int
	err          error
}

type hideErrorMessageMsg = struct {
	id      int64
	message string
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
		if _, err := ezex.UpdateAccountBalance(db, transaction.AccountID, transaction.AmountInCents); err != nil {
			return createNewTransactionMsg{err: err}
		}

		now := time.Now()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)
		transactions := ezex.GetTransactions(db, transaction.AccountID, monthStart, monthEnd)

		return createNewTransactionMsg{
			transactions:  transactions,
			newPayee:      payee,
			newCategory:   category,
			amountInCents: transaction.AmountInCents,
			err:           nil,
		}
	}
}

func deleteAccountCmd(db *sql.DB, id int, index int) tea.Cmd {
	return func() tea.Msg {
		if _, err := ezex.DeleteAccount(db, id); err != nil {
			return deleteTransactionMsg{
				err: err,
			}
		}

		return deleteAccountMsg{
			deletedID:    id,
			deletedIndex: index,
			err:          nil,
		}
	}
}

func deleteTransactionCmd(db *sql.DB, accountID int, id int, amountInCents int64, index int) tea.Cmd {
	return func() tea.Msg {
		var err error

		if _, err = ezex.DeleteTransaction(db, id); err != nil {
			return deleteTransactionMsg{err: err}
		}
		if _, err = ezex.UpdateAccountBalance(db, accountID, amountInCents); err != nil {
			return createNewTransactionMsg{err: err}
		}

		return deleteTransactionMsg{
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

func hideErrorMessageCmd(id int64, message string) tea.Cmd {
	return tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
		return hideErrorMessageMsg{
			message: message,
			id:      id,
		}
	})
}
