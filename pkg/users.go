package models

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

// Insert - Метод для создания
func (m *UserModel) Insert(username, password, email string) error {

	stmt := `INSERT INTO user (username, password, email, created_at)
    VALUES(?, ?, ?, current_date)`

	_, err := m.DB.Exec(stmt, username, password, email)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) CheckUser(email string) (User, error) {
	stmt := `SELECT * FROM user WHERE email = ?`
	s := User{}
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&s.ID, &s.Username, &s.Password, &s.Email, &s.Created_at)
	if err != nil {
		return s, nil
	}

	return s, err

}
