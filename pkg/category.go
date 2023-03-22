package models

import (
	"database/sql"
	"fmt"
	"log"
)

type CategoryModel struct {
	DB *sql.DB
}

// Get all categories
func (m *CategoryModel) GetCategories() []string {
	var name string
	var allCategories []string

	rows, err := m.DB.Query(`SELECT name FROM category`)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(rows)

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Println("Something wrong with db")
			log.Println(err)
			return nil
		}
		fmt.Println("Here1")
		allCategories = append(allCategories, name)
	}
	
	fmt.Println("Here2")
	fmt.Println(allCategories)

	return allCategories
}
