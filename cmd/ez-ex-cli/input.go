package main

import "github.com/charmbracelet/bubbles/textinput"

type standardTextInput struct {
	model         textinput.Model
	previousInput string
	errorMsg      string
	label         string
}
