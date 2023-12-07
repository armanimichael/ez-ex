package main

import (
	"database/sql"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"regexp"
	"strconv"
	"strings"
)

type accountCreatorModel struct {
	db                   *sql.DB
	stage                int
	newAccount           ezex.Account
	existingAccountNames map[string]struct{}
	input                standardTextInput
}

const (
	accountNewNameStage = iota
	accountNewDescriptionStage
	accountNewInitialBalanceStage
)

func initAccountCreatorModel(db *sql.DB, accounts []ezex.Account) accountCreatorModel {
	names := accountsToNamesMap(accounts)

	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 100
	ti.Placeholder = "Account name"
	ti.Focus()

	return accountCreatorModel{
		db:                   db,
		newAccount:           ezex.Account{},
		existingAccountNames: names,
		input: standardTextInput{
			model:         ti,
			previousInput: "",
			errorMsg:      "",
			label:         "Account name",
		},
	}
}

func (m accountCreatorModel) Update(msg tea.Msg) (accountCreatorModel, tea.Cmd) {
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
	}

	m.input.errorMsg = m.validateInput()
	m.input.previousInput = m.input.model.Value()
	m.input.model, cmd = m.input.model.Update(msg)

	return m, cmd
}

func (m accountCreatorModel) View() string {
	if m.input.errorMsg != "" {
		return m.input.label + ": " + m.input.model.View() + "\n" + errorMessageStyle.Render(m.input.errorMsg)
	}

	return m.input.label + ": " + m.input.model.View()
}

func (m *accountCreatorModel) reset(accounts []ezex.Account) {
	m.existingAccountNames = accountsToNamesMap(accounts)
	m.stage = accountNewNameStage
	m.input.label = "Account name"
	m.input.model.Placeholder = "Account name"
	m.input.model.SetValue("")
}

func (m accountCreatorModel) validateInput() string {
	input := m.input.model.Value()

	// Avoid multiple checks for same inputs (because of update)
	if input == m.input.previousInput && m.input.previousInput != "" {
		return m.input.errorMsg
	}

	switch m.stage {
	case accountNewInitialBalanceStage:
		// Strict balance format 0:00
		re := regexp.MustCompile(`^-?(?P<integer>\d+)(\.(?P<cents>\d{2}))+$`)
		isValidBalance := re.MatchString(input)

		if !isValidBalance {
			return "invalid balance format, should look like `n.xx`"
		}
	case accountNewNameStage:
		if input == "" {
			return "account name must have at least 1 char"
		}

		// Don't allow duplicate account names
		if _, exists := m.existingAccountNames[input]; exists {
			return fmt.Sprintf("there's already an account named: %v", input)
		}
	}

	return ""
}

func accountsToNamesMap(accounts []ezex.Account) map[string]struct{} {
	names := make(map[string]struct{}, len(accounts))
	for _, account := range accounts {
		names[account.Name] = struct{}{}
	}
	return names
}
