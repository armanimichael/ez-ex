package main

import (
	"database/sql"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type transactionCreatorModel struct {
	db         *sql.DB
	stage      int
	accountID  int
	payees     []ezex.Payee
	categories []ezex.Category
	inputs     []standardTextInput
	suggestion struct {
		autocompleteSuggestion string
		payee                  ezex.Payee
		category               ezex.Category
	}
}

const (
	transactionDateStage = iota
	transactionAmountStage
	transactionPayeeStage
	transactionCategoryStage
	transactionNoteStage
)

func initTransactionCreator(
	db *sql.DB,
	accountID int,
	payees []ezex.Payee,
	categories []ezex.Category,
) transactionCreatorModel {
	inputs := make([]standardTextInput, 5)
	inputs[transactionDateStage] = createTransactionInput(transactionDateStage)
	inputs[transactionAmountStage] = createTransactionInput(transactionAmountStage)
	inputs[transactionPayeeStage] = createTransactionInput(transactionPayeeStage)
	inputs[transactionCategoryStage] = createTransactionInput(transactionCategoryStage)
	inputs[transactionNoteStage] = createTransactionInput(transactionNoteStage)

	return transactionCreatorModel{
		db:        db,
		accountID: accountID,
		stage:     transactionAmountStage,
		inputs:    inputs,
		suggestion: struct {
			autocompleteSuggestion string
			payee                  ezex.Payee
			category               ezex.Category
		}{autocompleteSuggestion: ""},
		payees:     payees,
		categories: categories,
	}
}

func (m transactionCreatorModel) Update(msg tea.Msg) (transactionCreatorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !areStandardTextInputsValid(m.inputs) {
				break
			}

			notes := m.inputs[transactionNoteStage].model.Value()

			// Create new transaction
			return m, createNewTransactionCmd(
				m.db,
				ezex.Transaction{
					CategoryID:          m.suggestion.category.ID,
					PayeeID:             m.suggestion.payee.ID,
					AccountID:           m.accountID,
					AmountInCents:       decodeCents(m.inputs[transactionAmountStage].model.Value()),
					TransactionDateUnix: decodeUnixDate(m.inputs[transactionDateStage].model.Value()),
					Notes: sql.NullString{
						String: notes,
						Valid:  notes != "",
					},
				},
				ezex.Payee{
					ID:   m.suggestion.payee.ID,
					Name: m.inputs[transactionPayeeStage].model.Value(),
				},
				ezex.Category{
					ID:   m.suggestion.category.ID,
					Name: m.inputs[transactionCategoryStage].model.Value(),
				})
		case "tab":
			if m.suggestion.autocompleteSuggestion == "" {
				break
			}

			switch m.stage {
			case transactionPayeeStage:
				m.inputs[m.stage].model.SetValue(m.suggestion.payee.Name)
				m.inputs[m.stage].model.SetCursor(len(m.suggestion.payee.Name))
			case transactionCategoryStage:
				m.inputs[m.stage].model.SetValue(m.suggestion.category.Name)
				m.inputs[m.stage].model.SetCursor(len(m.suggestion.category.Name))
			}

			m.suggestion.autocompleteSuggestion = ""
		case "up", "down":
			return m.switchTransaction(msg)
		}
	}

	for i := range m.inputs {
		errMsg := m.validateInput(i)
		m.inputs[i].previousInput = m.inputs[i].model.Value()
		m.inputs[i].errorMsg = errMsg
	}
	currentInput := &m.inputs[m.stage]
	currentInput.model, cmd = currentInput.model.Update(msg)

	// No value = no autosuggestion
	val := currentInput.model.Value()
	prevVal := currentInput.previousInput
	if val == "" {
		m.suggestion.autocompleteSuggestion = ""

		if m.stage == transactionPayeeStage {
			m.suggestion.payee.ID = 0
		} else if m.stage == transactionCategoryStage {
			m.suggestion.category.ID = 0
		}

		return m, cmd
	}

	switch m.stage {
	case transactionPayeeStage:
		if val == prevVal {
			break
		}

		if match, ok := autocomplete(m.payees, val); ok {
			m.suggestion.autocompleteSuggestion = match.Name[len(val):]
			m.suggestion.payee = *match
		} else {
			m.suggestion.autocompleteSuggestion = ""
			if val != m.suggestion.payee.Name {
				m.suggestion.payee.ID = 0
			}
		}
	case transactionCategoryStage:
		if val == prevVal {
			break
		}

		if match, ok := autocomplete(m.categories, val); ok {
			m.suggestion.autocompleteSuggestion = match.Name[len(val):]
			m.suggestion.category = *match
		} else {
			m.suggestion.autocompleteSuggestion = ""
			if val != m.suggestion.category.Name {
				m.suggestion.category.ID = 0
			}
		}
	default:
		m.suggestion.autocompleteSuggestion = ""
	}

	return m, cmd
}

func (m transactionCreatorModel) View() string {
	return standardTextInputView(m.stage, m.inputs, m.suggestion.autocompleteSuggestion)
}

func (m transactionCreatorModel) switchTransaction(msg fmt.Stringer) (transactionCreatorModel, tea.Cmd) {
	m.inputs[m.stage].model.Blur()
	m.stage = handleSwitchInputStage(msg.String() == "up", m.stage, transactionDateStage, transactionNoteStage)
	m.inputs[m.stage].model.SetCursor(0)
	m.inputs[m.stage].model.Focus()

	return m, textinput.Blink
}

func (m transactionCreatorModel) validateInput(stage int) string {
	currentInput := &m.inputs[stage]
	value := currentInput.model.Value()

	// Avoid multiple checks for same inputs (because of update)
	if value == currentInput.previousInput && currentInput.previousInput != "" {
		return currentInput.errorMsg
	}

	switch stage {
	case transactionDateStage:
		if err := validateDateString(value); err != nil {
			return err.Error()
		}
	case transactionAmountStage:
		if !moneyFormatRegex.MatchString(value) {
			return "invalid amount format, should look like `0.00` or `-0.00`"
		}
	case transactionPayeeStage:
		if value == "" {
			return "payee field is required"
		}
	}

	return ""
}

func (m transactionCreatorModel) reset() transactionCreatorModel {
	m.stage = transactionAmountStage
	m.suggestion.payee.ID = 0
	m.suggestion.category.ID = 0
	m.suggestion.autocompleteSuggestion = ""

	m.inputs[transactionDateStage] = createTransactionInput(transactionDateStage)
	m.inputs[transactionAmountStage] = createTransactionInput(transactionAmountStage)
	m.inputs[transactionPayeeStage] = createTransactionInput(transactionPayeeStage)
	m.inputs[transactionCategoryStage] = createTransactionInput(transactionCategoryStage)
	m.inputs[transactionNoteStage] = createTransactionInput(transactionNoteStage)
	return m
}

func createTransactionInput(stage int) standardTextInput {
	ti := textinput.New()
	ti.Prompt = ""

	switch stage {
	case transactionDateStage:
		now := encodeUnixDate(time.Now().Unix())
		ti.Placeholder = now
		ti.SetValue(now)

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Transaction date*",
		}
	case transactionAmountStage:
		ti.Placeholder = "0.00"
		ti.Focus()

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Amount*",
		}
	case transactionPayeeStage:
		ti.Placeholder = "..."

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Payee*",
		}
	case transactionCategoryStage:
		ti.Placeholder = "No category"

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Category",
		}
	case transactionNoteStage:
		ti.Placeholder = "<NO NOTES>"

		return standardTextInput{
			model:    ti,
			errorMsg: "",
			label:    "Note",
		}
	}

	panic("unsupported transaction creation stage")
}
