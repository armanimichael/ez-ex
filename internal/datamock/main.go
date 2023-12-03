package main

import (
	"database/sql"
	"flag"
	"fmt"
	ezex "github.com/armanimichael/ez-ex"
	"github.com/brianvoe/gofakeit/v6"
	"log"
	"strings"
	"time"
)

var faker *gofakeit.Faker

func main() {
	accountCount := flag.Int("accounts", 10, "number of mock accounts")
	categoryCount := flag.Int("categories", 10, "number of mock categories")
	payeeCount := flag.Int("payees", 10, "number of mock payees")
	transactionsCount := flag.Int("transactions", 100, "number of mock transactions per account")
	flag.Parse()

	faker = gofakeit.NewUnlocked(0)

	db, err := ezex.OpenDB()
	if err != nil {
		log.Fatalf("Error opening the DB: %s", err)
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			log.Fatalf("Error closing the DB: %s", err)
		}
	}(db)

	if err = ezex.MigrateDB(db); err != nil {
		log.Fatalf("Error migrating the DB: %s", err)
	}

	accounts := createAccounts(db, *accountCount)
	categories := createCategories(db, *categoryCount)
	payees := createPayees(db, *payeeCount)
	createTransactions(db, *transactionsCount, accounts, categories, payees)
}

func createCategories(db *sql.DB, n int) []ezex.Category {
	str := strings.Builder{}
	// language=text
	str.WriteString(`insert into categories (name, description) VALUES`)

	for i := 0; i < n; i++ {
		name := escape(faker.Word() + " " + faker.AdjectiveDescriptive())
		description := escape(faker.Phrase())

		str.WriteString(fmt.Sprintf(
			"('%s', '%s')",
			name,
			description,
		))

		if i == n-1 {
			str.WriteString(" ON CONFLICT DO NOTHING")
		} else {
			str.WriteString(",\n")
		}
	}

	query := str.String()
	_, _ = db.Exec(query)
	return ezex.GetCategories(db)
}

func createPayees(db *sql.DB, n int) []ezex.Payee {
	str := strings.Builder{}
	// language=text
	str.WriteString(`insert into payees (name, description) VALUES`)

	for i := 0; i < n; i++ {
		name := escape(faker.Company() + " " + faker.AdjectiveDescriptive())
		description := escape(faker.Phrase())

		str.WriteString(fmt.Sprintf(
			"('%s', '%s')",
			name,
			description,
		))

		if i == n-1 {
			str.WriteString(" ON CONFLICT DO NOTHING")
		} else {
			str.WriteString(",\n")
		}
	}

	query := str.String()
	_, _ = db.Exec(query)
	return ezex.GetPayees(db)
}

func createTransactions(db *sql.DB, n int, accounts []ezex.Account, categories []ezex.Category, payees []ezex.Payee) {
	str := strings.Builder{}
	// language=text
	initQuery := `insert into transactions (category_id, payee_id, account_id, amount_in_cents, transaction_date_unix, notes) VALUES`

	for _, account := range accounts {
		str.WriteString(initQuery)
		for j := 0; j < n; j++ {
			categoryID := categories[faker.IntRange(0, len(categories)-1)].ID
			payeeID := payees[faker.IntRange(0, len(payees)-1)].ID
			accountID := account.ID
			amount := faker.IntRange(1, 10_000*100)
			dateUnix := faker.DateRange(time.Now(), time.Now().AddDate(1, 0, 0)).Unix()
			notes := escape(faker.Phrase())

			str.WriteString(fmt.Sprintf(
				"(%d, %d, %d, %d, %d, '%s')",
				categoryID,
				payeeID,
				accountID,
				amount,
				dateUnix,
				notes,
			))

			if j == n-1 {
				str.WriteString(" ON CONFLICT DO NOTHING")
			} else {
				str.WriteString(",\n")
			}
		}
		query := str.String()
		_, _ = db.Exec(query)
		str.Reset()
	}
}

func createAccounts(db *sql.DB, n int) []ezex.Account {
	str := strings.Builder{}
	// language=text
	str.WriteString(`insert into accounts (name, description, balance_in_cents, initial_balance_in_cents) VALUES`)

	for i := 0; i < n; i++ {
		name := escape(faker.Word() + " " + faker.AdjectiveDescriptive())
		description := escape(faker.Phrase())
		balance := faker.IntRange(1, 100_000*100)
		str.WriteString(fmt.Sprintf(
			"('%s', '%s', %d, %d)",
			name,
			description,
			balance,
			balance,
		))

		if i == n-1 {
			str.WriteString(" ON CONFLICT DO NOTHING")
		} else {
			str.WriteString(",\n")
		}
	}

	query := str.String()
	_, _ = db.Exec(query)
	return ezex.GetAccounts(db)
}

func escape(str string) string {
	return strings.Replace(str, "'", "''", -1)
}
