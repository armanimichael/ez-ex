package main

import (
	"database/sql"
	ezex "github.com/armanimichael/ez-ex"
	tea "github.com/charmbracelet/bubbletea"
)

type switchModelMsg = int

type createNewAccountMsg = struct {
	newAccount ezex.Account
	err        error
}

type deleteAccountMsg = struct {
	deletedID int
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

func switchModelCmd(id int) tea.Cmd {
	return func() tea.Msg {
		return id
	}
}
