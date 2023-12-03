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
	deletedID int
}

type switchTransactionsMonthMsg = struct {
	month        time.Month
	year         int
	transactions []ezex.TransactionView
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

func deleteAccountCmd(db *sql.DB, id int) tea.Cmd {
	return func() tea.Msg {
		_ = ezex.DeleteAccount(db, id)

		return deleteAccountMsg{deletedID: id}
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
