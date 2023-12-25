package command

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type SwitchModelMsg = struct {
	ModelID   int
	AccountID int
}

type HideErrorMessageMsg = struct {
	ID      int64
	Message string
}

func SwitchModelCmd(id int, accountID int) tea.Cmd {
	return func() tea.Msg {
		return SwitchModelMsg{
			ModelID:   id,
			AccountID: accountID,
		}
	}
}

func HideErrorMessageCmd(id int64, message string) tea.Cmd {
	return tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
		return HideErrorMessageMsg{
			Message: message,
			ID:      id,
		}
	})
}
