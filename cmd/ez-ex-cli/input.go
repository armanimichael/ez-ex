package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"strings"
)

type standardTextInput struct {
	model         textinput.Model
	previousInput string
	errorMsg      string
	label         string
}

func areStandardTextInputsValid(inputs []standardTextInput) bool {
	for _, input := range inputs {
		if input.errorMsg != "" {
			// Errors, cannot submit
			return false
		}
	}

	return true
}

func standardTextInputView(stage int, inputs []standardTextInput, autocompleteSuggestion string) string {
	lastIndex := len(inputs) - 1
	inputListStr := strings.Builder{}
	errorsListStr := strings.Builder{}

	for i, input := range inputs {
		var render func(s ...string) string
		if i == stage {
			render = inputBoxSelectedStyle.Render
		} else {
			render = inputBoxStyle.Render
		}

		var mark string
		if input.errorMsg == "" {
			mark = successMessageStyle.Render("🗸\t")
		} else {
			mark = errorMessageStyle.Render("🞪\t")
		}
		label := mark + fmt.Sprintf("%-19s", input.label+": ")

		inputListStr.WriteString(render(label) + input.model.View())

		if autocompleteSuggestion != "" && i == stage && input.model.Value() != "" {
			inputListStr.WriteString(lowOpacityForegroundStyle.Render(autocompleteSuggestion))
		}

		if errMsg := input.errorMsg; errMsg != "" {
			errorsListStr.WriteString(errorMessageStyle.Render("- "+errMsg) + "\n")
		}

		if i != lastIndex {
			inputListStr.WriteString("\n")
		}
	}

	if errorsListStr.Len() > 0 {
		return inputListStr.String() + "\n\n" + errorsListStr.String()
	} else {
		return inputListStr.String()
	}
}

func handleSwitchInputStage(movingUp bool, currentStage int, firstStage int, lastStage int) int {
	if movingUp {
		if currentStage > firstStage {
			return currentStage - 1
		}

		return lastStage
	}

	if currentStage < lastStage {
		return currentStage + 1
	}
	return firstStage
}
