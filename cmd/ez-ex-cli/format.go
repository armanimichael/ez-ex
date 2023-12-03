package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
	"time"
)

func formatCents(cents int64, pad bool) string {
	if pad {
		return message.NewPrinter(language.English).Sprintf("%10.2f", float64(cents)/100.0)
	}

	return message.NewPrinter(language.English).Sprintf("%.2f", float64(cents)/100.0)
}

func formatUnixDate(unix int64) string {
	return time.Unix(unix, 0).Format(time.DateOnly)
}

func formatKeySuggestions(commands [][]string) string {
	str := strings.Builder{}
	for _, pair := range commands {
		str.WriteString(
			fmt.Sprintf(
				"%s\t\t%s\n",
				keySuggestionStyle.Render(pair[0]),
				keySuggestionNoteStyle.Render(pair[1]),
			),
		)
	}

	return str.String()
}
