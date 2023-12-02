package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	foreground         = lipgloss.Color("240")
	selectedForeground = lipgloss.Color("229")
	selectedBackground = lipgloss.Color("32")
	errorForeground    = lipgloss.Color("124")
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var errorMessageStyle = lipgloss.NewStyle().
	Foreground(errorForeground)

var keySuggestionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("248"))

var keySuggestionNoteStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240"))
