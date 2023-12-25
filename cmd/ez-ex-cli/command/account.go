package command

import (
	"database/sql"
	ezex "github.com/armanimichael/ez-ex"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateNewAccountMsg = struct {
	NewAccount ezex.Account
	Err        error
}

type DeleteAccountMsg = struct {
	DeletedID    int
	DeletedIndex int
	Err          error
}

func CreateNewAccountCmd(db *sql.DB, account ezex.Account) tea.Cmd {
	return func() tea.Msg {
		id, err := ezex.AddAccount(db, account)
		account.ID = id

		return CreateNewAccountMsg{
			NewAccount: account,
			Err:        err,
		}
	}
}
func DeleteAccountCmd(db *sql.DB, id int, index int) tea.Cmd {
	return func() tea.Msg {
		if _, err := ezex.DeleteAccount(db, id); err != nil {
			return DeleteTransactionMsg{
				Err: err,
			}
		}

		return DeleteAccountMsg{
			DeletedID:    id,
			DeletedIndex: index,
			Err:          nil,
		}
	}
}
