package main

import (
	ezex "github.com/armanimichael/ez-ex"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"strconv"
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

func accountsToTableRows(accounts ...ezex.Account) []table.Row {
	var rows []table.Row

	for _, account := range accounts {
		desc := account.Description.String
		if !account.Description.Valid {
			desc = "<NO DESCRIPTION>"
		}

		rows = append(
			rows,
			table.Row{
				strconv.Itoa(account.ID),
				account.Name,
				encodeCents(account.BalanceInCents, true),
				desc,
			})
	}

	return rows
}

func transactionsToTableRows(transactions ...ezex.TransactionView) []table.Row {
	var rows []table.Row

	for _, transaction := range transactions {
		date := encodeUnixDate(transaction.TransactionDateUnix)

		notes := transaction.Notes.String
		if !transaction.Notes.Valid {
			notes = "<NO NOTES>"
		}

		rows = append(
			rows,
			table.Row{
				strconv.Itoa(transaction.ID),
				date,
				encodeCents(transaction.AmountInCents, true),
				transaction.PayeeName,
				transaction.CategoryName,
				notes,
			})
	}

	return rows
}
