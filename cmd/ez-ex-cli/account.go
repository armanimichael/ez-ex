package main

import (
	"database/sql"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"regexp"
	"strconv"
	"strings"
)

type accountModel struct {
	db         *sql.DB
	stage      int
	newAccount ezex.Account
	accounts   []ezex.Account
	table      struct {
		model      table.Model
		selectedID int
	}
	input struct {
		model         textinput.Model
		previousInput string
		errorMsg      string
		label         string
	}
}

const (
	accountSelectionStage = iota
	accountNewNameStage
	accountNewDescriptionStage
	accountNewInitialBalanceStage
)

var accountTableKeySuggestions = buildAccountTableKeySuggestions()

func initAccountModel(db *sql.DB) (m accountModel) {
	m.db = db
	m = m.createAccountsTable(ezex.GetAccounts(db))

	if len(m.accounts) > 0 {
		m.stage = accountSelectionStage
	} else {
		m.stage = accountNewNameStage
	}

	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 100
	ti.Placeholder = "Account name"
	m.input.label = "Account name"
	ti.Focus()

	m.input.model = ti

	return m
}

func (m accountModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m accountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.stage {
	case accountSelectionStage:
		return m.handleAccountSelection(msg)
	default:
		return m.handleAccountCreation(msg)
	}
}

func (m accountModel) View() string {
	if m.stage == accountSelectionStage {
		return baseStyle.Render(m.table.model.View()) + "\n" + accountTableKeySuggestions + "\n"
	}

	if m.input.errorMsg != "" {
		return m.input.label + ": " + m.input.model.View() + "\n" + errorMessageStyle.Render(m.input.errorMsg)
	}

	return m.input.label + ": " + m.input.model.View()
}

func (m accountModel) createAccountsTable(accounts []ezex.Account) accountModel {
	m.newAccount = ezex.Account{}
	m.accounts = accounts
	t := table.New(
		table.WithColumns(
			[]table.Column{
				{Title: "ID", Width: 10},
				{Title: "Account name", Width: 20},
				{Title: "Balance", Width: 10},
				{Title: "Description", Width: 50},
			}),
		table.WithRows(accountsToTableRows(accounts...)),
		table.WithFocused(true),
		table.WithHeight(6),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(foreground).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(selectedForeground).
		Background(selectedBackground).
		Bold(false)
	t.SetStyles(s)

	m.table.model = t
	return m
}

func buildAccountTableKeySuggestions() string {
	commands := [][]string{
		{"^C", "quit"},
		{"{enter}", "select account"},
		{"d", "delete account"},
		{"n", "create account"},
	}

	str := strings.Builder{}
	for _, pair := range commands {
		str.WriteString(
			fmt.Sprintf(
				"%s\t\t%s\n",
				keySuggestionStyle.Render(pair[0]),
				keySuggestionNoteStyle.Render(pair[1]),
			),
		)
	}

	return str.String()
}

func (m accountModel) validateInput() string {
	input := m.input.model.Value()

	// Avoid multiple checks for same inputs (because of update)
	if input == m.input.previousInput && m.input.previousInput != "" {
		return m.input.errorMsg
	}

	switch m.stage {
	case accountNewInitialBalanceStage:
		// Strict balance format 0:00
		re := regexp.MustCompile(`(?P<integer>\d+)(\.(?P<cents>\d{2}))+$`)
		isValidBalance := re.MatchString(input)

		if !isValidBalance {
			return "invalid balance format, should look like `n.xx`"
		}
	case accountNewNameStage:
		if input == "" {
			return "account name must have at least 1 char"
		}

		// Don't allow duplicate account names
		for _, account := range m.accounts {
			if input == account.Name {
				return fmt.Sprintf("there's already an account named: %v", input)
			}
		}
	}

	return ""
}

func (m accountModel) handleAccountSelection(msg tea.Msg) (accountModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table.model, cmd = m.table.model.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			logger.Debug("Switch to transaction model")
			return m, switchModelCmd(transactionModelID)
		case "d":
			var selectedID int
			if len(m.accounts) == 1 {
				selectedID = m.accounts[0].ID
			} else {
				selectedID = m.table.selectedID
			}

			return m, tea.Batch(deleteAccountCmd(m.db, selectedID), cmd)
		case "n":
			m.input.model.Placeholder = "Account name"
			m.input.label = "Account name"
			m.stage = accountNewNameStage

			return m, tea.Batch(m.input.model.Focus(), cmd)
		case "down", "up":
			r := m.table.model.SelectedRow()
			selectedID, _ := strconv.ParseInt(r[0], 10, 32)
			m.table.selectedID = int(selectedID)
		}
	case deleteAccountMsg:
		logger.Debug(fmt.Sprintf("Delete account (ID: %v)", m.table.selectedID))
		updatedAccounts := make([]ezex.Account, 0, len(m.accounts)-1)

		for _, account := range m.accounts {
			if account.ID != msg.deletedID {
				updatedAccounts = append(updatedAccounts, account)
			}
		}

		if len(updatedAccounts) == 0 {
			m.input.model.Placeholder = "Account name"
			m.input.label = "Account name"
			m.stage = accountNewNameStage

			return m.createAccountsTable(updatedAccounts), textinput.Blink
		}

		return m.createAccountsTable(updatedAccounts), cmd
	}

	return m, cmd
}

func (m accountModel) handleAccountCreation(msg tea.Msg) (accountModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.input.errorMsg != "" {
				break
			}

			value := m.input.model.Value()
			m.input.model.SetValue("")

			switch m.stage {
			case accountNewNameStage:
				m.newAccount.Name = value
				m.input.label = "Account description (can be empty)"
				m.input.model.Placeholder = ""
				m.stage = accountNewDescriptionStage
			case accountNewDescriptionStage:
				m.newAccount.Description = sql.NullString{
					String: value,
					Valid:  len(value) > 0,
				}
				m.input.label = "Account balance"
				m.input.model.Placeholder = "0.00"
				m.stage = accountNewInitialBalanceStage
				m.input.model.SetValue("0.00")
			case accountNewInitialBalanceStage:
				balanceInCentsStr := strings.Replace(value, ".", "", 1)
				balanceInCents, _ := strconv.ParseInt(balanceInCentsStr, 10, 64)
				m.newAccount.InitialBalanceInCents = balanceInCents
				m.newAccount.BalanceInCents = balanceInCents

				return m, createNewAccountCmd(m.db, m.newAccount)
			}
		}
	case createNewAccountMsg:
		if msg.err != nil {
			logger.Fatal(fmt.Sprintf("Error creating account: %v", msg.err))
			break
		}

		logger.Debug(fmt.Sprintf("Create new account: %v", msg.newAccount))
		m.newAccount = msg.newAccount
		m.input.model.SetValue("")
		m.accounts = append(m.accounts, m.newAccount)
		m.table.model.SetRows(append(m.table.model.Rows(), accountsToTableRows(m.newAccount)[0]))
		m.stage = accountSelectionStage
	}

	m.input.errorMsg = m.validateInput()
	m.input.previousInput = m.input.model.Value()
	m.input.model, cmd = m.input.model.Update(msg)

	return m, cmd
}

func accountsToTableRows(accounts ...ezex.Account) []table.Row {
	var rows []table.Row

	for _, account := range accounts {
		balance := strconv.FormatFloat(float64(account.BalanceInCents)/100.0, 'f', 2, 64)

		desc := account.Description.String
		if !account.Description.Valid {
			desc = "<NO DESCRIPTION>"
		}

		rows = append(
			rows,
			table.Row{
				strconv.Itoa(account.ID),
				account.Name,
				balance,
				desc,
			})
	}

	return rows
}
