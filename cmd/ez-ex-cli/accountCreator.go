package main

import (
	"database/sql"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type accountCreatorModel struct {
	db                   *sql.DB
	stage                int
	existingAccountNames map[string]struct{}
	inputs               []standardTextInput
}

const (
	accountNewNameStage = iota
	accountNewDescriptionStage
	accountNewInitialBalanceStage
)

func initAccountCreatorModel(db *sql.DB, accounts []ezex.Account) accountCreatorModel {
	names := accountsToNamesMap(accounts)

	inputs := make([]standardTextInput, 3)
	inputs[accountNewNameStage] = createAccountInput(accountNewNameStage)
	inputs[accountNewDescriptionStage] = createAccountInput(accountNewDescriptionStage)
	inputs[accountNewInitialBalanceStage] = createAccountInput(accountNewInitialBalanceStage)

	return accountCreatorModel{
		db:                   db,
		existingAccountNames: names,
		inputs:               inputs,
	}
}

func (m accountCreatorModel) Update(msg tea.Msg) (accountCreatorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			isValid := true
			for _, input := range m.inputs {
				if input.errorMsg != "" {
					// Errors, cannot submit
					isValid = false
					break
				}
			}
			if !isValid {
				break
			}

			description := m.inputs[accountNewDescriptionStage].model.Value()
			balance := decodeCents(m.inputs[accountNewInitialBalanceStage].model.Value())

			return m, createNewAccountCmd(m.db, ezex.Account{
				Name: m.inputs[accountNewNameStage].model.Value(),
				Description: sql.NullString{
					String: description,
					Valid:  description != "",
				},
				InitialBalanceInCents: balance,
				BalanceInCents:        balance,
			})
		case "up", "down":
			return m.switchAccount(msg)
		}
	}

	currentInput := &m.inputs[m.stage]
	for i := range m.inputs {
		errMsg := m.validateInput(i)
		m.inputs[i].previousInput = m.inputs[i].model.Value()
		m.inputs[i].errorMsg = errMsg
	}
	currentInput.model, cmd = currentInput.model.Update(msg)

	return m, cmd
}

func (m accountCreatorModel) View() string {
	lastIndex := len(m.inputs) - 1
	inputListStr := strings.Builder{}
	errorsListStr := strings.Builder{}

	for i, input := range m.inputs {
		var render func(s ...string) string
		if i == m.stage {
			render = inputBoxSelectedStyle.Render
		} else {
			render = inputBoxStyle.Render
		}

		var mark string
		if input.errorMsg == "" {
			mark = successMessageStyle.Render("🗸\t")
		} else {
			mark = errorMessageStyle.Render("🞪\t")
		}
		label := mark + fmt.Sprintf("%-19s", input.label+": ")

		inputListStr.WriteString(render(label) + input.model.View())

		if errMsg := input.errorMsg; errMsg != "" {
			errorsListStr.WriteString(errorMessageStyle.Render("- "+errMsg) + "\n")
		}

		if i != lastIndex {
			inputListStr.WriteString("\n")
		}
	}

	if errorsListStr.Len() > 0 {
		return inputListStr.String() + "\n\n" + errorsListStr.String()
	} else {
		return inputListStr.String()
	}
}

func (m accountCreatorModel) switchAccount(msg fmt.Stringer) (accountCreatorModel, tea.Cmd) {
	m.inputs[m.stage].model.Blur()

	if msg.String() == "up" {
		if m.stage > accountNewNameStage {
			m.stage--
		} else {
			m.stage = accountNewInitialBalanceStage
		}
	} else {
		if m.stage < accountNewInitialBalanceStage {
			m.stage++
		} else {
			m.stage = accountNewNameStage
		}
	}

	m.inputs[m.stage].model.SetCursor(0)
	m.inputs[m.stage].model.Focus()

	return m, textinput.Blink
}

func (m accountCreatorModel) reset(accounts []ezex.Account) accountCreatorModel {
	m.stage = accountNewNameStage
	m.existingAccountNames = accountsToNamesMap(accounts)

	m.inputs[accountNewNameStage] = createAccountInput(accountNewNameStage)
	m.inputs[accountNewDescriptionStage] = createAccountInput(accountNewDescriptionStage)
	m.inputs[accountNewInitialBalanceStage] = createAccountInput(accountNewInitialBalanceStage)

	return m
}

func (m accountCreatorModel) validateInput(stage int) string {
	currentInput := &m.inputs[stage]
	value := currentInput.model.Value()

	// Avoid multiple checks for same inputs (because of update)
	if value == currentInput.previousInput && currentInput.previousInput != "" {
		return currentInput.errorMsg
	}

	switch stage {
	case accountNewNameStage:
		if value == "" {
			return "account name must have at least 1 char"
		}

		// Don't allow duplicate account names
		if _, exists := m.existingAccountNames[value]; exists {
			return fmt.Sprintf("there's already an account named: %v", value)
		}
	case accountNewInitialBalanceStage:
		if !moneyFormatRegex.MatchString(value) {
			return "invalid balance format, should look like `0.00` or `-0.00`"
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

func createAccountInput(stage int) standardTextInput {
	ti := textinput.New()
	ti.Prompt = ""

	switch stage {
	case accountNewNameStage:
		ti.Placeholder = "..."
		ti.SetValue("")
		ti.Focus()

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Account name*",
		}
	case accountNewDescriptionStage:
		ti.Placeholder = "<NO DESCRIPTION>"

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Description",
		}
	case accountNewInitialBalanceStage:
		ti.Placeholder = "0.00"

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Balance*",
		}
	}

	panic("unsupported account creation stage")
}
