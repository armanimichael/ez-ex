package main

import (
	"database/sql"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	accountModelID = iota
	categoryModelID
	payeeModelID
	transactionModelID
)

type model struct {
	currentModel int
	accountModel accountModel
}

func initialModel(db *sql.DB) model {
	return model{
		currentModel: accountModelID,
		accountModel: initAccountModel(db),
	}
}

func (m model) Init() tea.Cmd {
	switch m.currentModel {
	case accountModelID:
		return m.accountModel.Init()
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentModel {
	case accountModelID:
		m.accountModel, cmd = m.accountModel.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.currentModel {
	case accountModelID:
		return m.accountModel.View()
	}

	return ""
}
