package models

import (
	"database/sql"
	"log"
)

type CategoryModel struct {
	DB *sql.DB
}

// Get all categories
func (m *CategoryModel) GetCategories() ([]*Category, error) {
	var allCategories []*Category
	rows, err := m.DB.Query(`SELECT * FROM category`)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		category := &Category{}
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return allCategories, err
		}
		allCategories = append(allCategories, category)
	}
	return allCategories, nil
}
