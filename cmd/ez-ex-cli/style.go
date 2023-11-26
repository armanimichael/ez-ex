package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	foreground = lipgloss.Color("240")
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var selectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")).
	Bold(false)

var tableStyle table.Styles

func init() {
	ds := table.DefaultStyles()
	ds.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(foreground).
		BorderBottom(true).
		Bold(false)
	ds.Selected = selectedStyle

	tableStyle = ds
}
