package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMoneyFormatRegex(t *testing.T) {
	cases := []struct {
		value   string
		isValid bool
	}{
		{"0.00", true},
		{"1.00", true},
		{"-1.00", true},
		{"123.12", true},
		{"-123.12", true},
		{"123.00", true},
		{"0,00", false},
		{"1.0", false},
		{".0", false},
		{"123", false},
		{"-123", false},
		{".10", false},
		{"0", false},
		{"-.12", false},
	}
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			isValid := moneyFormatRegex.MatchString(c.value)
			assert.Equal(t, c.isValid, isValid)
		})
	}
}

func TestValidateDateString(t *testing.T) {
	cases := []struct {
		value   string
		isValid bool
	}{
		{"2023-01-31", true},
		{"2023-02-28", true},
		{"2023-03-31", true},
		{"2023-04-30", true},
		{"2023-05-31", true},
		{"2023-06-30", true},
		{"2023-07-31", true},
		{"2023-08-31", true},
		{"2023-09-30", true},
		{"2023-10-31", true},
		{"2023-11-30", true},
		{"2023-12-31", true},
		{"2024-12-29", true},
		{"2023/01/32", false},
		{"2023-01-32", false},
		{"2023-02-29", false},
		{"2023-03-32", false},
		{"2023-04-31", false},
		{"2023-05-32", false},
		{"2023-06-31", false},
		{"2023-07-32", false},
		{"2023-08-32", false},
		{"2023-09-31", false},
		{"2023-10-32", false},
		{"2023-11-31", false},
		{"2023-12-32", false},
		{"2023-13-10", false},
		{"2024-02-30", false},
	}
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			err := validateDateString(c.value)
			isValid := err == nil
			assert.Equal(t, c.isValid, isValid)
		})
	}
}
