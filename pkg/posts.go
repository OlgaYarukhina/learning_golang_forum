package models

import (
	"database/sql"
	"fmt"
)

type PostModel struct {
	DB *sql.DB
}

// Insert - Метод для создания
func (m *PostModel) Insert(header, description string) error {
	stmt := `INSERT INTO post (header, description, created_at)
    VALUES(?, ?, current_date)`

	result, err := m.DB.Exec(stmt, header, description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return  err
	}
	fmt.Println(id)

	return nil
}

