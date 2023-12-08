package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	foreground         = lipgloss.Color("240")
	selectedForeground = lipgloss.Color("229")
	selectedBackground = lipgloss.Color("32")
	errorForeground    = lipgloss.Color("124")
	successForeground  = lipgloss.Color("2")
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var errorMessageStyle = lipgloss.NewStyle().
	Foreground(errorForeground)

var successMessageStyle = lipgloss.NewStyle().
	Foreground(successForeground)

var keySuggestionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("248"))

var lowOpacityForegroundStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240"))

var inputBoxSelectedStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(selectedBackground)

var inputBoxStyle = lipgloss.NewStyle().
	PaddingLeft(1)
