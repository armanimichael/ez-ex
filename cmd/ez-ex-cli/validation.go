package main

import (
	"errors"
	"regexp"
	"strconv"
)

var moneyFormatRegex = regexp.MustCompile(`^-?(?P<integer>\d+)(\.(?P<cents>\d{2}))+$`)
var dateFormatRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func validateDateString(value string) error {
	if !dateFormatRegex.MatchString(value) {
		return errors.New("invalid date format, should be YYYY-MM-DD")
	}

	year, _ := strconv.Atoi(value[:4])
	month, _ := strconv.Atoi(value[5:7])
	day, _ := strconv.Atoi(value[8:])

	if month > 12 || month == 0 {
		return errors.New("invalid month (01-12)")
	}

	isLeapYear := year%4 == 0
	if month == 2 && isLeapYear && day > 29 {
		return errors.New("invalid day (01-29)")
	}
	if month == 2 && !isLeapYear && day > 28 {
		return errors.New("invalid day (01-28)")
	}

	// 31 days months
	longMonths := map[int]struct{}{1: {}, 3: {}, 5: {}, 7: {}, 8: {}, 10: {}, 12: {}}
	if _, ok := longMonths[month]; ok && day > 31 {
		return errors.New("invalid day (01-31)")
	} else if !ok && day > 30 {
		return errors.New("invalid day (01-30)")
	}

	return nil
}
