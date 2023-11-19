package ezex

import "database/sql"

type Category struct {
	ID          int
	Name        string
	Description sql.NullString
}

// AddCategory creates a new category and returns the new category ID if successful
func AddCategory(db *sql.DB, category Category) (int, error) {
	return dbAdd(db, `INSERT INTO categories (name, description) VALUES ($name, $description)`, category.Name, category.Description)
}

// DeleteCategory deletes a category, trying to delete ID 0 is not allowed and will be noop
// returns the number of affected rows
func DeleteCategory(db *sql.DB, id int) int {
	if id == 0 {
		return 0
	}

	return dbDelete(db, `DELETE FROM categories WHERE id = $id`, id)
}

// UpdateCategory updates a category, trying to update ID 0 is not allowed and will be noop
func UpdateCategory(db *sql.DB, category Category) (int, error) {
	if category.ID == 0 {
		return 0, nil
	}

	return dbUpdate(
		db,
		`
		UPDATE  categories
		SET     name 		= $name,
				description = $description
		WHERE   id = $id
		`,
		category.Name,
		category.Description,
		category.ID,
	)
}

func GetCategories(db *sql.DB) []Category {
	return dbGet[Category](
		db,
		`SELECT id, name, description FROM categories ORDER BY name DESC`,
	)
}
