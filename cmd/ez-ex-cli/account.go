package main

import (
	"database/sql"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/armanimichael/ez-ex/internal/slice"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
	"time"
)

type accountModel struct {
	db             *sql.DB
	stage          int
	accounts       []ezex.Account
	accountCreator accountCreatorModel
	err            struct {
		id  int64
		msg string
	}
	table struct {
		model      table.Model
		selectedID int
	}
}

const (
	accountSelectionStage = iota
	accountCreationStage
)

var accountTableKeySuggestions = formatKeySuggestions([][]string{
	{"^C", "quit"},
	{"{enter}", "select account"},
	{"d", "delete account"},
	{"n", "create account"},
})

func initAccountModel(db *sql.DB) (m accountModel) {
	m.db = db
	m = m.createAccountsTable(ezex.GetAccounts(db))

	if len(m.accounts) > 0 {
		m.table.selectedID = m.accounts[0].ID
		m.stage = accountSelectionStage
	} else {
		m.stage = accountCreationStage
	}

	m.accountCreator = initAccountCreatorModel(db, m.accounts)

	return m
}

func (m accountModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m accountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case hideErrorMessageMsg:
		if msg.message == m.err.msg && msg.id == m.err.id {
			m.err.msg = ""
		}
	case deleteAccountMsg:
		if msg.err != nil {
			logger.Err(fmt.Sprintf("Error deleting account: %v", msg.err))
			m.err.msg = msg.err.Error()
			m.err.id = time.Now().UnixMicro()

			return m, hideErrorMessageCmd(m.err.id, m.err.msg)
		}

		logger.Debug(fmt.Sprintf("Delete account (ID: %v)", m.table.selectedID))

		// Remove deleted row from the table
		updatedAccounts := make([]ezex.Account, 0, len(m.accounts)-1)
		for _, account := range m.accounts {
			if account.ID != msg.deletedID {
				updatedAccounts = append(updatedAccounts, account)
			}
		}
		m.accounts = updatedAccounts
		m.accountCreator.reset(m.accounts)

		if len(m.accounts) == 0 {
			// No accounts left, go back to creation
			m.stage = accountCreationStage
		} else {
			// Pre-select fist row
			m.table.selectedID = m.accounts[0].ID
			m.table.model.SetCursor(0)

			m.table.model.SetRows(accountsToTableRows(m.accounts...))
		}
		m.table.model.GotoTop()
	case createNewAccountMsg:
		if msg.err != nil {
			logger.Fatal(fmt.Sprintf("Error creating account: %v", msg.err))
			return m, tea.Quit
		}

		logger.Debug(fmt.Sprintf("Create new account: %v", msg.newAccount))
		m.accounts = slice.Prepend(m.accounts, msg.newAccount, 10)
		m.accountCreator.reset(m.accounts)
		m.table.model.SetRows(accountsToTableRows(m.accounts...))
		m.table.selectedID = msg.newAccount.ID
		m.stage = accountSelectionStage
		m.table.model.GotoTop()
	}

	if m.stage == accountCreationStage {
		m.accountCreator, cmd = m.accountCreator.Update(msg)

		return m, cmd
	}

	return m.handleAccountSelectionCommands(msg)
}

func (m accountModel) View() string {
	if m.stage == accountSelectionStage {
		msg := ""
		if m.err.msg != "" {
			msg = errorMessageStyle.Render("Error: "+m.err.msg) + "\n"
		}

		return baseStyle.Render(m.table.model.View()) + "\n" + accountTableKeySuggestions + "\n" + msg
	}

	return m.accountCreator.View()
}

func (m accountModel) createAccountsTable(accounts []ezex.Account) accountModel {
	m.accounts = accounts
	m.table.model = createStandardTable(
		[]table.Column{
			{Title: "ID", Width: 5},
			{Title: "Account name", Width: 20},
			{Title: "Balance", Width: 10},
			{Title: "Description", Width: 50},
		},
		accountsToTableRows(accounts...),
	)

	return m
}

func (m accountModel) handleAccountSelectionCommands(msg tea.Msg) (accountModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table.model, cmd = m.table.model.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			logger.Debug(fmt.Sprintf("Select account ID %v", m.table.selectedID))
			return m, switchModelCmd(transactionModelID, m.table.selectedID)
		case "d":
			var selectedID int
			if len(m.accounts) == 1 {
				selectedID = m.accounts[0].ID
			} else {
				selectedID = m.table.selectedID
			}

			return m, tea.Batch(deleteAccountCmd(m.db, selectedID, m.table.model.Cursor()), cmd)
		case "n":
			m.stage = accountCreationStage
			return m, textinput.Blink
		case "down", "up":
			r := m.table.model.SelectedRow()
			selectedID, _ := strconv.ParseInt(r[0], 10, 32)
			m.table.selectedID = int(selectedID)
			m.err.msg = ""
		}

	}

	return m, cmd
}
