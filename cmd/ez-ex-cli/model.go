package main

import (
	"database/sql"
	"fmt"
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
				m.currentModel = initTransactionModel(m.db, m.accountID)
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
