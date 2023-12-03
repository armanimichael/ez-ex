package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func createStandardTable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(foreground).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(selectedForeground).
		Background(selectedBackground).
		Bold(false)
	t.SetStyles(s)

	return t
}
