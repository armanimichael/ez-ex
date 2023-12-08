package ezex

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPayee(t *testing.T) {
	id1, err1 := AddPayee(testDB, Payee{
		Name:        "TestAddPayee1",
		Description: sql.NullString{},
	})
	id2, err2 := AddPayee(testDB, Payee{
		Name:        "TestAddPayee2",
		Description: sql.NullString{},
	})

	assert.Nil(t, err1)
	assert.Greater(t, id1, 0)
	assert.Nil(t, err2)
	assert.Greater(t, id2, 0)
}

func TestAddPayee_Unique(t *testing.T) {
	_, _ = AddPayee(testDB, Payee{
		Name:        "TestAddPayee_Unique",
		Description: sql.NullString{},
	})
	_, err := AddPayee(testDB, Payee{
		Name:        "TestAddPayee_Unique",
		Description: sql.NullString{},
	})

	assert.Error(t, err)
}

func TestDeletePayee(t *testing.T) {
	id, _ := AddPayee(testDB, Payee{
		Name:        "TestDeletePayee",
		Description: sql.NullString{},
	})
	n := DeletePayee(testDB, id)

	assert.Greater(t, n, 0)
}

func TestUpdatePayee(t *testing.T) {
	id, _ := AddPayee(testDB, Payee{
		Name:        "TestUpdatePayee",
		Description: sql.NullString{},
	})
	n, err := UpdatePayee(testDB, Payee{
		ID:          id,
		Name:        "TestUpdatePayee2",
		Description: sql.NullString{},
	})

	assert.Greater(t, n, 0)
	assert.Nil(t, err)
}

func TestUpdatePayee_Unique(t *testing.T) {
	_, _ = AddPayee(testDB, Payee{
		Name:        "TestUpdatePayee_Unique",
		Description: sql.NullString{},
	})
	id, _ := AddPayee(testDB, Payee{
		Name:        "TestUpdatePayee_Unique2",
		Description: sql.NullString{},
	})
	n, err := UpdatePayee(testDB, Payee{
		ID:          id,
		Name:        "TestUpdatePayee_Unique",
		Description: sql.NullString{},
	})

	assert.Equal(t, 0, n)
	assert.Error(t, err)
}

func TestGetPayees(t *testing.T) {
	payee1 := Payee{
		Name: "TestGetPayees1",
		Description: sql.NullString{
			String: "TestGetPayees1",
			Valid:  true,
		},
	}
	payee2 := Payee{
		Name: "TestGetPayees2",
		Description: sql.NullString{
			String: "TestGetPayees2",
			Valid:  true,
		},
	}

	_, _ = AddPayee(testDB, payee1)
	_, _ = AddPayee(testDB, payee2)

	payees := GetPayees(testDB)

	type PayeeWithoutID struct {
		name        string
		description sql.NullString
	}
	payeesWithoutID := make([]PayeeWithoutID, len(payees))
	for i, payee := range payees {
		payeesWithoutID[i] = PayeeWithoutID{
			name:        payee.Name,
			description: payee.Description,
		}
	}

	assert.Contains(t, payeesWithoutID, PayeeWithoutID{
		name:        payee1.Name,
		description: payee1.Description,
	})
	assert.Contains(t, payeesWithoutID, PayeeWithoutID{
		name:        payee2.Name,
		description: payee2.Description,
	})
}

func TestPayee_GetName(t *testing.T) {
	payee := Payee{Name: "TestPayee_GetName"}
	assert.Equal(t, "TestPayee_GetName", payee.GetName())
}
