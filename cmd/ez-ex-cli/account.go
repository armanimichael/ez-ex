package main

import (
	"database/sql"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
)

type accountModel struct {
	db        *sql.DB
	accounts  []ezex.Account
	textInput textinput.Model
	table     table.Model
}

func initAccountModel(db *sql.DB) (m accountModel) {
	m.db = db
	m.accounts = ezex.GetAccounts(db)

	ti := textinput.New()
	ti.Placeholder = "Create a new account"
	ti.Focus()

	t := table.New(
		table.WithColumns(
			[]table.Column{
				{Title: "ID", Width: 5},
				{Title: "Account name", Width: 20},
				{Title: "Balance", Width: 10},
				{Title: "Description", Width: 100},
			}),
		table.WithRows(accountsToTableRows(m.accounts...)),
		table.WithFocused(true),
		table.WithHeight(6),
	)
	t.SetStyles(tableStyle)

	m.table = t
	m.textInput = ti

	return m
}

func (m accountModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m accountModel) Update(msg tea.Msg) (accountModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.textInput.Value() != "" {
				account := ezex.Account{
					Name: m.textInput.Value(),
					Description: sql.NullString{
						String: "",
						Valid:  false,
					},
					InitialBalanceInCents: 0,
					BalanceInCents:        0,
				}

				id, _ := ezex.AddAccount(m.db, account)
				account.ID = id
				m.accounts = append(m.accounts, account)

				m.table.SetRows(append(m.table.Rows(), accountsToTableRows(account)[0]))
			}
		}
	}

	if len(m.accounts) > 0 {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m accountModel) View() string {
	if len(m.accounts) > 0 {
		return baseStyle.Render(m.table.View())
	}

	return m.textInput.View()
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
