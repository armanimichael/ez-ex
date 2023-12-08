package ezex

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddCategory(t *testing.T) {
	id1, err1 := AddCategory(testDB, Category{
		Name:        "TestAddCategory1",
		Description: sql.NullString{},
	})
	id2, err2 := AddCategory(testDB, Category{
		Name:        "TestAddCategory2",
		Description: sql.NullString{},
	})

	assert.Nil(t, err1)
	assert.Greater(t, id1, 0)
	assert.Nil(t, err2)
	assert.Greater(t, id2, 0)
}

func TestAddCategory_Unique(t *testing.T) {
	_, _ = AddCategory(testDB, Category{
		Name:        "TestAddCategory_Unique",
		Description: sql.NullString{},
	})
	_, err := AddCategory(testDB, Category{
		Name:        "TestAddCategory_Unique",
		Description: sql.NullString{},
	})

	assert.Error(t, err)
}

func TestDeleteCategory(t *testing.T) {
	id, _ := AddCategory(testDB, Category{
		Name:        "TestDeleteCategory",
		Description: sql.NullString{},
	})
	n := DeleteCategory(testDB, id)

	assert.Greater(t, n, 0)
}

func TestDeleteCategory_DefaultCategory(t *testing.T) {
	n := DeleteCategory(testDB, 0)
	assert.Equal(t, 0, n)
}

func TestUpdateCategory(t *testing.T) {
	id, _ := AddCategory(testDB, Category{
		Name:        "TestUpdateCategory",
		Description: sql.NullString{},
	})
	n, err := UpdateCategory(testDB, Category{
		ID:          id,
		Name:        "TestUpdateCategory2",
		Description: sql.NullString{},
	})

	assert.Greater(t, n, 0)
	assert.Nil(t, err)
}

func TestUpdateCategory_Unique(t *testing.T) {
	_, _ = AddCategory(testDB, Category{
		Name:        "TestUpdateCategory_Unique",
		Description: sql.NullString{},
	})
	id, _ := AddCategory(testDB, Category{
		Name:        "TestUpdateCategory_Unique2",
		Description: sql.NullString{},
	})
	n, err := UpdateCategory(testDB, Category{
		ID:          id,
		Name:        "TestUpdateCategory_Unique",
		Description: sql.NullString{},
	})

	assert.Equal(t, 0, n)
	assert.Error(t, err)
}

func TestUpdateCategory_DefaultCategory(t *testing.T) {
	n, err := UpdateCategory(testDB, Category{
		ID:          0,
		Name:        "Anything",
		Description: sql.NullString{},
	})

	assert.Equal(t, 0, n)
	assert.Nil(t, err)
}

func TestGetCategories(t *testing.T) {
	cat1 := Category{
		Name: "TestGetCategories1",
		Description: sql.NullString{
			String: "TestGetCategories1",
			Valid:  true,
		},
	}
	cat2 := Category{
		Name: "TestGetCategories2",
		Description: sql.NullString{
			String: "TestGetCategories2",
			Valid:  true,
		},
	}

	_, _ = AddCategory(testDB, cat1)
	_, _ = AddCategory(testDB, cat2)

	categories := GetCategories(testDB)

	type CategoryWithoutID struct {
		name        string
		description sql.NullString
	}
	categoriesWithoutID := make([]CategoryWithoutID, len(categories))
	for i, category := range categories {
		categoriesWithoutID[i] = CategoryWithoutID{
			name:        category.Name,
			description: category.Description,
		}
	}

	assert.Contains(t, categoriesWithoutID, CategoryWithoutID{
		name: "no category",
		description: sql.NullString{
			String: "",
			Valid:  false,
		},
	})
	assert.Contains(t, categoriesWithoutID, CategoryWithoutID{
		name:        cat1.Name,
		description: cat1.Description,
	})
	assert.Contains(t, categoriesWithoutID, CategoryWithoutID{
		name:        cat2.Name,
		description: cat2.Description,
	})
}

func TestCategory_GetName(t *testing.T) {
	cat := Category{Name: "TestCategory_GetName"}
	assert.Equal(t, "TestCategory_GetName", cat.GetName())
}
