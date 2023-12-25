package main

import (
	"database/sql"
	"fmt"
	"github.com/armanimichael/ez-ex/cmd/ez-ex-cli/command"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	accountModelID = iota
	transactionModelID
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
	case command.SwitchModelMsg:
		if m.currentModelID != msg.ModelID {
			logger.Debug(fmt.Sprintf("Switch to transaction model ID %v", msg.ModelID))

			m.currentModelID = msg.ModelID
			if msg.AccountID != 0 {
				m.accountID = msg.AccountID
			}

			switch msg.ModelID {
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
