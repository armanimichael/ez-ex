package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	"strings"
	"time"
)

func encodeCents(cents int64, pad bool) string {
	if pad {
		return message.NewPrinter(language.English).Sprintf("%10.2f", float64(cents)/100.0)
	}

	return message.NewPrinter(language.English).Sprintf("%.2f", float64(cents)/100.0)
}

func decodeCents(cents string) int64 {
	str := strings.Replace(cents, ".", "", 1)
	c, _ := strconv.ParseInt(str, 10, 64)

	return c
}

func encodeUnixDate(unix int64) string {
	return time.Unix(unix, 0).Format(time.DateOnly)
}

func decodeUnixDate(unix string) int64 {
	year, _ := strconv.Atoi(unix[:4])
	month, _ := strconv.Atoi(unix[5:7])
	day, _ := strconv.Atoi(unix[8:])

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
}

func formatKeySuggestions(commands [][]string) string {
	str := strings.Builder{}
	for _, pair := range commands {
		str.WriteString(
			fmt.Sprintf(
				"%s\t\t%s\n",
				keySuggestionStyle.Render(pair[0]),
				lowOpacityForegroundStyle.Render(pair[1]),
			),
		)
	}

	return str.String()
}
