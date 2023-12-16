package main

import (
	"database/sql"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

const (
	accountModelID = iota
	transactionModelID
	categoryModelID
	payeeModelID
)

type model struct {
	db             *sql.DB
	currentModelID int
	accountID      int
	currentModel   tea.Model
}

func initialModel(db *sql.DB) model {
	return model{
		db:             db,
		currentModelID: accountModelID,
		currentModel:   initAccountModel(db),
	}
}

func (m model) Init() tea.Cmd {
	switch m.currentModelID {
	case accountModelID:
		return m.currentModel.Init()
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case switchModelMsg:
		if m.currentModelID != msg.modelID {
			logger.Debug(fmt.Sprintf("Switch to transaction model ID %v", msg.modelID))

			m.currentModelID = msg.modelID
			if msg.accountID != 0 {
				m.accountID = msg.accountID
			}

			switch msg.modelID {
			case accountModelID:
				m.currentModel = initAccountModel(m.db)
			case transactionModelID:
				var err error
				m.currentModel, err = initTransactionModel(m.db, m.accountID)
				if err != nil {
					return m, tea.Quit
				}
			}

			return m, cmd
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.currentModel.View()
}

func areStandardTextInputsValid(inputs []standardTextInput) bool {
	for _, input := range inputs {
		if input.errorMsg != "" {
			// Errors, cannot submit
			return false
		}
	}

	return true
}

func standardTextInputView(stage int, inputs []standardTextInput, autocompleteSuggestion string) string {
	lastIndex := len(inputs) - 1
	inputListStr := strings.Builder{}
	errorsListStr := strings.Builder{}

	for i, input := range inputs {
		var render func(s ...string) string
		if i == stage {
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

		if autocompleteSuggestion != "" && i == stage && input.model.Value() != "" {
			inputListStr.WriteString(lowOpacityForegroundStyle.Render(autocompleteSuggestion))
		}

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

func handleSwitchInputStage(movingUp bool, currentStage int, firstStage int, lastStage int) int {
	if movingUp {
		if currentStage > firstStage {
			return currentStage - 1
		}

		return lastStage
	}

	if currentStage < lastStage {
		return currentStage + 1
	}
	return firstStage
}
