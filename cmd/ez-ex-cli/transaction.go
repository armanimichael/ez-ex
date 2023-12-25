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
	db                 *sql.DB
	newTransaction     ezex.Transaction
	account            ezex.Account
	transactions       []ezex.TransactionView
	stage              int
	transactionCreator transactionCreatorModel
	err                struct {
		id  int64
		msg string
	}
	table struct {
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

const (
	transactionSelectionStage = iota
	transactionCreationStage
)

var transactionTableKeySuggestions = formatKeySuggestions([][]string{
	{"^C", "quit"},
	{"{esc}", "accounts list"},
	{"{right}", "next month"},
	{"{left}", "previous month"},
	{"r", "reset month"},
	{"d", "delete transaction"},
	{"n", "create transaction"},
})

func initTransactionModel(db *sql.DB, accountID int) (m transactionModel, err error) {
	m.db = db
	m.stage = transactionSelectionStage
	m.transactionCreator = initTransactionCreator(db, accountID, ezex.GetPayees(db), ezex.GetCategories(db))
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)
	m.table.selectedMonth = monthStart.Month()
	m.table.selectedYear = monthStart.Year()
	m = m.createTransactionsTable(ezex.GetTransactions(db, accountID, monthStart, monthEnd))
	if len(m.transactions) > 0 {
		m.table.selectedID = m.transactions[0].ID
	}

	m.account, err = ezex.GetAccount(db, accountID)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Cannot get account ID = %d: %v", accountID, err))
	}

	return m, err
}

func (m transactionModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m transactionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case hideErrorMessageMsg:
		if msg.message == m.err.msg && msg.id == m.err.id {
			m.err.msg = ""
		}
	case switchTransactionsMonthMsg:
		month := msg.month
		year := msg.year
		transactions := msg.transactions

		m.table.selectedMonth = month
		m.table.selectedYear = year

		if len(transactions) > 0 {
			m.table.selectedID = msg.transactions[0].ID
		}
		m.table.model.SetCursor(0)

		return m.createTransactionsTable(transactions), cmd
	case createNewTransactionMsg:
		if msg.err != nil {
			logger.Fatal(fmt.Sprintf("Error creating new transaction: %v", msg.err))
			m.err.msg = msg.err.Error()
			m.err.id = time.Now().UnixMicro()

			return m, hideErrorMessageCmd(m.err.id, m.err.msg)
		}

		m.transactionCreator = m.transactionCreator.reset()
		m.transactions = msg.transactions
		m.account.BalanceInCents += msg.amountInCents
		m.table.model.SetRows(transactionsToTableRows(msg.transactions...))
		m.stage = transactionSelectionStage
		m.table.selectedID = msg.transactions[0].ID
		m.table.model.SetCursor(0)
	case deleteTransactionMsg:
		if msg.err != nil {
			logger.Fatal(fmt.Sprintf("Error deleting transaction: %v", msg.err))
			m.err.msg = msg.err.Error()
			m.err.id = time.Now().UnixMicro()

			return m, hideErrorMessageCmd(m.err.id, m.err.msg)
		}

		// Remove deleted row from the table
		updatedTransactions := make([]ezex.TransactionView, 0, len(m.transactions)-1)
		for _, transaction := range m.transactions {
			if transaction.ID != msg.deletedID {
				updatedTransactions = append(updatedTransactions, transaction)
			} else {
				m.account.BalanceInCents -= transaction.AmountInCents
			}
		}
		m.transactions = updatedTransactions

		if len(m.transactions) != 0 {
			// Pre-select fist row
			m.table.selectedID = m.transactions[0].ID
			m.table.model.SetCursor(0)
			m.table.model.SetRows(transactionsToTableRows(m.transactions...))
		} else {
			m.table.model.SetRows([]table.Row{})
		}
		m.table.model.GotoTop()
	}

	if m.stage == transactionCreationStage {
		m.transactionCreator, cmd = m.transactionCreator.Update(msg)
		return m, cmd
	} else {
		m.table.model, cmd = m.table.model.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			logger.Debug("Go back to account list")
			return m, switchModelCmd(accountModelID, 0)
		case "right":
			next := m.table.selectedMonth + 1
			logger.Debug(fmt.Sprintf("Switch to %v", time.Date(m.table.selectedYear, next, 0, 0, 0, 0, 0, time.Local).Format("January 2006")))
			return m, tea.Batch(cmd, switchTransactionsMonthCmd(m.db, m.account.ID, m.table.selectedYear, m.table.selectedMonth+1))
		case "left":
			prev := m.table.selectedMonth - 1
			logger.Debug(fmt.Sprintf("Switch to %v", time.Date(m.table.selectedYear, prev, 0, 0, 0, 0, 0, time.Local).Format("January 2006")))
			return m, tea.Batch(cmd, switchTransactionsMonthCmd(m.db, m.account.ID, m.table.selectedYear, prev))
		case "r":
			return m, tea.Batch(cmd, switchTransactionsMonthCmd(m.db, m.account.ID, time.Now().Year(), time.Now().Month()))
		case "d":
			if len(m.transactions) == 0 {
				break
			}

			cursor := m.table.model.Cursor()
			deletedTransaction := m.transactions[cursor]
			return m, tea.Batch(
				deleteTransactionCmd(
					m.db, deletedTransaction.AccountID,
					deletedTransaction.ID,
					deletedTransaction.AmountInCents,
					cursor,
				),
				cmd,
			)
		case "n":
			m.stage = transactionCreationStage
			return m, textinput.Blink
		case "down", "up":
			r := m.table.model.SelectedRow()
			if r != nil {
				selectedID, _ := strconv.ParseInt(r[0], 10, 32)
				m.table.selectedID = int(selectedID)
			}
		}
	}

	return m, cmd
}

func (m transactionModel) View() string {
	if m.stage == transactionCreationStage {
		return m.transactionCreator.View()
	}

	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("Account:\t(ID: %d) %s\n", m.account.ID, m.account.Name))
	if m.account.Description.Valid {
		str.WriteString(fmt.Sprintf("Description:\t%s\n", m.account.Description.String))
	}
	str.WriteString(fmt.Sprintf("Balance:\t%s\n\n", encodeCents(m.account.BalanceInCents, false)))
	str.WriteString(fmt.Sprintf("Month:\t\t%s %d\n", m.table.selectedMonth.String(), m.table.selectedYear))
	str.WriteString(fmt.Sprintf("Count:\t\t%d\n", len(m.transactions)))
	str.WriteString(baseStyle.Render(m.table.model.View()) + "\n")
	str.WriteString(transactionTableKeySuggestions)

	if m.err.msg != "" {
		str.WriteString(errorMessageStyle.Render("Error: "+m.err.msg) + "\n")
	}

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
