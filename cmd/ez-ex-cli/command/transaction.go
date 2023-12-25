package command

import (
	"database/sql"
	ezex "github.com/armanimichael/ez-ex"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type SwitchTransactionsMonthMsg = struct {
	Month        time.Month
	Year         int
	Transactions []ezex.TransactionView
}

type CreateNewTransactionMsg = struct {
	Transactions  []ezex.TransactionView
	NewPayee      ezex.Payee
	NewCategory   ezex.Category
	AmountInCents int64
	Err           error
}

type DeleteTransactionMsg = struct {
	DeletedID    int
	DeletedIndex int
	Err          error
}

func CreateNewTransactionCmd(db *sql.DB, transaction ezex.Transaction, payee ezex.Payee, category ezex.Category) tea.Cmd {
	return func() tea.Msg {
		if payee.ID == 0 {
			id, err := ezex.AddPayee(db, payee)
			if err != nil {
				return CreateNewTransactionMsg{Err: err}
			}
			payee.ID = id
			transaction.PayeeID = id
		}
		if category.ID == 0 && category.Name != "" {
			id, err := ezex.AddCategory(db, category)
			if err != nil {
				return CreateNewTransactionMsg{Err: err}
			}
			category.ID = id
			transaction.CategoryID = id
		}

		if _, err := ezex.AddTransaction(db, transaction); err != nil {
			return CreateNewTransactionMsg{Err: err}
		}
		if _, err := ezex.UpdateAccountBalance(db, transaction.AccountID, transaction.AmountInCents); err != nil {
			return CreateNewTransactionMsg{Err: err}
		}

		now := time.Now()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)
		transactions := ezex.GetTransactions(db, transaction.AccountID, monthStart, monthEnd)

		return CreateNewTransactionMsg{
			Transactions:  transactions,
			NewPayee:      payee,
			NewCategory:   category,
			AmountInCents: transaction.AmountInCents,
			Err:           nil,
		}
	}
}

func DeleteTransactionCmd(db *sql.DB, accountID int, id int, amountInCents int64, index int) tea.Cmd {
	return func() tea.Msg {
		var err error

		if _, err = ezex.DeleteTransaction(db, id); err != nil {
			return DeleteTransactionMsg{Err: err}
		}
		if _, err = ezex.UpdateAccountBalance(db, accountID, amountInCents); err != nil {
			return CreateNewTransactionMsg{Err: err}
		}

		return DeleteTransactionMsg{
			DeletedID:    id,
			DeletedIndex: index,
			Err:          err,
		}
	}
}

func SwitchTransactionsMonthCmd(db *sql.DB, accountID int, year int, month time.Month) tea.Cmd {
	return func() tea.Msg {
		monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		monthEnd := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
		transactions := ezex.GetTransactions(db, accountID, monthStart, monthEnd)

		return SwitchTransactionsMonthMsg{
			Month:        monthStart.Month(),
			Year:         monthStart.Year(),
			Transactions: transactions,
		}
	}
}
