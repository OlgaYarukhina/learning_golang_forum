package models

import (
	"database/sql"
	"fmt"
	"log"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(header, category, description string, userId int) error {
	stmt := `INSERT INTO post (header, description,user_id, created_at)
    VALUES(?,?,?, current_date)`

	result, err := m.DB.Exec(stmt, header, description, userId)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Println(id)

	return nil
}


// Get all categories
func (m *PostModel) GetCategories() []string {
	var name string
	var allCategories []string

	rows, err := m.DB.Query(`SELECT * FROM category`)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(rows)

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Println("Something wrong with db")
			return nil
		}
		fmt.Println("Here1")
		allCategories = append(allCategories, name)
	}
	
	fmt.Println("Here2")
	fmt.Println(allCategories)

	return allCategories
}



// Get one user posts
func (m *PostModel) GetUserPosts(userID int) map[string][]string {
	var header, description, created_at string
	var allPosts map[string][]string

	rows, err := m.DB.Query(`SELECT header, description, created_at from post where user_id = $userID`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&header, &description, &created_at)
		if err != nil {
			log.Println("Something wrong with db")
			return nil
		}

		var val []string
		val = append(val, description)
		val = append(val, created_at)
		allPosts[header] = val
	}

	return allPosts
}
