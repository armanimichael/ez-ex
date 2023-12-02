package main

import (
	"database/sql"
	tea "github.com/charmbracelet/bubbletea"
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
		m.currentModelID = msg
	case accountModel:
		m.currentModel = initAccountModel(m.db)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentModelID {
	case accountModelID:
		m.currentModel, cmd = m.currentModel.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	return m.currentModel.View()
}
