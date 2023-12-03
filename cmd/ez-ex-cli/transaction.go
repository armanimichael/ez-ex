package main

import (
	"database/sql"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
	"strings"
	"time"
)

type transactionModel struct {
	db             *sql.DB
	newTransaction ezex.Transaction
	account        ezex.Account
	transactions   []ezex.TransactionView
	table          struct {
		model         table.Model
		selectedID    int
		selectedMonth time.Month
		selectedYear  int
	}
	input struct {
		model         textinput.Model
		previousInput string
		errorMsg      string
		label         string
	}
}

var transactionTableKeySuggestions = formatKeySuggestions([][]string{
	{"^C", "quit"},
	{"{esc}", "accounts list"},
	{"{enter}", "select transaction"},
	{"{right}", "next month"},
	{"{left}", "previous month"},
	{"r", "reset month"},
	{"d", "delete transaction"},
	{"n", "create transaction"},
})

func initTransactionModel(db *sql.DB, accountID int) (m transactionModel) {
	m.db = db
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)
	m.table.selectedMonth = monthStart.Month()
	m.table.selectedYear = monthStart.Year()
	m = m.createTransactionsTable(ezex.GetTransactions(db, accountID, monthStart, monthEnd))

	var err error
	m.account, err = ezex.GetAccount(db, accountID)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Cannot get account ID = %d: %v", accountID, err))
	}

	return m
}

func (m transactionModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m transactionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table.model, cmd = m.table.model.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			logger.Debug(fmt.Sprintf("Go back to account list"))
			return m, switchModelCmd(accountModelID, 0)
		case "right":
			return m, tea.Batch(cmd, switchTransactionsMonthCmd(m.db, m.account.ID, m.table.selectedYear, m.table.selectedMonth+1))
		case "left":
			return m, tea.Batch(cmd, switchTransactionsMonthCmd(m.db, m.account.ID, m.table.selectedYear, m.table.selectedMonth-1))
		case "r":
			return m, tea.Batch(cmd, switchTransactionsMonthCmd(m.db, m.account.ID, time.Now().Year(), time.Now().Month()))
		}
	case switchTransactionsMonthMsg:
		month := msg.month
		year := msg.year
		transactions := msg.transactions

		m.table.selectedMonth = month
		m.table.selectedYear = year
		return m.createTransactionsTable(transactions), cmd
	}

	return m, cmd
}

func (m transactionModel) View() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("Account:\t(ID: %d) %s\n", m.account.ID, m.account.Name))
	if m.account.Description.Valid {
		str.WriteString(fmt.Sprintf("Description:\t%s\n", m.account.Description.String))
	}
	str.WriteString(fmt.Sprintf("Balance:\t%s\n\n", formatCents(m.account.BalanceInCents, false)))
	str.WriteString(fmt.Sprintf("Month:\t\t%s %d\n", m.table.selectedMonth.String(), m.table.selectedYear))
	str.WriteString(fmt.Sprintf("Count:\t\t%d\n", len(m.transactions)))
	str.WriteString(baseStyle.Render(m.table.model.View()) + "\n")
	str.WriteString(transactionTableKeySuggestions)

	return str.String()
}

func (m transactionModel) createTransactionsTable(transactions []ezex.TransactionView) transactionModel {
	m.newTransaction = ezex.Transaction{}
	m.transactions = transactions
	m.table.model = createStandardTable(
		[]table.Column{
			{Title: "ID", Width: 5},
			{Title: "Date", Width: 10},
			{Title: "Amount", Width: 10},
			{Title: "Payee", Width: 20},
			{Title: "Category", Width: 20},
			{Title: "Notes", Width: 40},
		},
		transactionsToTableRows(transactions...),
	)
	return m
}

func transactionsToTableRows(transactions ...ezex.TransactionView) []table.Row {
	var rows []table.Row

	for _, transaction := range transactions {
		date := formatUnixDate(transaction.TransactionDateUnix)

		notes := transaction.Notes.String
		if !transaction.Notes.Valid {
			notes = "<NO NOTES>"
		}

		rows = append(
			rows,
			table.Row{
				strconv.Itoa(transaction.ID),
				date,
				formatCents(transaction.AmountInCents, true),
				transaction.PayeeName,
				transaction.CategoryName,
				notes,
			})
	}

	return rows
}
